package redis

import (
    "context"
    "log"
    "time"
    "github.com/go-redis/redis/v8"
)

type RedisCon struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCon(addr string, password string, db int) *RedisCon {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    // Check that it's reachable or exit with Fatal err
    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal(err)
    }

    return &RedisCon{rdb, ctx}
}

func (r *RedisCon) Push(key string, msg interface{}) {
    r.client.RPush(r.ctx, key, msg)
}

// Blocks until there is data to pop
func (r *RedisCon) Get(key string) ([]string, error) {
    data, err := r.client.BLPop(r.ctx, 0*time.Second, key).Result()
    return data, err
}
