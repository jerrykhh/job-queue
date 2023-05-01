package server_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/redis/go-redis/v9"
)

type JobQueue struct {
	pb.JobQueue
	Id              string
	Name            string
	RunEvery        time.Duration
	Seed            int
	DequeueCount    int
	count           int
	PerviousRunTime time.Time
	Pause           bool
	Redis           *redis.Client
}

func NewJobQueue(name string, runEverySec, seed, dequeueCount int, redisAddr string, redisPort int) (*JobQueue, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed create the job queueu")
	}

	jobQueue := &JobQueue{
		Id:           id.String(),
		Name:         name,
		RunEvery:     time.Duration(runEverySec),
		Seed:         seed,
		DequeueCount: dequeueCount,
		Pause:        false,
		count:        0,
	}

	jobQueue.Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisAddr, redisPort),
	})
	return jobQueue, nil
}

func (jobQueue *JobQueue) ToPB() *pb.JobQueue {
	return &pb.JobQueue{
		Name:         jobQueue.Name,
		RunEverySec:  int32(jobQueue.RunEvery.Seconds()),
		Seed:         int32Ptr(jobQueue.Seed),
		DequeueCount: int32Ptr(jobQueue.DequeueCount),
	}
}

func (queue *JobQueue) IsEmpty() bool {
	return queue.count == 0
}

func (queue *JobQueue) Start(seed int64) {

	queue.Pause = false

	var bias int64 = 0
	rand.Seed(seed)
	if queue.Seed != -1 {
		bias = rand.Int63()
	}

	fmt.Println("bias", bias)

	for {
		if queue.Pause {
			break
		}

		if !queue.IsEmpty() {
			jobs, err := queue.Dequeue()
			if err != nil {
				fmt.Println("Error dequeuing:", err, queue.Id)
			} else {
				for _, job := range jobs {
					fmt.Println(job)
					fmt.Printf("Running %s (%s)\n", job.Script, queue.Id)
				}
				queue.PerviousRunTime = time.Now()
			}
		}
		fmt.Println((queue.RunEvery + time.Duration(bias)))
		time.Sleep((queue.RunEvery + time.Duration(bias)) * time.Second)

	}
}

func (queue *JobQueue) Enqueue(job *Job) error {
	jobJson, err := json.Marshal(job)

	if err != nil {
		return err
	}

	node := redis.Z{
		Score:  float64(job.Priority),
		Member: string(jobJson),
	}

	err = queue.Redis.ZAdd(context.Background(), queue.Id, node).Err()
	if err != nil {
		return err
	}

	queue.count++
	return nil
}

func (queue *JobQueue) Dequeue() ([]*Job, error) {

	if queue.count == 0 {
		return nil, fmt.Errorf("Queue(%s) is empty", queue.Id)
	}

	ctx := context.Background()
	result, err := queue.Redis.ZRangeByScoreWithScores(ctx, queue.Id, &redis.ZRangeBy{
		Min:   "-20",
		Max:   "19",
		Count: int64(queue.DequeueCount),
	}).Result()

	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, len(result))
	for i, data := range result {
		jobJson := data.Member.(string)
		var job Job
		err := json.Unmarshal([]byte(jobJson), &job)
		if err != nil {
			log.Println(err)
		}

		jobs[i] = &job
		_, err = queue.Redis.ZRem(ctx, queue.Id, jobJson).Result()
		if err != nil {
			log.Println(err)
		}
		queue.count--
	}
	return jobs, nil
}

func (queue *JobQueue) List() ([]redis.Z, error) {
	ctx := context.Background()
	fmt.Println("ctx")
	results, err := queue.Redis.ZRangeWithScores(ctx, queue.Id, 0, -1).Result()
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	fmt.Println(results)
	return results, nil
}

func (queue *JobQueue) RemoveJob(job *Job) error {
	ctx := context.Background()

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return err
	}
	queue.count--
	return queue.Redis.ZRem(ctx, queue.Id, string(jobJSON)).Err()
}
