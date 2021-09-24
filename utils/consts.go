package utils

type JobStatus uint

const (
	Assigned   JobStatus = 1
	Processing JobStatus = 2
	Finished   JobStatus = 3
	Error      JobStatus = 4
)
