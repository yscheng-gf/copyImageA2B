package helper

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(host string) *mongo.Client {
	mongoCli, err := mongo.Connect(
		context.TODO(),
		options.Client().SetHosts(strings.Split(host, ",")),
	)
	if err != nil {
		panic(err)
	}
	return mongoCli
}
