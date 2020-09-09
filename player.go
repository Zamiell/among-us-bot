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
	NumImposterGames int
	ImposterWins     int
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
		p.Stats.NumImposterGames++
		if win {
			p.Stats.ImposterWins++
		}
	}

	err := models.Players.UpdateStats(p)
	return err
}