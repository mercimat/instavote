package core

import (
    "errors"
    "net/http"
)

var voteHandlerTestCases = []struct{
    description string
    voterid     string
    option      string
}{
    {
        description:    "Vote \"a\" and no voter-id cookie",
        voterid:        "",
        option:         "a",
    },
    {
        description:    "Vote \"b\" and no voter-id cookie",
        voterid:        "",
        option:         "b",
    },
    {
        description:    "Vote \"a\" and voter-id cookie",
        voterid:        "3f1ed0b1-3e31-4c39-8a00-84ba26e70a6b",
        option:         "a",
    },
    {
        description:    "Vote \"b\" and voter-id cookie",
        voterid:        "3f1ed0b1-3e31-4c39-8a00-84ba26e70a6b",
        option:         "b",
    },
}

var resultsHandlerTestCases = []struct{
    description string
    a           int
    b           int
    err         error
    status      int
    content     string
}{
    {
        description:    "Results with no error",
        a:              10,
        b:              1000,
        err:            nil,
        status:         http.StatusOK,
        content:        `
        <div class="results">
            <div class="result">
                <h2 id="a">10</h2>
            </div>

            <div class="result">
                <h2 id="b">1000</h2>
            </div>
        </div>`,
    },
    {
        description:    "Results with error",
        a:              10,
        b:              1000,
        err:            errors.New("failed to get results"),
        status:         http.StatusInternalServerError,
        content:        "failed to get results",
    },
}

var apiResultsHandlerTestCases = []struct{
    description string
    a           int
    b           int
    err         error
    status      int
    content     string
}{
    {
        description:    "API results with no error",
        a:              10,
        b:              1000,
        err:            nil,
        status:         http.StatusOK,
        content:        `{"a":10,"b":1000}`,
    },
    {
        description:    "API results with error",
        a:              10,
        b:              1000,
        err:            errors.New("failed to get results"),
        status:         http.StatusInternalServerError,
        content:        "failed to get results",
    },
}
