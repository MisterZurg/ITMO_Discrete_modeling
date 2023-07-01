package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"time"
)

func MakeCustumers(qSize uint32, queue chan model.Customer, done chan struct{}) {
	for i := 0; uint32(i) < qSize; i++ {
		queue <- model.Customer{
			ID:                         i,
			MetricArrivedAtCashBoxLine: time.Now(),
		}
	}
	close(queue)
	done <- struct{}{}
}

func MakeStaticCustumers(qSize int) []*model.Customer {
	q := make([]*model.Customer, qSize)
	for i := range q {
		tmp := &model.Customer{ID: i, MetricArrivedAtCashBoxLine: time.Now()}
		q[i] = tmp
	}
	return q
}

func MakeConcurrentCustomers(customers chan<- model.Customer, capacity int) {
	for i := 0; i < capacity; i++ {
		newDude := model.Customer{
			ID:                         i,
			PeopleBefore:               0,
			MetricCreatedTime:          time.Now(),
			MetricArrivedAtCashBoxLine: time.Now(),
			MetricArrivedToCashierTime: time.Now(),
			MetricLeftTime:             time.Now(),
		}
		customers <- newDude
	}
}
