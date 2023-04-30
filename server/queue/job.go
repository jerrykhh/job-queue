package server_queue

type Job struct {
	Script   string
	Parma    string
	Priority int8
}

func NewJob(script, param string) (*Job, error) {
	return &Job{
		Script:   script,
		Parma:    param,
		Priority: 1,
	}, nil
}
