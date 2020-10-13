package main

import (
    "log"
    "net/http"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/redis"
)

func main() {

    err := core.Init(
        "InstaVote App",
        "Dogs",
        "Cats",
        "templates/vote.html",
    )
    if err != nil {
        panic(err)
    }

    rdb := redis.NewRedisCon(
        "localhost:6379",
        "", // no password set
        0, // use default DB
    )

    url := ":8090"

    http.HandleFunc("/", core.MakeVoteHandler(core.VoteHandler, rdb))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    log.Fatal(http.ListenAndServe(url, nil))
}
