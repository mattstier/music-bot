package filePlayer

import (
	"time"
)

type Song struct {
	name     string
	duration time.Duration
	path     string
}

func (s Song) GetName() string {
	return s.name
}

func (s Song) GetDuration() time.Duration {
	return s.duration
}

func (s Song) GetFilePath() string {
	//TODO implement me
	panic("implement me")
}
