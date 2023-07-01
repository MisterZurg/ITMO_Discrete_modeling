# -*- coding: cp1251 -*-
import numpy

from DEVS import *
import random
import numpy as np
import matplotlib.pyplot as plt
from scipy import stats
import scipy as sp
import pandas as pd

maxAngents = 100
arrivalRateMin = 1
arrivalRateMax = 8
# service_xk = np.arange(6) + 1
# service_pk = (0.1, 0.2, 0.3, 0.25, 0.1, 0.05)

# service_xk = np.arange(5)+1
# service_pk = (0.1, 0.2, 0.3, 0.25, 0.15)


service_xk = np.arange(3) + 1
service_pk = (0.2, 0.3, 0.5)

custm = stats.rv_discrete(name='custm', values=(service_xk, service_pk))


# ---- Customer Statistics ----
class customerStat:
    def __init__(self):
        self.id = -1
        self.arrivalTime = -1
        self.serviceTime = -1
        self.interArrivalTime = 0
        self.serviceBegins = -1
        self.waitingTimeInQueue = 0
        self.serviceEnds = -1
        self.timeInSystem = -1
        self.idleTimeOfServer = 0


# ---- Arrival Event ----
class ArrivalEvent:
    def __init__(self):
        self.eTime = 0.0

    def Execute(self):
        customer = customerStat()
        customer.id = DEVS.newId
        customer.arrivalTime = self.eTime
        if len(DEVS.stats) > 0:
            customer.interArrivalTime = customer.arrivalTime - DEVS.stats[-1].arrivalTime

        # print("Time %d" % self.eTime, " Arrival Event of agent {0}".format(customer.id))
        if DEVS.newId < maxAngents - 1:
            NextArrival = ArrivalEvent()
            NextArrival.eTime = self.eTime + random.randint(arrivalRateMin, arrivalRateMax)
            DEVS.EQ.AddEvent(NextArrival)

        # server is Free
        if DEVS.serverIdle == True:
            DEVS.serverIdle = False
            # print("server is Busy")
            Service = ServiceEvent()
            serviceTime = custm.rvs()
            customer.serviceTime = serviceTime
            customer.serviceBegins = self.eTime  # current time
            Service.eTime = self.eTime + serviceTime
            Service.id = customer.id
            DEVS.EQ.AddEvent(Service)
        # server is Busy
        else:
            # increase waiting line
            DEVS.customerQueue.append(customer.id)
            # print("customerQueue = %d" % len(DEVS.customerQueue))

        DEVS.newId = DEVS.newId + 1
        DEVS.stats.append(customer)


# ---- Service (END) Event ----
class ServiceEvent:
    def __init__(self):
        self.eTime = 0.0
        self.id = 0

    def Execute(self):
        ind = [i for i, val in enumerate(DEVS.stats) if val.id == self.id][0]
        DEVS.stats[ind].serviceEnds = self.eTime
        DEVS.stats[ind].timeInSystem = DEVS.stats[ind].serviceEnds - DEVS.stats[ind].arrivalTime
        DEVS.stats[ind].waitingTimeInQueue = DEVS.stats[ind].serviceBegins - DEVS.stats[
            ind].arrivalTime  # 0 without queue
        DEVS.stats[ind].idleTimeOfServer = DEVS.stats[ind].serviceBegins - DEVS.lastServedTime

        # print("Time %d" % self.eTime, "Service finished")
        if len(DEVS.customerQueue) > 0:
            qid = DEVS.customerQueue.pop(0)
            qind = [i for i, val in enumerate(DEVS.stats) if val.id == qid][0]
            Service = ServiceEvent()
            serviceTime = custm.rvs()
            Service.eTime = self.eTime + serviceTime
            Service.id = qid
            DEVS.stats[qind].serviceBegins = self.eTime
            DEVS.stats[qind].serviceTime = serviceTime
            DEVS.EQ.AddEvent(Service)
            # print("take new customer from the queue")
        else:
            DEVS.serverIdle = True
            # print("server is Idle (do nothing)")

        DEVS.lastServedTime = self.eTime


def runSimulation():
    # run simulation
    AE = ArrivalEvent()
    DEVS.EQ.AddEvent(AE)

    # simulation attributes
    DEVS.customerQueue = []
    DEVS.stats = []
    DEVS.newId = 0
    DEVS.serverIdle = True
    DEVS.lastServedTime = 0  # for Idle time

    # --- SIMULATION ---
    while DEVS.EQ.QueueSize() > 0:
        DEVS.ProcessNextEvent()

    # --- STATISTICS ---

    #  --- store all in file  ---
    f = open('output.csv', 'w')
    f.write(
        "Id;Interarrival Time;Arrival Time;Service Time;Time Service Begins;Waiting time in Queue;Time Service Ends;Time Customer Spends in System;Idle time of Server\n")
    for s in DEVS.stats:
        f.write("{0};{1};{2};{3};{4};{5};{6};{7};{8}\n".format(s.id, s.interArrivalTime, s.arrivalTime, s.serviceTime,
                                                               s.serviceBegins, s.waitingTimeInQueue, s.serviceEnds,
                                                               s.timeInSystem, s.idleTimeOfServer))
    f.close()

    numOfCustWhoWait = len([x for x in DEVS.stats if x.waitingTimeInQueue > 0])
    avTimeWhoWait = sum([x.waitingTimeInQueue for x in DEVS.stats]) / numOfCustWhoWait

    return avTimeWhoWait


def mean_confidence_interval(sample_df, eps=0.01, confidence=0.95):
    a = 1.0 * np.array(sample_df)
    n = len(a)
    m = np.mean(a)
    se = np.std(a, ddof=1) / np.sqrt(n)
    h = se * sp.stats.t._ppf((1 + confidence) / 2.0, n - 1)

    r = (sp.stats.t._ppf((1 + confidence) / 2.0, n - 1) * np.std(a, ddof=1) / eps) ** 2

    # print("sample_df")
    # print(sample_df)
    # print(m, h, r)
    # print("===")

    return m, h, r


def getConfidenceInterval(metric_df, metric, eps_percent=1, START_EXP_NUMS=5):
    eps = eps_percent * metric_df[metric].mean() / 100  # 1% epsilon

    # data = metric_df.sample(START_EXP_NUMS)[metric]
    start = len(metric_df)

    m, h, r = mean_confidence_interval(metric_df[metric], eps=eps)
    if h < eps:
        end = len(metric_df)
        print(f"eps: {eps :.3f}")
        print(f"ans: {m :.3f}+-{h :.3f}")
        print(f"R from: {start} to {end}")
        return True, 0

    if h >= eps:
        assert r < 1_000_000
        try:
            add_data = metric_df.sample(int(r + 1) - START_EXP_NUMS)[metric]
            metric_df = metric_df.append(add_data, ignore_index=True)
            START_EXP_NUMS += int(r + 1)
            if h < eps:
                end = len(metric_df)
                print(f"eps: {eps :.3f}")
                print(f"ans: {m :.3f}+-{h :.3f}")
                print(f"R from: {start} to {end}")
                return True, 0
        except:
            print(f"Want: {eps :.3f}")
            print("\nBigger R than population requested!\n")
            print(f"m={m}, h={h}, r={r}")
            print("\nContinue Experiments!\n")
            return False, int(int(r + 1) - len(metric_df))


# Core Logic
EXPERIMENT_NUMBER = 1000
OBSERVATION_PRECISION = 5  # Percent
OBSERVATION_NAME = "avg_wait_time"
START_EXP_NUMBER = 10

done = False

experiments_df = pd.DataFrame(columns=[OBSERVATION_NAME])
for experiment in range(EXPERIMENT_NUMBER):
    # If not done perform experements filling
    if not done:
        avTimeWhoWait_experiment = runSimulation()
        experiments_df.loc[experiment] = avTimeWhoWait_experiment

    # Filling DF
    if experiment < START_EXP_NUMBER:
        continue

    done, do_sample = getConfidenceInterval(experiments_df, OBSERVATION_NAME, OBSERVATION_PRECISION,
                                            int(START_EXP_NUMBER))

    if done:
        print(f"Confidence Interval found {experiment}/{EXPERIMENT_NUMBER}")
        break
    else:
        print(f"Going to increase sample")
        while experiment < START_EXP_NUMBER + do_sample:
            experiment += 1
            avTimeWhoWait_experiment = runSimulation()
            # list_row = [OBSERVATION_NAME, avTimeWhoWait_experiment]
            # df2 = df.append(new_row, ignore_index=True)
            new_row_dict = {OBSERVATION_NAME : avTimeWhoWait_experiment}
            #
            # experiments_df = experiments_df.append(new_row_dict, ignore_index=True)
            # experiments_df.loc[len(experiments_df)] = new_row_dict
            pd.concat([
                experiments_df,
                pd.DataFrame([pd.Series(new_row_dict)])]
            ).reset_index(drop=True)


        done, do_sample = getConfidenceInterval(experiments_df, OBSERVATION_NAME, OBSERVATION_PRECISION,
                                                 int(START_EXP_NUMBER) + do_sample)
        if done:
            print(f"Confidence Interval found {experiment}/{EXPERIMENT_NUMBER}")
            break

# print(experiments_df)
