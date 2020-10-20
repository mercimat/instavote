package core

import (
    "encoding/json"
    "html/template"
    "net/http"
    "os"

    "github.com/google/uuid"
)

type TemplateData struct {
    Title   string
    Host    string
    OptionA string
    OptionB string
    Vote    string
    ResultA int
    ResultB int
}

var templates *template.Template
var handlerData TemplateData
var hostname, _ = os.Hostname()

func Init(title string, optionA string, optionB string, filenames ...string) error {
    handlerData = TemplateData{
        Title:   title,
        Host:    hostname,
        OptionA: optionA,
        OptionB: optionB,
    }
    files, err := template.ParseFiles(filenames...)
    if err != nil {
        return err
    }
    templates = template.Must(files, nil)
    return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
    err := templates.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func resultsAddDefaults(results map[string]int) {
    for _, k := range []string{"a", "b"} {
        if _, ok := results[k]; !ok {
            results[k] = 0
        }
    }
}

func ResultsHandler(w http.ResponseWriter, r *http.Request, con DatabaseCon) {
    results, err := con.Count("$vote")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    resultsAddDefaults(results)

    data := handlerData
    data.ResultA = results["a"]
    data.ResultB = results["b"]

    renderTemplate(w, "results.html", data)
}

func ApiResultsHandler(w http.ResponseWriter, r *http.Request, con DatabaseCon) {
    results, err := con.Count("$vote")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    resultsAddDefaults(results)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func VoteHandler(w http.ResponseWriter, r *http.Request, con RedisCon) {
    voterId, err := r.Cookie("voterId")
    if err != nil {
        newUuid := uuid.Must(uuid.NewRandom())
        voterId = &http.Cookie{Name:"voterId", Value:newUuid.String()}
    }

    var vote string
    if r.Method == http.MethodPost {
        vote = r.FormValue("vote")
        // build json data
        msg, _ := json.Marshal(&VoteMessage{VoterId: voterId.Value, Vote: vote})
        con.Push("votes", msg)
    }

    http.SetCookie(w, voterId)

    data := handlerData
    data.Vote = vote
    renderTemplate(w, "vote.html", data)
}

func MakeVoteHandler(fn func(w http.ResponseWriter, r *http.Request, con RedisCon), con RedisCon) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fn(w, r, con)
    }
}

func MakeResultsHandler(fn func(w http.ResponseWriter, r *http.Request, con DatabaseCon), con DatabaseCon) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fn(w, r, con)
    }
}

