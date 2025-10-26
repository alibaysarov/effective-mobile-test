package model

import "time"

type Subscription struct {
	ID          string
	ServiceName string
	UserId      string
	StartDate   time.Time
	Price       float64
}
