package main

import (
    "log"
    "os"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/db"
    "github.com/mercimat/instavote/redis"
)

func main() {

    redisServer := "localhost:6379"
    if _, ok := os.LookupEnv("REDIS_HOST"); ok {
        redisServer = os.ExpandEnv("${REDIS_HOST}:6379")
    }
    rdb := redis.NewRedisCon(
        redisServer,
        "", // no password set
        0, // use default DB
    )

    mongoServer := "mongodb://localhost:27017/"
    if _, ok := os.LookupEnv("MONGODB_HOST"); ok {
        mongoServer = os.ExpandEnv("mongodb://${MONGODB_HOST}:27017/")
    }
    mdb := db.NewMongoDB(
        mongoServer,
        "instavote",
        "votes",
    )

    for {
        err := core.Redis2Mongo(rdb, mdb)
        if err != nil {
            log.Printf("Failed to get or parse Redis message: %s\n", err)
        }
    }
}
