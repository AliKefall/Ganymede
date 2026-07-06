CREATE TABLE ratings(
    user_id UUID NOT NULL,
    rating_type TEXT NOT NULL,
    rating INT NOT NULL DEFAULT 1500,
    games_played INT NOT NULL DEFAULT 0,
    PRIMARY KEY(user_id, rating_type),
    FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
