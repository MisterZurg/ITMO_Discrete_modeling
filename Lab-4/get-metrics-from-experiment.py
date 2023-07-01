import sys

import numpy as np
import pandas as pd
import scipy as sp
import scipy.stats
import warnings
warnings.filterwarnings("ignore")

def getMetricsDf(path, exp_no, min_cashiers, max_cashiers):
    first_exp = 1
    cols = list(getMetrics(path, first_exp, 1, min_cashiers, max_cashiers).keys())
    metrics_df = pd.DataFrame(columns=cols)

    for i in range(1, 5):
        metrics = getMetrics(path, exp_no, i, min_cashiers, max_cashiers)
        print(metrics)
        metrics_df = metrics_df._append(metrics, ignore_index=True)

    metrics_df.to_csv(f"./metrics_exp_{exp_no}.csv", index=False)

def getMetrics(path, exp_no, sub_no, min_cashiers, max_cashiers):
    # start_time = time.time()
    #  /data/notebook_files/logs/static/four_lines/ex_0_clients_dynamic_4.csv
    df_cli = pd.read_csv(f"./experiments/ex_{exp_no}_{sub_no}_clients_dynamic_{min_cashiers}.csv")
    df_serv = pd.read_csv(f"./experiments/ex_{exp_no}_{sub_no}_servers_dynamic_{min_cashiers}.csv")

    metrics = {}
    waiters = df_cli[df_cli["ArrivedCashier"] - df_cli["ArrivedQueue"] > 0]
    metrics["avg_wait_time"] = (df_cli["ArrivedCashier"] - df_cli["ArrivedQueue"]).mean()
    metrics["prob_wait"] = len(waiters) / len(df_cli)
    metrics["avg_service_time"] = (df_cli["Left"] - df_cli["ArrivedCashier"]).mean()
    metrics["avg_wait_time_waiters"] = (waiters["ArrivedCashier"] - waiters["ArrivedQueue"]).mean()
    metrics["avg_sys_time"] = (df_cli["Left"] - df_cli["ArrivedQueue"]).mean()

    arrivals = []
    idles = []
    works = []

    # Preporcessing cashiers logs
    for i in range(max_cashiers):
        cur_cli = df_cli[df_cli["Cashier"] == i]
        cur_serv = df_serv[df_serv["Server"] == i]

        if len(cur_cli) == 0:
            continue

        arrivals.append(cur_cli.diff()["ArrivedCashier"].mean())

        to_drop = [idx for idx in cur_serv.index[cur_serv["Action"] == "DELETED"]]
        to_drop.extend([idx+1 for idx in to_drop])
        cur_serv.drop(to_drop, inplace=True)
        # left with created-started-finished-started-finished-...-finished
        #            odd[0]-even[0]- odd[1] -even[1]- odd[2] -...- odd[n]
        odd_serv = cur_serv.iloc[::2]["Time"].to_numpy()
        even_serv = cur_serv.iloc[1::2]["Time"].to_numpy()
        assert len(odd_serv) == len(even_serv)+1, "Server logs even after deleting DELETE-CREATE pairs"
        idle = (even_serv-odd_serv[:-1]).sum()
        work = (odd_serv[1:]-even_serv).sum()
        idles.append(idle)
        works.append(work)

    metrics["avg_time_arrivals"] = np.array(arrivals).mean()
    metrics["prob_idle"] = np.array(idles).sum() / (np.array(idles).sum() + np.array(works).sum())
    # end_time = time.time()
    return metrics

"""
From the lecture # 8
"""
def mean_confidence_interval(data, eps=0.01, confidence=0.95):
    a = 1.0 * np.array(data)


    n = len(a)
    m = np.mean(a)
    se = np.std(a, ddof=1) / (np.sqrt(n))
    h = se * sp.stats.t._ppf((1+confidence)/2.0, n-1)
    r = (sp.stats.t._ppf((1+confidence)/2.0, n-1) * np.std(a, ddof=1) / eps)**2
    return m, h, r


import os

# os.chdir("./Lab-4")

"""
QUASI-CONSTANTS
"""
EXPERIMENT_NUMBER = sys.argv[1]
OBSERVATION_AVG_WAIT_TIME = sys.argv[2]
OBSERVATION_PRECISION = sys.argv[3]
OBSERVATION_NAME = sys.argv[4]

"""
DEBUG
"""
# EXPERIMENT_NUMBER         = 1
# OBSERVATION_AVG_WAIT_TIME = "0.20"
# OBSERVATION_PRECISION     = "0.01"
# OBSERVATION_NAME          = "avg_wait_time"

MIN_CSH, MAX_CSH = 4, 5
EXPERIMENT_PATH = "Lab-4/experiments"


getMetricsDf("./experiments", EXPERIMENT_NUMBER, MIN_CSH, MAX_CSH)
# getMetricsDf(EXPERIMENT_PATH, EXPERIMENT_NUMBER, MIN_CSH, MAX_CSH)


epsilon = float(OBSERVATION_PRECISION)
metrics_df = pd.read_csv(f"./metrics_exp_{EXPERIMENT_NUMBER}.csv")

m, h, r = mean_confidence_interval(metrics_df[OBSERVATION_NAME], eps=epsilon)
if h <= epsilon:
    print(f"eps: {epsilon :.3f}")
    print(f"ans: {m :.3f}+-{h :.3f}")
    print(f"R from: {EXPERIMENT_NUMBER} to {500}")
