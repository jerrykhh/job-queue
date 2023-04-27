package main

import (
	"cron-queue/queue"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	q := map[string]*queue.CronQueue{}
	q["Test1"] = queue.NewCronQueue("Test1", rdb, 1*time.Second, 0, 3)
	// q["Test2"] = queue.NewCronQueue("Test2", rdb, 1*time.Second, 0, 3)
	// q["Test3"] = queue.NewCronQueue("Test3", rdb, 3*time.Second, 0, 3)
	for _, v := range q {
		v.Start()
	}
	fmt.Println("Enqueue")
	q["Test1"].Enqueue(queue.NewCron("test1", 1))
	q["Test1"].Enqueue(queue.NewCron("test2", 1))
	q["Test1"].Enqueue(queue.NewCron("test23", 1))
	q["Test1"].Enqueue(queue.NewCron("test3", 1))
	q["Test1"].Enqueue(queue.NewCron("test4", 1))
	q["Test1"].Enqueue(queue.NewCron("test5", 1))
	q["Test1"].Enqueue(queue.NewCron("test6", -19))
	fmt.Println("Enqueue finish")

	select {}
}
