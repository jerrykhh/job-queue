package server_queue

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobQueue struct {
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
