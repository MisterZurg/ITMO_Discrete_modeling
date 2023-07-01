package model

import "time"

type Customer struct {
	ID           int
	PeopleBefore int

	MetricCreatedTime          time.Time
	MetricArrivedAtCashBoxLine time.Time
	MetricArrivedToCashierTime time.Time
	MetricLeftTime             time.Time
}
