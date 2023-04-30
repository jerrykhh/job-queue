package server

import (
	"context"

	pb "github.com/jerrykhh/job-queue/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Start(ctx context.Context, req *pb.JobQueueRequest) (*pb.JobQueue, error) {
	queueId := req.GetQueueId()

	if q, ok := server.queues[queueId]; ok {
		go q.Start()
		return q.ToPB(), nil
	} else {
		return nil, status.Error(codes.NotFound, "queue id not found")
	}
}

func (server *Server) Create(ctx context.Context, req *pb.JobQueue) (*pb.CreateJobQueueResponse, error) {
	q, err := server.NewJobQueue(req.GetName(), int(req.GetRunEverySec()), int(req.GetRunEverySec()), int(req.GetDequeueCount()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create the new job queue")
	}

	return &pb.CreateJobQueueResponse{
		QueueId:  q.Id,
		JobQueue: q.ToPB(),
	}, nil

}
