package db

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
    Collection  *mongo.Collection
    Ctx         context.Context
}

type Counter struct {
    ID      string  `bson:"_id,omitempty"`
    Total   int     `bson:"total,omitempty"`
}

func NewMongoDB(server string, database string, collection string) *MongoDB {
    clientOptions := options.Client().ApplyURI(server)
    ctx := context.TODO()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    return &MongoDB{
        Collection: client.Database(database).Collection(collection),
        Ctx: ctx,
    }
}

//filter := db.FilterBson("voterid", message.VoterId)
//update := db.UpdateBson("vote", message.Vote)
func (m MongoDB) Set(elm interface{}, filter interface{}, update interface{}) error {
    opts := options.FindOneAndUpdate().SetUpsert(true)
    err := m.Collection.FindOneAndUpdate(m.Ctx, filter, update, opts).Decode(&elm)
    if err != nil && err != mongo.ErrNoDocuments {
        return err
    }
    return nil
}

// filter := bson.D{{"voterid", id}}
func (m MongoDB) Get(elm interface{}, filter interface{}) error {
    return m.Collection.FindOne(m.Ctx, filter).Decode(elm)
}

func (m MongoDB) Count(key string) (map[string]int, error) {
    results := make(map[string]int)

    pipeline := []bson.M{bson.M{"$group": bson.M{"_id": key, "total": bson.M{"$sum": 1}}}}
    cursor, err := m.Collection.Aggregate(m.Ctx, pipeline)
    if err != nil {
        return results, err
    }

    var counters []Counter
    err = cursor.All(m.Ctx, &counters)
    if err != nil {
        return results, err
    }

    for _, counter := range counters {
        results[counter.ID] = counter.Total
    }

    return results, nil
}

func (m MongoDB) QueryFilter(key string, value interface{}) interface{} {
    return bson.D{{key, value}}
}

func (m MongoDB) QueryUpdate(key string, value interface{}) interface{} {
    return bson.D{{"$set", bson.D{{key, value}}}}
}

