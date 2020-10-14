package main

import (
    "flag"
    "log"
    "net/http"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/db"
)

func main() {
    mongoServer := "mongodb://localhost:27017/"
    mdb := db.NewMongoDB(mongoServer, "instavote", "votes")

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

    http.HandleFunc("/", core.MakeResultsHandler(core.ResultHandler, mdb))
    http.HandleFunc("/api/results", core.MakeResultsHandler(core.ApiResultHandler, mdb))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    log.Fatal(http.ListenAndServe(url, nil))
}
