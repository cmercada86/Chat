package cache

import (
	"Chat/model"
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
)

const maxConnections = 10
const expireTime = 10000

var redisPool *redis.Pool

func ConnectToCache(redisAddress string) {
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisAddress)

		if err != nil {
			return nil, err
		}

		return c, err
	}, maxConnections)
}

func Get(key model.UUID) model.User {
	c := redisPool.Get()
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		log.Println("Redis error: ", err)
		return model.User{}
	}
	var user model.User

	json.Unmarshal([]byte(value), &user)

	return user

}

func Set(key model.UUID, value model.User) {
	c := redisPool.Get()
	defer c.Close()

	data, _ := json.Marshal(&value)
	c.Do("SET", key, string(data))
}

func Contains(key model.UUID) bool {
	c := redisPool.Get()
	defer c.Close()

	val, err := redis.Int(c.Do("EXISTS", key))
	if err != nil {
		log.Println("Redis error: ", err)
		return false
	}
	return val > 0

}

func Expire(key model.UUID) {
	c := redisPool.Get()
	defer c.Close()

	c.Do("EXPIRE", key, expireTime)

}

// func Flush() {
// 	c := redisPool.Get()
// 	defer c.Close()

// 	c.Do("FLUSHALL")
// }

func Close() {
	redisPool.Close()
}
