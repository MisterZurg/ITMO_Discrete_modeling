package internal

import (
	"ITMO_Discrete_modeling/Lab-3/internal/constants"
	"ITMO_Discrete_modeling/Lab-3/internal/model"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func ConcurrentDynamicShortestQueue(isDynamic bool, fCli, fServ *os.File) {
	users := make(chan model.Customer, constants.CUSTUMERS_NUMBER)
	var qs [constants.MAX_CASHIERS]*chan model.Customer
	for idx := 0; idx < constants.MIN_CASHIERS; idx++ {
		curChan := make(chan model.Customer, constants.CUSTUMERS_NUMBER)
		qs[idx] = &curChan
	}

	if isDynamic {
		ConcurrentDynamicQueue(users, qs, fCli, fServ)
	} else {
		ConcurrentStaticQueue(users, qs, fCli, fServ)
	}
}

func ConcurrentDynamicQueue(users chan model.Customer, qs [constants.MAX_CASHIERS]*chan model.Customer, fCli, fServ *os.File) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(users)
		MakeConcurrentCustomers(users, constants.CUSTUMERS_NUMBER)
	}()

	for i := 0; i < constants.MIN_CASHIERS; i++ {
		wg.Add(1)
		idx := i
		go func(idx int) {
			defer wg.Done()
			//fmt.Println("Cashier", idx, "started.")
			ConcurrentCashier(qs, idx, constants.MAX_WORKTIME, fCli, fServ)
		}(idx)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		FanOutUsersDynamic(users, qs, constants.CUSTUMER_ARRIVE_TIME, &wg, fCli, fServ)
	}()

	wg.Wait()
}

func ConcurrentCashier(qs [constants.MAX_CASHIERS]*chan model.Customer, id, maxTime, active_workers int, fCli, fServ *os.File) {
	servMsg := fmt.Sprintf("%d,%s,%d\n", id, "CREATED", time.Now().UnixMilli())
	fServ.WriteString(servMsg)
	// lastWorked := time.Now()
	for {
		// Cashier idle
		// if time.Now().Sub(lastWorked).Seconds() >= constants.SLEEP_TIME && id >= constants.MIN_CASHIERS {
		//if time.Now().Sub(lastWorked).Seconds() >= constants.SLEEP_TIME {
		//	qs[id] = nil
		//	servMsg = fmt.Sprintf("%d,%s,%d\n", id, "DELETED", time.Now().UnixMilli())
		//	fServ.WriteString(servMsg)
		//	return
		//}
		//rand.Intn(20)
		if rand.Intn(20) >= 10 && active_workers > constants.MIN_CASHIERS {
			qs[id] = nil
			servMsg = fmt.Sprintf("%d,%s,%d\n", id, "DELETED", time.Now().UnixMilli())
			fServ.WriteString(servMsg)
			return
		}

		cur, ok := <-*qs[id]
		if ok {
			cur.MetricArrivedToCashierTime = time.Now()
			workingTime := rand.Intn(maxTime)
			time.Sleep(time.Duration(workingTime) * time.Millisecond)
			cur.MetricLeftTime = time.Now()
			// Customer, Cashier, Created, ArrivedQueue, ArrivedCashier, Left
			msg := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d\n",
				cur.ID,
				id,
				cur.MetricCreatedTime.UnixMilli(),
				cur.MetricArrivedAtCashBoxLine.UnixMilli(),
				cur.MetricArrivedToCashierTime.UnixMilli(),
				cur.MetricLeftTime.UnixMilli(),
				cur.PeopleBefore,
			)
			fCli.WriteString(msg)
			// lastWorked = time.Now()
			servMsg = fmt.Sprintf("%d,%s,%d\n", id, "STARTED", cur.MetricArrivedToCashierTime.UnixMilli())
			fServ.WriteString(servMsg)
			servMsg = fmt.Sprintf("%d,%s,%d\n", id, "FINISHED", cur.MetricLeftTime.UnixMilli())
			fServ.WriteString(servMsg)
		} else {
			return
		}
	}
}

func FanOutUsersDynamic(customers chan model.Customer, qs [constants.MAX_CASHIERS]*chan model.Customer, maxTime int, wg *sync.WaitGroup, fCli, fServ *os.File) {
	for user := range customers {
		available := make([]int, 0)
		for i, q := range qs {
			if q == nil {
				continue
			}
			available = append(available, i)
		}
		userDest := rand.Intn(len(available))
		minLen := len(*qs[available[userDest]])

		shoppingTime := rand.Intn(maxTime)
		time.Sleep(time.Duration(shoppingTime) * time.Millisecond)

		// user looks at all the q's and considers, to wait or not
		emptyQ := false
		for _, q := range qs {
			if q == nil {
				continue
			}
			if len(*q) == 0 {
				emptyQ = true
				break
			}
		}

		isWaiting := true
		if !emptyQ {
			isWaiting = rand.Intn(2) == 0
			//fmt.Println("isWaiting", isWaiting)
		}

		if !isWaiting {
			// create new cashier, assign there
			newIdx := -1
			for idx, val := range qs {
				if val != nil {
					continue
				}
				// this one is nil
				curChan := make(chan model.Customer, constants.CUSTUMERS_NUMBER)
				qs[idx] = &curChan
				newIdx = idx
				break
			}

			if newIdx != -1 {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					//fmt.Println("Dynamic cashier", idx, "started.")
					ConcurrentCashier(qs, idx, constants.MAX_WORKTIME, fCli, fServ)
				}(newIdx)
				user.PeopleBefore = 0
				user.MetricArrivedAtCashBoxLine = time.Now()
				*qs[newIdx] <- user
				continue
			}
			//fmt.Println("User", user.Id, "couldn't create a new queue")
		}

		for i, q := range qs {
			if q == nil {
				continue
			}
			if len(*q) < minLen {
				minLen = len(*q)
				userDest = i
			}
		}
		user.PeopleBefore = minLen
		user.MetricArrivedAtCashBoxLine = time.Now()
		*qs[userDest] <- user
		//fmt.Println("User", user.Id, "went to queue", userDest)
	}
	for _, q := range qs {
		if q == nil {
			continue
		}
		close(*q)
	}
}

func ConcurrentStaticQueue(users chan model.Customer, qs [constants.MAX_CASHIERS]*chan model.Customer, fCli, fServ *os.File) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(users)
		MakeConcurrentCustomers(users, constants.CUSTUMERS_NUMBER)
	}()

	for i := 0; i < constants.MIN_CASHIERS; i++ {
		wg.Add(1)
		idx := i
		go func(idx int) {
			defer wg.Done()
			//fmt.Println("Cashier", idx, "started.")
			ConcurrentCashier(qs, idx, constants.MAX_WORKTIME, fCli, fServ)
		}(idx)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		FanoutUsersStatic(users, qs, constants.CUSTUMER_ARRIVE_TIME)
	}()

	wg.Wait()
}

func FanoutUsersStatic(customers chan model.Customer, qs [constants.MAX_CASHIERS]*chan model.Customer, maxTime int) {
	for user := range customers {
		userDest := rand.Intn(constants.MIN_CASHIERS)
		minLen := len(*qs[userDest])

		shoppingTime := rand.Intn(maxTime)
		time.Sleep(time.Duration(shoppingTime) * time.Millisecond)

		for i := 0; i < constants.MIN_CASHIERS; i++ {
			if len(*qs[i]) < minLen {
				minLen = len(*qs[i])
				userDest = i
			}
		}
		user.PeopleBefore = minLen
		user.MetricArrivedAtCashBoxLine = time.Now()
		*qs[userDest] <- user
	}
	for _, q := range qs {
		if q == nil {
			continue
		}
		close(*q)
	}
}
