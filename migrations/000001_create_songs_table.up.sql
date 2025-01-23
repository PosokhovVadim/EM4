CREATE TABLE songs(
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    link VARCHAR(1024),
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (group_name, name)
);