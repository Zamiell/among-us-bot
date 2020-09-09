package main

import (
	"database/sql"
)

type PlayerList struct{}

func (*PlayerList) GetAll() ([]*Player, error) {
	players := make([]*Player, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			players.id,
			players.username,
			players.discord_id,
			players.total_games,
			players.num_crew_games,
			players.crew_wins,
			players.num_imposter_games,
			players.imposter_wins,
			player_list.playing
		FROM player_list
		JOIN players ON players.id = player_list.player_id
		ORDER BY player_list.id
	`); err != nil {
		return playerList, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		var player Player
		if err := rows.Scan(
			&player.ID,
			&player.Username,
			&player.DiscordID,
			&player.Stats.TotalGames,
			&player.Stats.NumCrewGames,
			&player.Stats.CrewWins,
			&player.Stats.NumImposterGames,
			&player.Stats.ImposterWins,
			&player.Playing,
		); err != nil {
			return nil, err
		}
		players = append(players, &player)
	}

	if err := rows.Err(); err != nil {
		return players, err
	}

	return players, nil
}

func (*PlayerList) Insert(player *Player) error {
	_, err := db.Exec(`
		INSERT INTO player_list (player_id, playing)
		VALUES (?, ?)
	`, player.ID, player.Playing)
	return err
}

func (*PlayerList) SetPlaying(player *Player) error {
	_, err := db.Exec(`
		UPDATE player_list
		SET playing = ?
		WHERE player_id = ?
	`, player.Playing, player.ID)
	return err
}

func (*PlayerList) Delete(player *Player) error {
	_, err := db.Exec(`
		DELETE FROM player_list
		WHERE player_id = ?
	`, player.ID)
	return err
}

func (*PlayerList) DeleteAll() error {
	_, err := db.Exec("DELETE FROM player_list")
	return err
}
