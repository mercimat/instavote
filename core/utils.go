package core

type DatabaseCon interface {
    Count(string) (map[string]int, error)
    Set(interface{}, interface{}, interface{}) error
    QueryFilter(key string, value interface{}) interface{}
    QueryUpdate(key string, value interface{}) interface{}
}

type RedisCon interface {
    Push(key string, msg interface{})
    Get(key string) ([]string, error)
}

type VoteMessage struct {
    VoterId string  `json:"voterid"`
    Vote    string  `json:"vote"`
}

