package server

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/jerrykhh/job-queue/grpc/pb"
	server_queue "github.com/jerrykhh/job-queue/server/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Start(ctx context.Context, req *pb.JobQueueRequest) (*pb.JobQueue, error) {
	q, err := server.GetJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}
	go q.Start(-1)
	return q.ToPB(), nil

}

func (server *Server) ListJob(req *pb.JobQueueRequest, stream pb.JobQueueService_ListJobServer) error {

	q, err := server.GetJobQueue(req.GetQueueId())

	if err != nil {
		return status.Error(codes.NotFound, "queue id not found")
	}

	results, err := q.List()
	if err != nil {
		return status.Error(codes.Internal, "failed to get List Queue")
	}

	for _, result := range results {
		jobJSON := result.Member.(string)
		var job server_queue.Job

		err := json.Unmarshal([]byte(jobJSON), &job)
		if err != nil {
			log.Println(err)
		}

		if err := stream.Send(job.ToPB()); err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) Create(ctx context.Context, req *pb.CreateJobQueueRequest) (*pb.JobQueue, error) {
	q, err := server.NewJobQueue(req.GetName(), int(req.GetRunEverySec()), int(req.GetSeed()), int(req.GetDequeueCount()))

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create the new job queue")
	}

	return q.ToPB(), nil
}

func (server *Server) Enqueue(ctx context.Context, req *pb.EnqueueRequest) (*pb.Job, error) {
	q, err := server.GetJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}

	job, err := server_queue.NewJob(req.GetScript(), req.GetParma())
	job.Priority = int8(req.GetPriority())

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create new job")
	}

	err = q.Enqueue(job)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to enqueue new job")
	}

	return job.ToPB(), nil

}

func (server *Server) Dequeue(ctx context.Context, req *pb.JobQueueRequest) (*pb.DequeueResponse, error) {
	q, err := server.GetJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}

	jobs, err := q.Dequeue()
	if err != nil {
		return nil, status.Error(codes.Aborted, "Queue is empty")
	}

	pbJobs := make([]*pb.Job, len(jobs))
	for i, job := range jobs {
		pbJobs[i] = job.ToPB()
	}

	return &pb.DequeueResponse{
		Items: pbJobs,
	}, nil

}

func (server *Server) Pause(ctx context.Context, req *pb.JobQueueRequest) (*pb.JobQueue, error) {
	q, err := server.GetJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}

	q.Pause = true
	return q.ToPB(), nil
}

func (server *Server) Remove(ctx context.Context, req *pb.JobQueueRequest) (*pb.JobQueue, error) {
	q, err := server.RemoveJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}

	return q.ToPB(), nil
}

func (server *Server) RemoveJob(ctx context.Context, req *pb.RemoveJobRequest) (*pb.Job, error) {
	q, err := server.RemoveJobQueue(req.GetQueueId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}

	removeJob := &server_queue.Job{
		Id:       req.GetJob().GetId(),
		Script:   req.Job.GetScript(),
		Parma:    req.GetJob().GetParma(),
		Priority: int8(req.GetJob().GetPriority()),
	}

	err = q.RemoveJob(removeJob)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to remove the job from queuue")
	}
	return req.GetJob(), nil

}

func (server *Server) List(ctx context.Context, req *pb.EmptyRequest) (*pb.ListRepsonse, error) {

	qs, err := server.ListJobQueue()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list queue")
	}

	pbQueues := make([]*pb.JobQueue, len(qs))
	for i, queue := range qs {
		pbQueues[i] = queue.ToPB()
	}

	return &pb.ListRepsonse{
		Items: pbQueues,
	}, nil
}
