package lib

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisPool *redis.Pool

var (
	redisServer   = "127.0.0.1:6379"
	redisPassword = ""
)

func init() {
	redisPool = newPool(redisServer, redisPassword)
}

func Redis() redis.Conn {
	return redisPool.Get()
}

func Get(key string) (string, error) {
	reply, err := Redis().Do("GET", key)
	return redis.String(reply, err)
}

func ZADD(key string, score int, field string) (int, error) {
	reply, err := Redis().Do("ZADD", key, score, field)
	return redis.Int(reply, err)
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			//Select db
			if _, err := c.Do("SELECT", 1); err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
