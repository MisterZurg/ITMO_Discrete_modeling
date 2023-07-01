package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"fmt"
	"math/rand"
	"time"
)

func DynamicShortestQueue(qSize, registers int) {
	queue := make(chan model.Customer, qSize)
	cashiers := make([]chan model.Customer, registers)

	for idx := range cashiers {
		cashiers[idx] = make(chan model.Customer, qSize)
	}

	done := make(chan struct{})

	go func() {
		MakeCustumers(uint32(qSize), queue, done)
	}()

	for i := 0; i < registers; i++ {
		idx := i
		go func(idx int) {
			Cashier(cashiers[idx], idx, 100, done)
		}(idx)
	}

	go func() {
		FanOutUsers(queue, cashiers, 25)
	}()

	<-done
	<-done
	<-done
	<-done
}

func Cashier(customers chan model.Customer, id, maxTime int, done chan struct{}) {
	for {
		msg, ok := <-customers
		if ok {
			workingTime := rand.Intn(maxTime)
			time.Sleep(time.Duration(workingTime) * time.Millisecond)
			fmt.Printf("%d,%d,%d,%d\n",
				msg.ID,
				id,
				workingTime,
				msg.PeopleBefore,
			)
		} else {
			// fmt.Printf("Cashier ID=%d received all jobs", id)
			done <- struct{}{}
			return
		}
	}
}

func FanOutUsers(customers chan model.Customer, qs []chan model.Customer, maxTime int) {
	for user := range customers {
		userDest := rand.Intn(len(qs))
		minLen := len(qs[userDest])

		shoppingTime := rand.Intn(maxTime)
		time.Sleep(time.Duration(shoppingTime) * time.Millisecond)

		for i, q := range qs {
			if len(q) < minLen {
				minLen = len(q)
				userDest = i
			}
		}
		user.PeopleBefore = minLen
		qs[userDest] <- user
	}
	for _, q := range qs {
		close(q)
	}
}
