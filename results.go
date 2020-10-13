package main

import (
    "log"
    "net/http"

    "github.com/mercimat/instavote/core"
    "github.com/mercimat/instavote/db"
)

func main() {
    mongoServer := "mongodb://localhost:27017/"
    mdb := db.NewMongoDB(mongoServer, "instavote", "votes")

    err := core.Init(
        "InstaVote Results",
        "Dogs",
        "Cats",
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
