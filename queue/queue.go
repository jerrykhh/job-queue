package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cron struct {
	TargetCmd string
	Priority  int
	Datetime  time.Time
}

func NewCron(cmd string, priority int) *Cron {
	return &Cron{
		TargetCmd: cmd,
		Priority:  priority,
		Datetime:  time.Now(),
	}
}

type CronQueue struct {
	Name                  string
	Redis                 *redis.Client
	RunEvery              time.Duration
	CreateDatetime        time.Time
	PreviousStartDatetime time.Time
	Pause                 bool
	Seed                  int
	DequeueCount          int
}

func NewCronQueue(name string, redisClient *redis.Client, runEvery time.Duration, seed int, dequeueCount int) *CronQueue {

	return &CronQueue{
		Name:           name,
		Redis:          redisClient,
		RunEvery:       runEvery,
		CreateDatetime: time.Now(),
		Seed:           seed,
		Pause:          false,
		DequeueCount:   dequeueCount,
	}

}

func (q *CronQueue) start() {
	for {

		if q.Pause {
			var count int = 0
			for {
				if !q.Pause {
					break
				}

				if count < 10 {
					time.Sleep(time.Duration(15) * time.Second)
				} else {
					time.Sleep(time.Duration(60) * time.Second)
				}
				count++
			}
		}

		crons, err := q.Dequeue()
		if err != nil {
			log.Println("Error dequeuing:", err, q.Name)
		} else {

			for _, cron := range crons {
				log.Printf("Running %s (%s)\n", cron.TargetCmd, q.Name)
			}
			q.PreviousStartDatetime = time.Now()
		}

		time.Sleep(q.RunEvery)

	}
}

func (q *CronQueue) Start() {
	q.Pause = false
	go q.start()
}

func (q *CronQueue) PauseCronQueue() {
	q.Pause = true
}

func (q *CronQueue) Enqueue(cron *Cron) error {
	ctx := context.Background()

	cronJSON, err := json.Marshal(cron)

	if err != nil {
		return err
	}

	node := redis.Z{
		Score:  float64(cron.Priority),
		Member: string(cronJSON),
	}
	return q.Redis.ZAdd(ctx, q.Name, node).Err()
}

func (q *CronQueue) Dequeue() ([]*Cron, error) {
	ctx := context.Background()
	results, err := q.Redis.ZRangeByScoreWithScores(ctx, q.Name, &redis.ZRangeBy{
		Min:   "-20",
		Max:   "19",
		Count: int64(q.DequeueCount),
	}).Result()

	if err != nil {
		return nil, err
	}

	// convert to Cron
	crons := make([]*Cron, len(results))
	for i, result := range results {
		cronJSON := result.Member.(string)
		var cron Cron

		err := json.Unmarshal([]byte(cronJSON), &cron)
		if err != nil {
			log.Println(err)
		}

		crons[i] = &cron
		_, err = q.Redis.ZRem(ctx, q.Name, cronJSON).Result()
		if err != nil {
			log.Println(err)
		}

	}

	return crons, nil
}
