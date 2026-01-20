# TEST ARCHITECTURE 

## How to run 
```bash
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
cd src
python -m load_manager_test.main
```

## Tree Structure 
```bash
.
├── README.md
├── requirements.txt
└── src
    └── load_manager_test
        ├── config
        │   ├── env.py
        │   ├── __init__.py
        │   ├── setup.py
        │   └── teardown.py
        ├── experiment
        │   ├── algorithm_exp.py
        │   ├── baseline_exp.py
        │   ├── exp.py
        │   ├── __init__.py
        │   ├── load_exp.py
        │   ├── scaling_exp.py
        │   ├── selector_exp.py
        │   └── strategy_exp.py
        ├── __init__.py
        ├── main.py
        └── plotting
            ├── __init__.py
            └── plot.py
```
