package main

import (
	"database/sql"
)

type Players struct{}

func (*Players) Get(discordID string) (bool, *Player, error) {
	var player Player
	if err := db.QueryRow(`
		SELECT
			id,
			username,
			discord_id,
			total_games,
			num_crew_games,
			crew_wins,
			num_impostor_games,
			impostor_wins
		FROM players
		WHERE discord_id = ?
	`, discordID).Scan(
		&player.ID,
		&player.Username,
		&player.DiscordID,
		&player.Stats.TotalGames,
		&player.Stats.NumCrewGames,
		&player.Stats.CrewWins,
		&player.Stats.NumImpostorGames,
		&player.Stats.ImpostorWins,
	); err == sql.ErrNoRows {
		return false, &player, nil
	} else if err != nil {
		return true, &player, err
	}

	return true, &player, nil
}

func (*Players) Insert(player *Player) (int64, error) {
	var res sql.Result
	if v, err := db.Exec(`
		INSERT INTO players (username, discord_id)
		VALUES (?, ?)
	`, player.Username, player.DiscordID); err != nil {
		return 0, err
	} else {
		res = v
	}

	var id int64
	if v, err := res.LastInsertId(); err != nil {
		return 0, err
	} else {
		id = v
	}

	return id, nil
}

func (*Players) UpdateStats(player *Player) error {
	_, err := db.Exec(
		`
			UPDATE players
			SET
				total_games = ?,
				num_crew_games = ?,
				crew_wins = ?,
				num_impostor_games = ?,
				impostor_wins = ?
			WHERE discord_id = ?
		`,
		player.Stats.TotalGames,
		player.Stats.NumCrewGames,
		player.Stats.CrewWins,
		player.Stats.NumImpostorGames,
		player.Stats.ImpostorWins,
		player.DiscordID,
	)
	return err
}
