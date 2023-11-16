# MongoDB Trigger for Go

MongoDB Trigger is a Go package designed to simplify the integration and management of MongoDB triggers, allowing seamless handling of insert, update, and delete operations via change streams.

## Installation

To use this package, ensure you have Go installed. Then, import the package into your Go project:

```bash
go get github.com/crettien/mongodbtrigger
```

## Usage

## Import the package in your Go code

```go
import (
    "github.com/crettien/mongodbtrigger"
)
```

## Example

## Initialize the MongoDB configuration and listen for operations

```go
package main

import (
    "github.com/your-username/mongodbtrigger"
)

func main() {
    // Configure MongoDB
    config := mongodbtrigger.MongoDBConfig{
        Cluster:    "localhost",
        Username:   "username",
        Password:   "password",
        Database:   "databaseName",
        Collection: "collectionName",
        Trigger:    2, // 1: Insert, 2: Update, 3: Delete
    }

    // Listen for MongoDB operations
    mongodbtrigger.ListenForOperations(config)
}
```

## MongoDBConfig Structure

The `MongoDBConfig` struct holds the MongoDB configuration details:

- `Cluster`: MongoDB cluster address.
- `Username`: MongoDB account username.
- `Password`: MongoDB account password.
- `Database`: Name of the MongoDB database.
- `Collection`: Name of the MongoDB collection.
- `Trigger`: Type of trigger operation (1: Insert, 2: Update, 3: Delete).

## Contributing

Contributions are welcome! Feel free to open issues or pull requests.

## License

This project is licensed under the MIT License.
