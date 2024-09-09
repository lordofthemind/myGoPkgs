package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB establishes a connection to the MongoDB database, retries if necessary,
// and allows you to specify the connection string as a parameter.
// Parameters:
// - ctx: The context to control cancellation and timeouts.
// - dsn: The MongoDB connection string (Data Source Name).
// - timeout: The maximum duration to wait before giving up on connecting.
// - maxRetries: The number of times to retry connecting in case of failure.
// Returns:
// - *mongo.Client: A pointer to the MongoDB client instance on a successful connection.
// - *mongo.Database: A pointer to the MongoDB database instance.
// - error: An error if the connection fails after all retries.
func ConnectToMongoDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int, dbName string) (*mongo.Client, *mongo.Database, error) {
	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Validate the DSN (connection string) input
	if dsn == "" {
		return nil, nil, fmt.Errorf("missing required MongoDB connection string (DSN)")
	}

	var client *mongo.Client
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, nil, fmt.Errorf("context timed out while trying to connect to MongoDB: %w", ctx.Err())
		default:
			// Try to establish a connection to MongoDB
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn))
			if err == nil {
				// Successfully connected, return the client and the database instance
				log.Println("Connected to MongoDB successfully")
				db := client.Database(dbName)
				return client, db, nil
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v\n", i+1, err)
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Return error if all retries fail
	return nil, nil, fmt.Errorf("failed to connect to MongoDB after %d retries: %w", maxRetries, err)
}

// ctx := context.Background()
// dsn := os.Getenv("MONGO_DOCKER_CONNECTION_URL")
// timeout := 30 * time.Second
// maxRetries := 5
// dbName := "polyglot" // Replace with your actual database name

// client, db, err := mypkg.ConnectMongoDB(ctx, dsn, timeout, maxRetries, dbName)
// if err != nil {
// 	log.Fatalf("Error connecting to MongoDB: %v", err)
// }
