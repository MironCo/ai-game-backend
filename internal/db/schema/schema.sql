CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    unity_id TEXT UNIQUE NOT NULL,
    phone_number TEXT,
);

CREATE TABLE messages {
    id SERIAL PRIMARY KEY,
    unity_id TEXT NOT NULL,
    message TEXT,
    sender VARCHAR(16)
    sent_to VARCHAR(16)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}

CREATE TABLE texts {
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    unity_id VARCHAR (255) NOT NULL,
    message TEXT NOT NULL,
    sender_number VARCHAR(50) NOT NULL,
    receiver_number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    player_numbeR VARCHAR (15)
}

CREATE TABLE events {
    id SERIAL PRIMARY KEY,
    unity_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_details TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
}