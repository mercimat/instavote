package main

import (
    "flag"
    "log"
    "net/http"
    "os"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/db"
)

func main() {
    mongoServer := "mongodb://localhost:27017/"
    if _, ok := os.LookupEnv("MONGODB_HOST"); ok {
        mongoServer = os.ExpandEnv("mongodb://${MONGODB_HOST}:27017/")
    }
    mdb := db.NewMongoDB(
        mongoServer,
        "instavote",
        "votes",
    )

    optA := flag.String("a", "Dogs", "Option A")
    optB := flag.String("b", "Cats", "Option B")
    flag.Parse()

    err := core.Init(
        "InstaVote Results",
        *optA,
        *optB,
        "templates/results.html",
    )
    if err != nil {
        panic(err)
    }

    url := ":8091"

    http.HandleFunc("/", core.MakeResultsHandler(core.ResultsHandler, mdb))
    http.HandleFunc("/api/results", core.MakeResultsHandler(core.ApiResultsHandler, mdb))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    log.Fatal(http.ListenAndServe(url, nil))
}
