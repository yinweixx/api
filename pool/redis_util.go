package pool

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

//Ping --
func Ping(conn redis.Conn) error {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db:%v", err)
	}
	return nil
}

//Get --
func Get(conn redis.Conn, key string) ([]byte, error) {
	// conn.Do("EXPIRE", key, 200*time.Millisecond)
	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

//Set --
func Set(conn redis.Conn, key string, value []byte) error {
	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

//Exists --
func Exists(conn redis.Conn, key string) (bool, error) {
	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

//Delete --
func Delete(conn redis.Conn, key string) error {
	_, err := conn.Do("DEL", key)
	return err
}

//GetKeys --
func GetKeys(conn redis.Conn, pattern string) ([]string, error) {
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

//Incr --
func Incr(conn redis.Conn, counterKey string) (int, error) {
	return redis.Int(conn.Do("INCR", counterKey))
}
