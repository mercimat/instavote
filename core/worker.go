package core

import (
    "encoding/json"
)

func Redis2Mongo(redis RedisCon, database DatabaseCon) error {
    vote, err := redis.Get("votes")
    if err != nil {
        return err
    }

    var message VoteMessage
    err = json.Unmarshal([]byte(vote[1]), &message)
    if err != nil {
        return err
    }

    filter := database.QueryFilter("voterid", message.VoterId)
    update := database.QueryUpdate("vote", message.Vote)
    database.Set(&message, filter, update)
    return nil
}
