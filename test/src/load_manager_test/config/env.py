from dotenv import load_dotenv
import os 

class Env: 
    def __init__(self) -> None:
        load_dotenv() 

        self._POSTGRES_USER = os.getenv('POSTGRES_USER')
        self._POSTGRES_PASSWORD = os.getenv('POSTGRES_PASSWORD')
        self._POSTGRES_DB = os.getenv('POSTGRES_DB')
        self._POSTGRES_PORT = os.getenv('POSTGRES_PORT')
        self._POSTGRES_HOST = os.getenv('POSTGRES_HOST')

    def get_db_env(self) -> str: 
        return f'postgresq://{self._POSTGRES_USER}:{self._POSTGRES_PASSWORD}\
        @{self._POSTGRES_HOST}:{self._POSTGRES_HOST}/{self._POSTGRES_DB}'

