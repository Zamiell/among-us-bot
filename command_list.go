package main

import (
	"github.com/bwmarrin/discordgo"
)

// /list - List the players on the list
func commandList(command string, args []string, m *discordgo.MessageCreate) {
	discordSend(m.ChannelID, playerListGetSummary())
}
