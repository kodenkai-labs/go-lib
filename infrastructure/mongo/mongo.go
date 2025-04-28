package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type DBMode int

const (
	DBModeWrite DBMode = 1
	DBModeRead  DBMode = 2
)

const timeout = 5 * time.Second

func New(cfg Config, mode DBMode) (*mongo.Database, error) {
	logrus.Info("Connecting to mongodb")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime)

	switch mode {
	case DBModeWrite:
		// https://www.mongodb.com/docs/manual/reference/write-concern/
		cOpts = cOpts.SetWriteConcern(writeconcern.W1())
	case DBModeRead:
		cOpts = cOpts.SetReadPreference(readpref.Secondary())
	default:
		return nil, fmt.Errorf("unknown DBMode: %v", mode)
	}

	client, err := mongo.Connect(ctx, cOpts)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	logrus.Info("Successfully established connection to mongodb")

	return client.Database(cfg.Name), nil
}
