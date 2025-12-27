import matplotlib.pyplot as plt 

"""
Plot a 2D graph, exactly one of the default argument must be true 
"""
def plot(result: dict, 
         fixReq=False, fixAlgo=False, fixSelector=False, fixNodes=False): 

    if sum([fixReq, fixAlgo, fixSelector, fixNodes]) != 1: 
        raise ValueError("Exactly one default argument must be true")

     
    x = ""
    if fixReq: 
        x = "Number of Requests"
    elif fixAlgo: 
        pass

