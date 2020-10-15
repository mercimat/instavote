package core

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
)

var title string = "instavote"
var optA string = "dogs"
var optB string = "cats"
var htmlFiles []string = []string{"../templates/vote.html", "../templates/results.html"}


func init() {
    Init("instavote", "dogs", "cats", htmlFiles...)
}


func TestVoteHandler(t *testing.T) {

    for _, tc := range voteHandlerTestCases {

        t.Logf("Test case: %s", tc.description)

        redisCon := &MockRedisCon{}

        form := url.Values{}
        form.Add("vote", tc.option)

        req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
        if err != nil {
            t.Fatal(err)
        }
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        if tc.voterid != "" {
            req.AddCookie(&http.Cookie{Name:"voterId", Value:tc.voterid})
        }

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(MakeVoteHandler(VoteHandler, redisCon))
        handler.ServeHTTP(rr, req)

        // Check that redis message is pushed to the "votes" queue
        if redisCon.Key != "votes" {
            t.Fatalf("Redis push executed with the wrong key: got %v want %v\n",
                redisCon.Key, "vote")
        }

        // Check that the msg contains a non-empty voter-id
        if redisCon.Msg.VoterId == "" {
            t.Fatalf("VoterId not set when pushing message to Redis.\n")
        }

        // If the voterId cookie is found in the request, it should be set in the redis message
        if tc.voterid != "" && redisCon.Msg.VoterId != tc.voterid {
            t.Fatalf("Redis push - message with incorrect VoterId: got %v want %v\n",
                redisCon.Msg.VoterId, tc.voterid)
        }

        // The redis message contains the specified option
        if redisCon.Msg.Vote != tc.option {
            t.Fatalf("Incorrect vote value in redis push message: got %v want %v",
                redisCon.Msg.Vote, tc.option)
        }

        // The voterId cookie is set in the response headers
        if _, ok := rr.HeaderMap["Set-Cookie"]; !ok {
            t.Fatalf("Cookie not set in the response headers: %v", rr.HeaderMap)
        }
        cookies := rr.HeaderMap["Set-Cookie"]
        expected_cookie := fmt.Sprintf("voterId=%s", redisCon.Msg.VoterId)
        if cookies[0] != expected_cookie {
            t.Fatalf("Unexpected cookie found in the response headers: got %v want %v\n",
                cookies[0], expected_cookie)
        }

        t.Logf("PASS: %s", tc.description)
    }
}


func TestResultsHandler(t *testing.T) {

    for _, tc := range resultsHandlerTestCases {

        t.Logf("Test case: %s", tc.description)

        mongoCon := &MockDatabaseCon{
            CountA: tc.a,
            CountB: tc.b,
            Err: tc.err,
        }

        req, err := http.NewRequest("GET", "/", nil)
        if err != nil {
            t.Fatal(err)
        }

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(MakeResultsHandler(ResultsHandler, mongoCon))
        handler.ServeHTTP(rr, req)

        // Check the status code of the response
        if rr.Code != tc.status {
            t.Fatalf("handler returned wrong status code: got %v want %v",
                rr.Code, tc.status)
        }

        // Check the content of the response, correct results or error message
        if body := rr.Body.String(); !strings.Contains(body, tc.content) {
            t.Fatalf("handler returned wrong body: \nbody does not contain: %s\nGot %s",
                tc.content, body)
        }

        t.Logf("PASS: %s", tc.description)
    }
}


func TestApiResultsHandler(t *testing.T) {

    for _, tc := range apiResultsHandlerTestCases {

        t.Logf("Test case: %s", tc.description)

        mongoCon := &MockDatabaseCon{
            CountA: tc.a,
            CountB: tc.b,
            Err: tc.err,
        }

        req, err := http.NewRequest("GET", "/", nil)
        if err != nil {
            t.Fatal(err)
        }

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(MakeResultsHandler(ApiResultsHandler, mongoCon))
        handler.ServeHTTP(rr, req)

        // Check the status code of the response
        if rr.Code != tc.status {
            t.Fatalf("handler returned wrong status code: got %v want %v",
                rr.Code, tc.status)
        }

        // Check that the Content-Type is set in the response headers
        if _, ok := rr.HeaderMap["Content-Type"]; !ok {
            t.Fatalf("Content-Type not found in the response headers: %v\n", rr.HeaderMap)
        }

        // Check that the Content-Type is set to application/json in the response headers
        if rr.HeaderMap["Content-Type"][0] != "application/json" {
            t.Fatalf("Wrong Content-Type found in the response headers: got %v want %v\n",
            rr.HeaderMap, "application/json")
        }

        // Check the content of the response, correct results or error message
        if body := rr.Body.String(); !strings.Contains(body, tc.content) {
            t.Fatalf("handler returned wrong body: \nbody does not contain: %s\nGot %s",
                tc.content, body)
        }

        t.Logf("PASS: %s", tc.description)
    }
}
