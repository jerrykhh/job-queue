package server_queue

func int32Ptr(v int) *int32 {
	p := int32(v)
	return &p
}
