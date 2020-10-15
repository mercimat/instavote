package core

import (
    "encoding/json"
)

type MockDatabaseCon struct {
    CountA int
    CountB int
    Err error
    Filter []string
    Update []string
    Msg VoteMessage
}

func (m *MockDatabaseCon) Count(_ string) (map[string]int, error) {
    res := make(map[string]int)
    if m.Err != nil {
        return res, m.Err
    }
    res["a"] = m.CountA
    res["b"] = m.CountB
    return res, nil
}

func (m *MockDatabaseCon) Set(message, _, _ interface{}) error {
    m.Msg = *(message.(*VoteMessage))
    if  m.Err != nil {
        return m.Err
    }
    return nil
}

func (m *MockDatabaseCon) QueryFilter(key string, value interface{}) interface{} {
    m.Filter = append(m.Filter, key)
    m.Filter = append(m.Filter, value.(string))
    return nil
}

func (m *MockDatabaseCon) QueryUpdate(key string, value interface{}) interface{} {
    m.Update = append(m.Update, key)
    m.Update = append(m.Update, value.(string))
    return nil
}

type MockRedisCon struct {
    Key string
    Msg VoteMessage
    Vote string
    Err error
}

func (r *MockRedisCon) Push(key string, msg interface{}) {
    r.Key = key
    var message VoteMessage
    _ = json.Unmarshal(msg.([]byte), &message)
    r.Msg = message
}

func (r *MockRedisCon) Get(key string) ([]string, error) {
    if r.Err != nil {
        return []string{}, r.Err
    }
    msg, _ := json.Marshal(&r.Msg)
    return []string{r.Key, string(msg)}, nil
}
