package types

import "time"

type Song interface {
	GetName() string
	GetDuration() time.Duration
	GetFilePath() string
}
