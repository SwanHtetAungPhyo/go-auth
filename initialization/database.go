package initialization

import (
	"context"
	"time"

	db "github.com/SwanHtetAungPhyo/go-auth/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Database(dsn string, redisAsSessionStore bool) *db.Store {
	pgxPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	store := db.NewStore(pgxPool)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := runMigrations(ctx, store, dsn, redisAsSessionStore); err != nil {
		log.Fatal().Err(err).Msg("goauth: database migrations failed")
	}

	return store
}

func runMigrations(ctx context.Context, store *db.Store, dsn string, redisAsSessionStore bool) error {
	log.Info().Msgf("goauth: creating database with %s", dsn)

	if err := store.CreateUserTable(ctx); err != nil {
		return err
	}
	if err := store.CreateUserIndexes(ctx); err != nil {
		return err
	}
	if err := store.CreateAccountTable(ctx); err != nil {
		return err
	}
	if err := store.CreateAccountIndexes(ctx); err != nil {
		return err
	}
	if err := store.CreateAuditLogTable(ctx); err != nil {
		return err
	}
	if !redisAsSessionStore {
		if err := store.CreateSessionTable(ctx); err != nil {
			return err
		}
		if err := store.CreateSessionIndexes(ctx); err != nil {
			return err
		}
		if err := store.CreatePasswordResetTable(ctx); err != nil {
			return err
		}
		if err := store.CreateEmailVerificationTable(ctx); err != nil {
			return err
		}

		log.Info().Msg("goauth: database migrations succeeded")
	}

	return nil
}
