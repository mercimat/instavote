package e2e

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "os"
    "testing"
    "time"
)

type Results struct {
    A int `json:"a"`
    B int `json:"b"`
}

var voteUrl string
var resultsUrl string

func init() {
    voteUrl = "http://localhost:8090/"
    resultsUrl = "http://localhost:8091/api/results"
    if _, ok := os.LookupEnv("JENKINS_HOST"); ok {
        voteUrl = os.ExpandEnv("http://${JENKINS_HOST}:8090/")
        resultsUrl = os.ExpandEnv("http://${JENKINS_HOST}:8091/api/results")
    }
}

func TestEnd2End(t *testing.T) {
    // Initial results should be 0, as docker-compose is configured to clear the volume associated with MongoDB
    t.Logf("Check initial results")
    checkResults(t, Results{})


    t.Logf("Add a vote for option \"a\"")
    _, err := http.PostForm(voteUrl, url.Values{"vote": {"a"}})
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(1000 * time.Millisecond)
    checkResults(t, Results{A:1, B:0})


    t.Logf("Add a new vote for option \"b\"")
    // First, create a http client to handle cookies
    jar, err := cookiejar.New(nil)
    if err != nil {
        t.Fatal(err)
    }
    client := &http.Client{
        Jar: jar,
    }

    // Then add a vote for b
    _, err = client.PostForm(voteUrl, url.Values{"vote": {"b"}})
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(1000 * time.Millisecond)
    checkResults(t, Results{A:1, B:1})

    t.Logf("Change the previous vote for option \"a\"")
    // Using client will allow to re-use the cookie that was set with the previous response
    _, err = client.PostForm(voteUrl, url.Values{"vote": {"a"}})
    if err != nil {
        t.Fatal(err)
    }
    time.Sleep(1000 * time.Millisecond)
    checkResults(t, Results{A:2, B:0})
}

func checkResults(t *testing.T, expected Results) {
    results, err := getApiResults()
    if err != nil {
        t.Fatal(err)
    }
    if results.A != expected.A || results.B != expected.B {
        t.Fatalf("Unexpected values for results: got %v but expected %v", results, expected)
    }
}

func getApiResults() (r Results, _ error) {
    resp, err := http.Get(resultsUrl)
    if err != nil {
        return r, err
    }

    body, _ := ioutil.ReadAll(resp.Body)
    err = json.Unmarshal(body, &r)
    if err != nil {
        return r, err
    }
    return r, nil
}
