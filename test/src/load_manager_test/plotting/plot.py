import matplotlib.pyplot as plt 
import numpy as np 

"""
General plot function

"""
def _plot(experiments: str, categories: list[str],
         throughputs: np.ndarray, title: str):

    # ensure category labels are strings and use numeric x positions
    cats = [str(c) for c in categories]
    x = np.arange(len(cats))
    width = 0.6

    # size figure based on number of categories to avoid overlap
    fig, ax = plt.subplots(figsize=(max(8, len(cats) * 0.6), 6))

    # draw bars with a fixed color and edge so they don't auto-cycle or appear stacked
    ax.bar(x, throughputs, width, color='tab:blue', edgecolor='k', zorder=3)

    ax.set_title(title)
    ax.set_xlabel(experiments)
    ax.set_ylabel('Throughput (ms)')

    ax.set_xticks(x)
    ax.set_xticklabels(cats, rotation=45, ha='right')

    plt.tight_layout()
    plt.savefig(title + '.png', bbox_inches='tight')


def plot(baseline: dict, results: dict, target: str):
    res_len = len(results['results'])
    throughputs = np.zeros(1 + res_len)
    throughputs[0] = np.sum(baseline['result'])
    categories = ['Base']

    for i in range(res_len):
        throughputs[i + 1] = np.sum(results['results'][i]['result'])
        categories.append(f'{target}: {results['results'][i][target]}')
    
    print(categories)
    _plot(results['experiment'], categories, throughputs, f'Base vs {target.capitalize()}')
