package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// ConnectPostgresDB establishes a connection to the PostgreSQL database using `database/sql`,
// retries if necessary, and returns a connection pool.
// Parameters:
// - ctx: The context to control cancellation and timeouts.
// - dsn: The Data Source Name (DSN) or database URL required to connect to the PostgreSQL database.
// - timeout: The maximum duration to wait before giving up on connecting.
// - maxRetries: The number of times to retry connecting in case of failure.
// Returns:
// - *sql.DB: A pointer to the database connection pool on a successful connection.
// - error: An error if the connection fails after all retries.
func ConnectPostgresDB(ctx context.Context, dsn string, timeout time.Duration, maxRetries int) (*sql.DB, error) {
	// Set a timeout for the connection operation using the context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Validate the DSN (database URL) input
	if dsn == "" {
		return nil, fmt.Errorf("missing required database URL (DSN)")
	}

	var db *sql.DB
	var err error
	retryDelay := 5 * time.Second // Time to wait between retries

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			// If context times out or is canceled, exit with an error
			return nil, fmt.Errorf("context timed out while trying to connect to database: %w", ctx.Err())
		default:
			// Try to open the connection using the standard library's sql package
			db, err = sql.Open("postgres", dsn)
			if err == nil {
				// Ping the database to ensure connection is established
				err = db.PingContext(ctx)
				if err == nil {
					log.Println("Connected to PostgreSQL successfully")
					return db, nil // Return the connected DB instance
				}
			}

			// Log the failure and retry after a delay
			log.Printf("Connection attempt %d failed: %v", i+1, err)
			time.Sleep(retryDelay) // Wait before the next retry
		}
	}

	// Return error if all retries fail
	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err)
}

// ctx := context.Background()
// dsn := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
// timeout := 30 * time.Second
// maxRetries := 3

// db, err := postgres.ConnectPostgresDB(ctx, dsn, timeout, maxRetries)
// if err != nil {
// 	log.Fatalf("Error connecting to the database: %v", err)
// }

// defer db.Close() // Always ensure to close the database connection when done
