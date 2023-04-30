package server_queue

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jerrykhh/job-queue/grpc/pb"
)

type JobQueue struct {
	pb.JobQueue
	Id              string
	Name            string
	RunEvery        time.Duration
	Seed            int
	DequeueCount    int
	PerviousRunTime time.Time
	Pause           bool
}

func NewJobQueue(name string, runEverySec, seed, dequeueCount int) (*JobQueue, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed create the job queueu")
	}

	return &JobQueue{
		Id:           id.String(),
		Name:         name,
		RunEvery:     time.Duration(runEverySec),
		Seed:         seed,
		DequeueCount: dequeueCount,
		Pause:        false,
	}, nil
}

func (jobQueue *JobQueue) ToPB() *pb.JobQueue {
	return &pb.JobQueue{
		Name:         jobQueue.Name,
		RunEverySec:  int32(jobQueue.RunEvery.Seconds()),
		Seed:         int32Ptr(jobQueue.Seed),
		DequeueCount: int32Ptr(jobQueue.DequeueCount),
	}
}

func int32Ptr(v int) *int32 {
	p := int32(v)
	return &p
}

func (queue *JobQueue) Start() {

	// queue.Pause = false

	// for {
	// 	if queue.Pause {
	// 		break
	// 	}

	// 	jobs, err := queue.Dequeue()
	// 	if err != nil {
	// 		log.Println("Error dequeuing:", err, queue.Name)
	// 	} else {
	// 		for _, cron := range jobs {
	// 			log.Printf("Running %s (%s)\n", cron.TargetCmd, q.Name)
	// 		}
	// 		queue.PerviousRunTime = time.Now()
	// 	}
	// 	time.Sleep(queue.RunEvery * time.Second)

	// }
}

// func (queue * JobQueue) Dequeue() ([]* Job, error) {

// }
