package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/stretchr/testify/suite"
)

type JobQueueSuite struct {
	suite.Suite
	jobClient pb.JobQueueServiceClient
}

func TestJobQueueService(t *testing.T) {
	suite.Run(t, new(JobQueueSuite))
}

func (s *JobQueueSuite) SetupSuite() {
	s.jobClient = ConnJobqueueServer()
}

func (s *JobQueueSuite) TearDownTest() {
	fmt.Println("tear down")
	ctx := context.Background()
	res, _ := s.jobClient.List(ctx, &pb.EmptyRequest{})
	for _, q := range res.Items {
		listJobRes, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: q.Id})
		if err == nil {
			for _, job := range listJobRes.Items {
				s.jobClient.RemoveJob(ctx, &pb.RemoveJobRequest{QueueId: q.Id, Job: job})
			}
		}
		s.jobClient.Remove(ctx, &pb.JobQueueRequest{QueueId: q.Id})
	}
}

func (s *JobQueueSuite) TestCreateQueue() {
	ctx := context.Background()
	fmt.Println("here1")
	request := &pb.CreateJobQueueRequest{
		Name:         "testqueue",
		RunEverySec:  1,
		Seed:         int32Ptr(-1),
		DequeueCount: int32Ptr(3),
	}

	expected := &pb.JobQueue{
		Name:         request.Name,
		RunEverySec:  request.RunEverySec,
		Seed:         request.Seed,
		DequeueCount: request.DequeueCount,
	}
	fmt.Println("here11")
	res, err := s.jobClient.Create(ctx, request)
	fmt.Println("here12")
	s.Nil(err)
	s.NotEmpty(res.Id)
	expected.Id = res.Id
	fmt.Println("here13å")
	s.Equal(expected.String(), res.String())
}

func (s *JobQueueSuite) TestRemoveQueue() {
	ctx := context.Background()

	request := &pb.CreateJobQueueRequest{
		Name:         "testRemove",
		RunEverySec:  1,
		Seed:         int32Ptr(-1),
		DequeueCount: int32Ptr(3),
	}

	expected := &pb.JobQueue{
		Name:         request.Name,
		RunEverySec:  request.RunEverySec,
		Seed:         request.Seed,
		DequeueCount: request.DequeueCount,
	}

	res, err := s.jobClient.Create(ctx, request)
	s.Nil(err)

	expected.Id = res.Id

	res, err = s.jobClient.Remove(ctx, &pb.JobQueueRequest{
		QueueId: res.Id,
	})

	s.Nil(err)
	s.Equal(expected.String(), res.String())
}

func (s *JobQueueSuite) TestRemoveQueueInvalidId() {
	ctx := context.Background()
	_, err := s.jobClient.Remove(ctx, &pb.JobQueueRequest{
		QueueId: "i go to school by bus",
	})
	s.Error(err)
}
