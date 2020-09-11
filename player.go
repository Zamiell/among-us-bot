package main

type Player struct {
	ID        int64
	Username  string
	DiscordID string
	Playing   bool
	Stats     Stats
}

type Stats struct {
	TotalGames       int
	NumCrewGames     int
	CrewWins         int
	NumImpostorGames int
	ImpostorWins     int
}

func (p *Player) Mention() string {
	return "<@" + p.DiscordID + ">"
}

func (p *Player) SetPlaying(playing bool) error {
	p.Playing = playing

	err := models.PlayerList.SetPlaying(p)
	return err
}

func (p *Player) UpdateStats(crew bool, win bool) error {
	p.Stats.TotalGames++
	if crew {
		p.Stats.NumCrewGames++
		if win {
			p.Stats.CrewWins++
		}
	} else {
		p.Stats.NumImpostorGames++
		if win {
			p.Stats.ImpostorWins++
		}
	}

	err := models.Players.UpdateStats(p)
	return err
}

func (p *Player) AdjustWin(positive bool, crew bool) error {
	if positive {
		p.Stats.TotalGames++
	} else {
		p.Stats.TotalGames--
	}

	if crew {
		if positive {
			p.Stats.NumCrewGames++
			p.Stats.CrewWins++
		} else {
			p.Stats.NumCrewGames--
			p.Stats.CrewWins--
		}
	} else {
		if positive {
			p.Stats.NumImpostorGames++
			p.Stats.ImpostorWins++
		} else {
			p.Stats.NumImpostorGames--
			p.Stats.ImpostorWins--
		}
	}

	err := models.Players.UpdateStats(p)
	return err
}
