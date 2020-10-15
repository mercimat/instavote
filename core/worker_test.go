package core

import (
    "strings"
    "testing"
)

func TestWorker(t *testing.T) {

    for _, tc := range workerTestCases {

        t.Logf("Test case: %s", tc.description)

        message := VoteMessage{tc.voterid, tc.vote}
        redisCon := MockRedisCon{Key:"votes", Msg: message, Vote: tc.vote, Err: tc.err}
        mongoCon := MockDatabaseCon{}

        res := Redis2Mongo(&redisCon, &mongoCon)

        if res != tc.err {
            t.Fatalf("Unexpected result for Redis2Mongo: got %v want %v", res, tc.err)
        }

        if tc.err != nil {
            continue
        }

        if mongoCon.Msg.VoterId != message.VoterId {
            t.Fatalf("Wrong voter-id in message: got %v want %v", mongoCon.Msg.VoterId, message.VoterId)
        }
        if mongoCon.Msg.Vote != message.Vote {
            t.Fatalf("Wrong vote in message: got %v want %v", mongoCon.Msg.Vote, message.Vote)
        }
        if strings.Join(mongoCon.Filter, ",") != strings.Join([]string{"voterid", message.VoterId}, ",") {
            t.Fatalf("Wrong mongodb filter: got %v want %v", mongoCon.Filter, []string{"voterid", message.VoterId})
        }
        if strings.Join(mongoCon.Update, ",") != strings.Join([]string{"vote", message.Vote}, ",") {
            t.Fatalf("Wrong mongodb update: got %v want %v", mongoCon.Update, []string{"vote", message.Vote})
        }
    }
}

