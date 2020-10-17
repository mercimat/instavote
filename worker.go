package main

import (
    "log"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/db"
    "github.com/mercimat/instavote/redis"
)

func main() {

    rdb := redis.NewRedisCon(
        "localhost:6379",
        "", // no password set
        0, // use default DB
    )

    mdb := db.NewMongoDB(
        "mongodb://localhost:27017/",
        "instavote",
        "votes",
    )

    // Call Redis2Mongo in a loop
    for {
        err := core.Redis2Mongo(rdb, mdb)
        if err != nil {
            log.Printf("Failed to get or parse Redis message: %s\n", err)
        }
    }
}
