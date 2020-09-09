package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// /win - Report a win for the crew members / a loss for the impostors
// /loss - Report a loss for the crew members / a win for the impostors
func commandWinLoss(command string, args []string, m *discordgo.MessageCreate) {
	if len(args) != 1 && len(args) != 2 {
		discordSend(m.ChannelID, "You must provide the names of the impostor(s) in the game when reporting a win or loss.")
	}

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
		win := (crew && command == "win") || (!crew && command == "loss")
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
	if command == "win" {
		msg += "crew"
	} else if command == "loss" {
		msg += "impostors"
	}
	msg += " won the game! Congratulations to: " + strings.Join(winners, " ")
	discordSend(m.ChannelID, msg)
}
