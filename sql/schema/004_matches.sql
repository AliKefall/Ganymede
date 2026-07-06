

CREATE TABLE matches(
    id UUID PRIMARY KEY,

    white_id UUID NOT NULL,
    black_id UUID NOT NULL,

    time_control TEXT NOT NULL,

    white_rating_before INT NOT NULL,
    black_rating_before INT NOT NULL,

    white_rating_after INT NOT NULL,
    black_rating_after INT NOT NULL,

    result TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,

    FOREIGN KEY(white_id)
        REFERENCES users(id),

    FOREIGN KEY (black_id)
        REFERENCES users(id)

);

CREATE TABLE match_moves(
    id BIGSERIAL PRIMARY KEY,
    match_id UUID NOT NULL,

    move_number INT NOT NULL,
    player_id UUID NOT NULL,

    san TEXT NOT NULL,
    uci TEXT NOT NULL,

    fen_after TEXT NOT NULL,

    white_time_ms BIGINT NOT NULL,
    black_time_ms BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,

    FOREIGN KEY(match_id)
        REFERENCES matches(id)
        ON DELETE CASCADE,

    FOREIGN KEY(player_id)
        REFERENCES users(id),

    UNIQUE(match_id, move_number)
);
