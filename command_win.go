package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// /crew - Report a win for the crew members / a loss for the impostors
// /impostor - Report a win for the impostors / a loss for the crew
func commandWin(command string, args []string, m *discordgo.MessageCreate) {
	if len(args) != 1 && len(args) != 2 {
		discordSend(m.ChannelID, "You must provide the names of the impostor(s) in the game when reporting a win or loss.")
	}

	crewWon := command == "crew" || command == "town"

	var impostor1 string
	var impostor2 string
	match1 := mentionRegExp.FindStringSubmatch(args[0])
	if match1 == nil || len(match1) <= 1 {
		discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid Discord user.")
		return
	}
	impostor1 = match1[1]

	if len(args) == 2 {
		match2 := mentionRegExp.FindStringSubmatch(args[1])
		if match2 == nil || len(match2) <= 1 {
			discordSend(m.ChannelID, "\""+args[1]+"\" is not a valid Discord user.")
			return
		}
		impostor2 = match2[1]
	}

	activePlayers := playerListGetActivePlayers()

	winners := make([]string, 0)
	for _, p := range activePlayers {
		crew := p.DiscordID != impostor1 && p.DiscordID != impostor2
		win := (crew && crewWon) || (!crew && !crewWon)
		if err := p.UpdateStats(crew, win); err != nil {
			logger.Error("Failed to update the stats for \""+p.Username+"\":", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		}
		if win {
			winners = append(winners, p.Username)
		}
	}

	msg := "The "
	if crewWon {
		msg += "crew"
	} else {
		msg += "impostors"
	}
	msg += " won the game! Congratulations to: " + strings.Join(winners, " ")
	discordSend(m.ChannelID, msg)
}
