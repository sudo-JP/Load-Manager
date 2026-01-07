import matplotlib.pyplot as plt 
import numpy as np 

"""
General plot function
"""
def plot(x_label: str, y_label: str, 
         category: np.ndarray,
         xs: np.ndarray, ys: np.ndarray, 
         exp_name: str):

    fig, ax = plt.subplots()
    width = 0.35

    bar1 = ax.bar(len(xs) - width / 2, xs,  width, yerr=xs,
                  label=x_label)
    bar2 = ax.bar(len(ys) + width / 2, ys,  width, yerr=ys)
    #ax.bar(x )

    plt.plot(xs, ys)
    plt.xlabel(x_label)
    plt.ylabel(y_label)

    plt.title(exp_name)

    ax.legend()
    plt.savefig(exp_name)


