package mongodbtrigger

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBConfig holds MongoDB configuration details.
type MongoDBConfig struct {
	Cluster    string
	Username   string
	Password   string
	Database   string
	Collection string
	Trigger    int
}

// OperationHandler defines a function to handle specific operations.
type OperationHandler func()

// ChangeStreamCallback defines a function to handle change stream documents.
type ChangeStreamCallback func(changeDoc bson.M)

func ListenForOperations(config MongoDBConfig, changeStreamCallback ChangeStreamCallback) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(createURI(config)))
	if err != nil {
		log.Fatalf("unable to connect to the database: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("error while disconnecting: %v", err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("unable to ping the database: %v", err)
	}

	dbCollection := client.Database(config.Database).Collection(config.Collection)

	triggerHandlers := map[string]OperationHandler{
		"insert": func() { waitForOperation(ctx, dbCollection, "insert", changeStreamCallback) },
		"update": func() { waitForOperation(ctx, dbCollection, "update", changeStreamCallback) },
		"delete": func() { waitForOperation(ctx, dbCollection, "delete", changeStreamCallback) },
	}

	operations := getTriggerOperations(config.Trigger)
	var wg sync.WaitGroup
	for _, operation := range operations {
		wg.Add(1)
		go func(op string) {
			defer wg.Done()
			triggerHandlers[op]()
		}(operation)
	}
	wg.Wait()
}

func waitForOperation(ctx context.Context, collection *mongo.Collection, operationType string, changeStreamCallback ChangeStreamCallback) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "operationType", Value: operationType}}}},
	}
	opts := options.ChangeStream().SetFullDocument("updateLookup")

	changeStream, err := collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Printf("error while creating change stream: %v", err)
		return
	}
	defer changeStream.Close(ctx)

	fmt.Printf("Watching for document %s events...\n", operationType)

	for {
		if changeStream.Next(ctx) {
			var changeDoc bson.M
			if err := changeStream.Decode(&changeDoc); err != nil {
				log.Println(err)
				continue
			}
			changeStreamCallback(changeDoc) // Pass the changeDoc to the callback function
		} else if err := changeStream.Err(); err != nil {
			log.Printf("error in change stream for %s operation: %v", operationType, err)
			break
		}
	}
}

func createURI(config MongoDBConfig) string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", config.Username, config.Password, config.Cluster)
}

func getTriggerOperations(trigger int) []string {
	switch trigger {
	case 1:
		return []string{"insert"}
	case 2:
		return []string{"update"}
	case 3:
		return []string{"insert", "update"}
	case 4:
		return []string{"delete"}
	case 5:
		return []string{"insert", "delete"}
	case 6:
		return []string{"update", "delete"}
	case 7:
		return []string{"insert", "update", "delete"}
	default:
		log.Printf("%d is not a valid trigger operation type", trigger)
		return []string{}
	}
}
