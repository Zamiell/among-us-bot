package main

import (
	"github.com/bwmarrin/discordgo"
)

// /help - Display a list of all commands
func commandHelp(command string, args []string, m *discordgo.MessageCreate) {
	msg := "The following commands are supported:\n"
	msg += "```\n"
	msg += "/list         Display a list of the current active players and the current waiting players.\n"
	msg += "/next         Add yourself to the waiting list.\n"
	msg += "/remove       Remove yourself from the playing list and/or waiting list.\n"
	msg += "/playing      Mark yourself as a player in the active game.\n"
	msg += "/notplaying   Move yourself from the playing list to the waiting list.\n"
	msg += "/ping         Send a Discord ping to the next person on the waiting list.\n"
	msg += "/clear        Remove all players from the playing list and waiting list.\n"
	msg += "```"
	discordSend(m.ChannelID, msg)
}
