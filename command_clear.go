package main

import (
	"github.com/bwmarrin/discordgo"
)

// /clear - Remove every player from the list
func commandClear(command string, args []string, m *discordgo.MessageCreate) {
	playerList = make([]*Player, 0)
	if err := models.PlayerList.DeleteAll(); err != nil {
		logger.Error("Failed to delete all the players on the player list:", err)
		discordSend(m.ChannelID, ErrorMsg)
		return
	}

	discordSend(m.ChannelID, "Successfully cleared the player list & waiting list.")
}
