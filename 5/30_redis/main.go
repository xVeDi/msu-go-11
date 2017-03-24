package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	c redis.Conn
)

func getRecord(mkey string) string {
	println("get", mkey)
	// получает запись, https://redis.io/commands/get
	item, err := redis.String(c.Do("GET", mkey))
	// если записи нету, то для этого есть специальная ошибка, её надо обрабатывать отдеьно, это почти штатная ситуация, а не что-то страшное
	if err == redis.ErrNil {
		fmt.Println("Record not found in redis (return value is nil)")
		return ""
	} else if err != nil {
		PanicOnErr(err)
	}
	return item
}

func main() {
	var err error
	// соединение
	c, err = redis.DialURL("redis://user:@localhost:6379/0")
	PanicOnErr(err)
	defer c.Close()

	mkey := "record_21"

	item := getRecord(mkey)
	fmt.Printf("first get %+v\n", item)

	ttl := 5
	// добавляет запись, https://redis.io/commands/set
	result, err := redis.String(c.Do("SET", mkey, 1, "EX", ttl))
	PanicOnErr(err)
	if result != "OK" {
		panic("result not ok: " + result)
	}

	time.Sleep(time.Microsecond)
	// time.Sleep( (ttl+1)*time.Second)

	item = getRecord(mkey)
	fmt.Printf("second get %+v\n", item)

	// https://redis.io/commands/incrby
	// https://redis.io/commands/incr - тут хорошо описан рейтилимитер
	n, _ := redis.Int(c.Do("INCRBY", mkey, 2))
	fmt.Println("INCRBY by 2 ", mkey, "is", n)

	// https://redis.io/commands/decrby
	n, _ = redis.Int(c.Do("DECRBY", mkey, 1))
	fmt.Println("DECRY by 1 ", mkey, "is", n)

	// если записи не было - редис создаст
	n, err = redis.Int(c.Do("INCR", mkey+"_not_exist"))
	fmt.Println("INCR (default by 1) ", mkey+"_not_exist", "is", n)
	PanicOnErr(err)

	keys := []interface{}{mkey, mkey + "_not_exist", "sure_not_exist"}

	reply, err := redis.Strings(c.Do("MGET", keys...))
	PanicOnErr(err)
	fmt.Println(reply)
}

//PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
