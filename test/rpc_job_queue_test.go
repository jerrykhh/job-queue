package test

import (
	"context"
	"log"
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
	ctx := context.Background()
	res, err := s.jobClient.List(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalln("TearDown: List job queue failed")
	}
	for _, q := range res.Items {
		_, err := s.jobClient.Remove(ctx, &pb.JobQueueRequest{
			QueueId: q.Id,
		})
		if err != nil {
			log.Fatalln("TearDown: Remove Queue failed")
		}

	}
}

func (s *JobQueueSuite) TestCreateQueue() {
	ctx := context.Background()

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

	res, err := s.jobClient.Create(ctx, request)
	s.Nil(err)
	s.NotEmpty(res.Id)
	expected.Id = res.Id
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

func (s *JobQueueSuite) TestEqueue() {

	job := &pb.EnqueueRequest{
		Script: "testscript",
		Parma:  "test parma",
	}

	expected := &pb.Job{
		Script: job.Script,
		Parma:  job.Parma,
	}

	ctx := context.Background()
	res, err := s.jobClient.Create(ctx, &pb.CreateJobQueueRequest{
		Name:         "testqueue",
		RunEverySec:  1,
		Seed:         int32Ptr(-1),
		DequeueCount: int32Ptr(3),
	})

	if s.Nil(err) {
		job.QueueId = res.Id
		res, err := s.jobClient.Enqueue(ctx, job)
		expected.Id = res.Id
		expected.Priority = int32Ptr(0)
		s.Nil(err)
		s.Equal(expected.String(), res.String())
	}

}
