package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/stretchr/testify/suite"
)

type JobQueueEnqueueSuite struct {
	suite.Suite
	jobClient pb.JobQueueServiceClient
	queue     *pb.JobQueue
}

func TestJobQueueEnqueueService(t *testing.T) {
	suite.Run(t, new(JobQueueEnqueueSuite))
}

func (s *JobQueueEnqueueSuite) SetupSuite() {
	s.jobClient = ConnJobqueueServer()
	ctx := context.Background()
	res, err := s.jobClient.Create(ctx, &pb.CreateJobQueueRequest{
		Name:         "testEnDequeue",
		RunEverySec:  1,
		Seed:         int32Ptr(-1),
		DequeueCount: int32Ptr(3),
	})

	if err != nil {
		log.Fatalln("SetupSuite: Create Queue failed")
	}
	s.queue = res
}

func (s *JobQueueEnqueueSuite) BeforeTest(suiteName, testName string) {
	ctx := context.Background()
	res, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})

	if err != nil {
		log.Fatalln("BeforeTest: List jobs failed", err)
	}

	for _, job := range res.Items {
		_, err := s.jobClient.RemoveJob(ctx, &pb.RemoveJobRequest{
			QueueId: s.queue.Id,
			Job:     job,
		})
		if err != nil {
			log.Fatalln("AfterTest: Remove Jobs failed,", err)
		}
	}

}

func (s *JobQueueEnqueueSuite) TearDownSuite() {
	ctx := context.Background()
	s.jobClient.Remove(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
}

func (s *JobQueueEnqueueSuite) TestEnqueue() {

	ctx := context.Background()

	request := &pb.EnqueueRequest{
		QueueId: s.queue.Id,
		Script:  "test.py",
		Parma:   "args",
	}

	expected := &pb.Job{
		Script:   request.Script,
		Parma:    request.Parma,
		Priority: int32Ptr(0),
	}
	res, err := s.jobClient.Enqueue(ctx, request)
	expected.Id = res.Id
	s.NoError(err)
	s.Equal(expected.String(), res.String())
}

func (s *JobQueueEnqueueSuite) TestEnqueueWithPriority() {
	ctx := context.Background()

	jobs := make([]*pb.Job, 5)
	expectedIndex := [5]int{3, 4, 0, 1, 2}

	jobs[0] = &pb.Job{
		Priority: int32Ptr(1),
		Script:   "test3",
		Parma:    "args",
	}
	jobs[1] = &pb.Job{
		Priority: int32Ptr(2),
		Script:   "test4",
		Parma:    "args",
	}
	jobs[2] = &pb.Job{
		Priority: int32Ptr(3),
		Script:   "test5",
		Parma:    "args",
	}
	jobs[3] = &pb.Job{
		Priority: int32Ptr(-1),
		Script:   "test1",
		Parma:    "args",
	}
	jobs[4] = &pb.Job{
		Priority: int32Ptr(-1),
		Script:   "test2",
		Parma:    "args",
	}

	for _, job := range jobs {
		res, err := s.jobClient.Enqueue(ctx, &pb.EnqueueRequest{
			QueueId:  s.queue.Id,
			Script:   job.Script,
			Parma:    job.Parma,
			Priority: job.Priority,
		})
		s.NoError(err)
		job.Id = res.Id
	}

	res, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})

	if s.NoError(err) {
		for i, job := range res.Items {
			fmt.Println(jobs[expectedIndex[i]].Script, job.Script)
			s.Equal(jobs[expectedIndex[i]].Script, job.Script)
		}
	}

}
