package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandVPN(command string, args []string, m *discordgo.MessageCreate) {
	discordSend(m.ChannelID, "Instructions for how to connect: <https://github.com/Zamiell/among-us-vpn/blob/master/README.md>")
}
