DROP TABLE IF EXISTS players;
CREATE TABLE players (
    id                  INTEGER  PRIMARY KEY,
    username            TEXT     NOT NULL     UNIQUE,
    discord_id          TEXT     NOT NULL     UNIQUE,
    total_games         INTEGER  NOT NULL     DEFAULT 0,
    num_crew_games      INTEGER  NOT NULL     DEFAULT 0,
    crew_wins           INTEGER  NOT NULL     DEFAULT 0,
    num_impostor_games  INTEGER  NOT NULL     DEFAULT 0,
    impostor_wins       INTEGER  NOT NULL     DEFAULT 0
);

DROP TABLE IF EXISTS player_list;
CREATE TABLE player_list (
    id         INTEGER  PRIMARY KEY,
    player_id  INTEGER  NOT NULL,
    playing    BOOLEAN  DEFAULT FALSE,
    FOREIGN KEY (player_id) REFERENCES players (id)
);
