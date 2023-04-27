package test

import (
	"cron-queue/queue"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type CronQueueSuit struct {
	suite.Suite
	q *queue.CronQueue
}

func (s *CronQueueSuit) SetupSuite() {
	s.q = &queue.CronQueue{
		Name:         "test1",
		Redis:        redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
		RunEvery:     time.Second,
		Seed:         0,
		DequeueCount: 3,
	}
}

func (s *CronQueueSuit) TestEnqueue() {
	crons := make([]*queue.Cron, 3)
	crons[0] = queue.NewCron("testcmd1", 1)
	crons[1] = queue.NewCron("testcmd2", 1)
	crons[2] = queue.NewCron("testcmd3", 1)

	for _, cron := range crons {
		s.q.Enqueue(cron)
	}

	s.q.DequeueCount = 3

	dequeueCrons, err := s.q.Dequeue()
	s.Assert().Equal(nil, err)

	for i, cron := range dequeueCrons {
		s.Assert().Equal(cron.TargetCmd, crons[i].TargetCmd)
	}

}

func (s *CronQueueSuit) TestDequeueWithPriority() {
	crons := make([]*queue.Cron, 3)
	crons[0] = queue.NewCron("testcmd1", 1)
	crons[1] = queue.NewCron("testcmd2", 2)
	crons[2] = queue.NewCron("testcmd3", -3)

	for _, cron := range crons {
		s.q.Enqueue(cron)
	}

	s.q.DequeueCount = 3

	dequeueCrons, err := s.q.Dequeue()
	s.Assert().Equal(nil, err)

	s.Assert().Equal(dequeueCrons[0].TargetCmd, crons[2].TargetCmd)
	s.Assert().Equal(dequeueCrons[1].TargetCmd, crons[0].TargetCmd)
	s.Assert().Equal(dequeueCrons[2].TargetCmd, crons[1].TargetCmd)
}

func TestStart(t *testing.T) {
	suite.Run(t, new(CronQueueSuit))
}
