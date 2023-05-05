package test

import (
	"context"
	"log"
	"math"
	"testing"
	"time"

	"github.com/jerrykhh/job-queue/grpc/pb"
	"github.com/stretchr/testify/suite"
)

type JobQueueStartSuite struct {
	suite.Suite
	jobClient pb.JobQueueServiceClient
	queue     *pb.JobQueue
	jobs      []*pb.Job
}

func TestJobQueueStartService(t *testing.T) {
	suite.Run(t, new(JobQueueStartSuite))
}

func (s *JobQueueStartSuite) SetupSuite() {
	s.jobClient = ConnJobqueueServer()
	ctx := context.Background()

	res, err := s.jobClient.Create(ctx, &pb.CreateJobQueueRequest{
		Name:         "testJobStart",
		RunEverySec:  2,
		Seed:         int32Ptr(-1),
		DequeueCount: int32Ptr(2),
	})
	if err != nil {
		log.Fatalln("SetupSuite: Create Queue failed")
	}
	s.queue = res
}

func (s *JobQueueStartSuite) BeforeTest(suitName, testName string) {
	ctx := context.Background()
	jobs := make([]*pb.Job, 5)
	jobs[0] = &pb.Job{
		Priority: int32Ptr(1),
		Script:   "test3",
		Parma:    "args",
	}
	jobs[1] = &pb.Job{
		Priority: int32Ptr(-2),
		Script:   "test4",
		Parma:    "args",
	}
	jobs[2] = &pb.Job{
		Priority: int32Ptr(2),
		Script:   "test5",
		Parma:    "args",
	}
	jobs[3] = &pb.Job{
		Priority: int32Ptr(-1),
		Script:   "test1",
		Parma:    "args",
	}
	jobs[4] = &pb.Job{
		Priority: int32Ptr(-19),
		Script:   "test1",
		Parma:    "args",
	}

	for _, job := range jobs {
		res, _ := s.jobClient.Enqueue(ctx, &pb.EnqueueRequest{
			QueueId:  s.queue.Id,
			Script:   job.Script,
			Priority: job.Priority,
			Parma:    job.Parma,
		})
		job.Id = res.Id
	}
	s.jobs = jobs

}

func (s *JobQueueStartSuite) AfterTest(suitName, testName string) {
	ctx := context.Background()

	res, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	if err != nil {
		log.Fatalln("BeforeTest: ListJob failed")
	}
	for _, job := range res.Items {
		s.jobClient.RemoveJob(ctx, &pb.RemoveJobRequest{QueueId: s.queue.Id, Job: job})
	}

}

func (s *JobQueueStartSuite) TearDownSuite() {
	ctx := context.Background()
	s.jobClient.Remove(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
}

func (s *JobQueueStartSuite) TestQueueStart() {

	ctx := context.Background()

	expectedListJobsIndex := [3]int{3, 0, 2}

	res, err := s.jobClient.Start(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})

	s.NoError(err)
	s.Equal(s.queue.String(), res.String())
	time.Sleep(1 * time.Second)
	listJobsRes, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	s.NoError(err)
	for i, job := range listJobsRes.Items {
		s.Equal(s.jobs[expectedListJobsIndex[i]].String(), job.String())
	}

	var countTime int = int(math.Ceil(float64(len(s.jobs)) / float64(*s.queue.DequeueCount)))
	time.Sleep(time.Duration(int(s.queue.RunEverySec) * countTime * int(time.Second)))
	listJobsRes, err = s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	s.NoError(err)
	s.Equal(0, len(listJobsRes.Items))

}

func (s *JobQueueStartSuite) TestQueueStartAndPause() {

	ctx := context.Background()

	res, err := s.jobClient.Start(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})

	s.NoError(err)
	s.Equal(s.queue.String(), res.String())
	time.Sleep(1 * time.Second)
	pauseRes, err := s.jobClient.Pause(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	s.NoError(err)
	s.Equal(s.queue.String(), pauseRes.String())

	listJobsRes, err := s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	s.NoError(err)
	var expected int = len(listJobsRes.Items)

	var countTime int = int(math.Ceil(float64(len(s.jobs)) / float64(*s.queue.DequeueCount)))
	time.Sleep(time.Duration(int(s.queue.RunEverySec) * (countTime + 1) * int(time.Second)))
	listJobsRes, err = s.jobClient.ListJob(ctx, &pb.JobQueueRequest{QueueId: s.queue.Id})
	s.NoError(err)
	s.Equal(expected, len(listJobsRes.Items))

}

func (s *JobQueueStartSuite) TestRemoveJob() {
	ctx := context.Background()
	for _, job := range s.jobs {
		res, err := s.jobClient.RemoveJob(ctx, &pb.RemoveJobRequest{QueueId: s.queue.Id, Job: job})
		s.NoError(err)
		s.Equal(job.String(), res.String())
	}
}
