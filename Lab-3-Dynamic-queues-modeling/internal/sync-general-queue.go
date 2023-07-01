package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"fmt"
	"time"
)

func GeneralQueue(qSize int) {
	q := MakeStaticCustumers(qSize)
	GeneralStaticCashier(q)

}

func GeneralStaticCashier(queue []*model.Customer) {
	for _, customer := range queue {
		fmt.Printf(
			"ID %d, arrived %v processed at %v\n",
			customer.ID,
			customer.MetricArrivedAtCashBoxLine.Nanosecond(),
			time.Now().Nanosecond(),
		)
	}
}
