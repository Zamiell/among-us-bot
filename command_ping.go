package main

import (
	"github.com/bwmarrin/discordgo"
)

// /ping - Send a Discord ping to the next person on the waiting list
func commandPing(command string, args []string, m *discordgo.MessageCreate) {
	player := playerListGetFirstWaiter()
	if player == nil {
		discordSend(m.ChannelID, "There is no-one currently on the waiting list.")
		return
	}

	discordSend(m.ChannelID, "Calling "+player.Mention()+" - there is now room for you in the current game!")
}
