CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY UNIQUE, 
    name TEXT NOT NULL, 
    age INT NOT NULL, 
    email TEXT NOT NULL, 
    password TEXT NOT NULL,
); 
