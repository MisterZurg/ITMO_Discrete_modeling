package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"fmt"
	"math/rand"
	"time"
)

func DynamicGeneralQueue(qSize, MAX_WORKTIME int) {
	queue := make(chan model.Customer, qSize)
	done := make(chan struct{})

	go func() {
		GeneralCashier(queue, 0, MAX_WORKTIME, done)
	}()

	go func() {
		MakeCustumers(500, queue, done)
	}()

	<-done
	<-done
}

func GeneralCashier(queue chan model.Customer, id, maxTime int, done chan struct{}) {
	for {
		msg, more := <-queue
		if more {
			workingTime := rand.Intn(maxTime)
			time.Sleep(time.Duration(workingTime) * time.Millisecond)
			fmt.Printf("%d,%d,%d,%d\n",
				msg.ID,
				id,
				workingTime,
				msg.PeopleBefore,
			)
		} else {
			fmt.Println("received all jobs")
			done <- struct{}{}
			return
		}
	}
}
