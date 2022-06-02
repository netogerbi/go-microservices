package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const (
	dbName         = "logger"
	collectionName = "logs"
	timeout        = 15 * time.Second
)

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"createAt" json:"createAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database(dbName).Collection(collectionName)

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error trying to insert LogEntry: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	collection := client.Database(dbName).Collection(collectionName)

	options := options.Find().SetSort(bson.D{
		primitive.E{Key: "createdAt", Value: -1},
	})

	cursor, err := collection.Find(ctx, options)
	if err != nil {
		log.Println("Error trying to find all LogEntry: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(item)
		if err != nil {
			log.Println("Error trying to decode LogEntry: ", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) FindOne(id string) (*LogEntry, error) {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()

	collection := client.Database(dbName).Collection(collectionName)

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error trying to convert string to ObjectID in LogEntry: ", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Println("Error trying to find or decode LogEntry: ", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()

	collection := client.Database(dbName).Collection(collectionName)

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error trying to convert string to ObjectID in LogEntry: ", err)
		return nil, err
	}

	r, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: l.Name},
			primitive.E{Key: "data", Value: l.Data},
			primitive.E{Key: "updatedAt", Value: time.Now()},
		}}},
	)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (l *LogEntry) dropCollection() error {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()

	collection := client.Database(dbName).Collection(collectionName)
	if err := collection.Drop(ctx); err != nil {
		log.Println("Error trying to drop collection: ", err)
		return err
	}

	return nil
}
