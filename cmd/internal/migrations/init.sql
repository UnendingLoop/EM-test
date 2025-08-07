-- Active: 1752936028168@@127.0.0.1@5432@emtest
CREATE TABLE IF NOT EXISTS subscriptions(
    user_id TEXT PRIMARY KEY UNIQUE,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);