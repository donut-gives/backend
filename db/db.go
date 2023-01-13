package db

import (
	"context"
	"donutBackend/config"
	"donutBackend/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

var db *mongo.Database

func Get() *mongo.Database {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB.Url))
			if err != nil {
				logger.Logger.Fatal(err)
			}
			db = client.Database(config.DB.Name)
		}

	}
	return db
}
