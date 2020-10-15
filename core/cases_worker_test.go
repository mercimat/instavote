package core

import (
    "errors"
)

var workerTestCases = []struct{
    description string
    voterid     string
    vote        string
    err         error
}{
    {
        description:    "From Redis to Mongodb",
        voterid:        "id1",
        vote:           "a",
        err:            nil,
    },
    {
        description:    "Error on Redis Get",
        voterid:        "id1",
        vote:           "b",
        err:            errors.New("redis get failed"),
    },
}


