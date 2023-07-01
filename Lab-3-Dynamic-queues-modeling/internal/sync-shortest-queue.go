package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"fmt"
	"math/rand"
	"time"
)

func ShortestStaticQueue(qSize, registers int) {
	q := MakeStaticCustumers(qSize)
	//  Spawn sep cashiers
	cashiers := make([]model.Cashier, registers)
	for _, csh := range cashiers {
		// tmp := make([]model.Customer, qSize/registers)
		csh.Queue = make([]model.Customer, qSize/registers)
		fmt.Println(len(csh.Queue))
	}

	for _, user := range q {
		userDest := rand.Intn(len(cashiers))
		minLen := cashiers[userDest].Size

		//shoppingTime := rand.Intn(maxTime)
		//time.Sleep(time.Duration(shoppingTime) * time.Millisecond)
		for i, csh := range cashiers {
			if csh.Size < minLen {
				minLen = csh.Size
				userDest = i
			}
		}

		cashiers[userDest].Queue = append(cashiers[userDest].Queue, *user)
		cashiers[userDest].Size++

		for idx, csh := range cashiers {
			ShortestStaticCashier(idx, csh.Queue)
		}
	}
}

func ShortestStaticCashier(cashierID int, queue []model.Customer) {
	for _, customer := range queue {
		fmt.Printf(
			"Cashier ID %d, procceed %d, arrived %v processed at %v\n",
			cashierID,
			customer.ID,
			customer.MetricArrivedAtCashBoxLine.Nanosecond(),
			time.Now().Nanosecond(),
		)
	}
}
