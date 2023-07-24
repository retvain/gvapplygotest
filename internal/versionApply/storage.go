package versionApply

type Repository interface {
	FindNewJob() (queue *Queue, err error)
}
