package main

import (
	"ITMO_Discrete_modeling/Lab-3/internal"
	"ITMO_Discrete_modeling/Lab-3/internal/constants"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	DYNAMIC_Q = true
	STATIC_Q  = true
)

func main() {
	Lab3()
	// Lab4()
}

func Lab3() {
	fCli, err := os.Create(fmt.Sprintf("Lab-3/cmd/logs/clients_dynamic_1-000_%d.csv", constants.MIN_CASHIERS))
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fCli.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	fServ, err := os.Create(fmt.Sprintf("Lab-3/cmd/logs/servers_dynamic_1-000_%d.csv", constants.MIN_CASHIERS))
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fServ.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	fCli.WriteString("Customer,Cashier,Created,ArrivedQueue,ArrivedCashier,Left,PeopleBefore\n")
	fServ.WriteString("Server,Action,Time\n")

	// 1) General queue
	//internal.GeneralQueue(CUSTUMERS_NUMBER)

	// 2) Shortest queue
	//internal.ShortestStaticQueue(CUSTUMERS_NUMBER, MIN_CASHIERS)

	// 3) Dynamic General queue
	//internal.DynamicGeneralQueue(CUSTUMERS_NUMBER, MAX_WORKTIME)

	// 4) Dynamic Shortest queue
	//internal.DynamicShortestQueue(CUSTUMERS_NUMBER, MIN_CASHIERS)

	// 5) Generic Queue
	internal.ConcurrentDynamicShortestQueue(DYNAMIC_Q, fCli, fServ)
}

const (
	EXPERIMENT_NUMBER         = 500
	OBSERVATION_AVG_WAIT_TIME = "0.14"
	OBSERVATION_PRECISION     = "0.02"
	OBSERVATION_NAME          = "avg_wait_time"
)

func Lab4() {
	fMetrics, err := os.Create(fmt.Sprintf("Lab-4/metrics.csv"))
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fMetrics.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	for expNo := 1; expNo <= EXPERIMENT_NUMBER; expNo++ {
		//for subExp := 1; subExp <= 5; subExp++ {
		//	makeExperiment(expNo, subExp)
		//}

		out, err := exec.Command(
			"python3",
			"Lab-4/get-metrics-from-experiment.py",
			strconv.Itoa(expNo),
			OBSERVATION_AVG_WAIT_TIME,
			OBSERVATION_PRECISION,
			OBSERVATION_NAME,
		).Output()
		if err != nil {
			log.Fatal(err)
		}

		fMetrics.WriteString(string(out))

		// fmt.Printf("hw output: %s\n", out)
	}
}

func makeExperiment(experimentNumber, subExp int) {
	fCli, err := os.Create(fmt.Sprintf("Lab-4/experiments/ex_%d_%d_clients_dynamic_%d.csv", experimentNumber, subExp, constants.MIN_CASHIERS))
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fCli.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	fServ, err := os.Create(fmt.Sprintf("Lab-4/experiments/ex_%d_%d_servers_dynamic_%d.csv", experimentNumber, subExp, constants.MIN_CASHIERS))
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fServ.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	fCli.WriteString("Customer,Cashier,Created,ArrivedQueue,ArrivedCashier,Left,PeopleBefore\n")
	fServ.WriteString("Server,Action,Time\n")

	internal.ConcurrentDynamicShortestQueue(DYNAMIC_Q, fCli, fServ)
}
