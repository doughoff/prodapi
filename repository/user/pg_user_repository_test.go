package user

import (
	"context"
	"fmt"
	"github.com/hoffax/prodapi/dbtypes"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func beforeAfter(t *testing.T) (*pgxpool.Pool, func(), func()) {
	ctx := context.Background()

	migrator, err := dbtypes.NewTestDBMigrator()
	if err != nil {
		t.Fatalf("Failed to create migrator: %v", err)
	}

	err = migrator.Prepare()
	if err != nil {
		t.Fatalf("Failed to prepare database: %v", err)
	}

	dbpool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}

	return dbpool, func() {
			err := migrator.Prepare()
			fmt.Printf("error preparing test-db\n%v\n", err)
		}, func() {
			dbpool.Close()
			err = migrator.Cleanup()
			if err != nil {
				t.Fatalf("Failed to cleanup database: %v", err)
			}
		}
}

func TestAllUsers(t *testing.T) {
	dbpool, migrator, cleanup := beforeAfter(t)
	defer cleanup()

	ctx := context.Background()
	repo := NewUserPgRepository(dbpool)

	migrator()

	users, err := repo.All(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// perform checks depending on your business logic, e.g.
	if len(users) != 0 {
		t.Errorf("expected no users, got %d", len(users))
	}
}
