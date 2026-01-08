import matplotlib.pyplot as plt 
import numpy as np 

"""
General plot function

"""
def _plot(experiments: str, categories: list,
         throughputs: np.ndarray, title: str):

    plt.bar(categories, throughputs)
    plt.title(title)
    plt.xlabel(experiments)
    plt.ylabel('Througput (ms)') # This is what we measuring anyway 

    plt.savefig(title + '.png')


def plot(baseline: dict, results: dict, target: str):
    res_len = len(results['results'])
    throughputs = np.zeros(1 + res_len)
    throughputs[0] = np.sum(baseline['result'])
    categories = ['Base']

    for i in range(res_len):
        throughputs[i + 1] = np.sum(results['results'][i]['result'])
        categories.append(results['results'][i][target])
    _plot(results['experiment'], categories, throughputs, f'Base vs {target.capitalize()}')
