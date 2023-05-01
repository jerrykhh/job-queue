package server_queue

import (
	"github.com/google/uuid"
	"github.com/jerrykhh/job-queue/grpc/pb"
)

type Job struct {
	Id       string
	Script   string
	Parma    string
	Priority int8
}

func NewJob(script, param string) (*Job, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Job{
		Id:       id.String(),
		Script:   script,
		Parma:    param,
		Priority: 1,
	}, nil
}

func (job *Job) ToPB() *pb.Job {
	return &pb.Job{
		Id:       job.Id,
		Script:   job.Script,
		Parma:    job.Parma,
		Priority: int32Ptr(int(job.Priority)),
	}
}
