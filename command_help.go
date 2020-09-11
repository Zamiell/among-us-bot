package main

import (
	"github.com/bwmarrin/discordgo"
)

// /help - Display a list of all commands
func commandHelp(command string, args []string, m *discordgo.MessageCreate) {
	msg := "List commands:\n"
	msg += "```\n"
	msg += "/list                                 Display a list of the current active players & people waiting.\n"
	msg += "/next                                 Add yourself to the waiting list.\n"
	msg += "/next [username1] [username2]         Add one or more people to the waiting list.\n"
	msg += "/remove                               Remove yourself from the playing list / waiting list.\n"
	msg += "/remove [username1] [username2]       Remove one or more people from the playing list / waiting list.\n"
	msg += "/playing                              Mark yourself as a player in the active game.\n"
	msg += "/playing [username1] [username2]      Mark one or more people as a player in the active game.\n"
	msg += "/notplaying                           Move yourself from the playing list to the waiting list.\n"
	msg += "/notplaying [username1] [username2]   Move one or more people from the playing list to the waiting list.\n"
	msg += "/ping                                 Send a Discord ping to the next person on the waiting list.\n"
	msg += "/clear                                Remove all players from the playing list and waiting list.\n"
	msg += "```\n"
	msg += "Stats commands:\n"
	msg += "```\n"
	msg += "/crew [impostor1] [impostor2]        Report a win for the crew members / a loss for the impostors.\n"
	msg += "/imposters [impostor1] [impostor2]   Report a loss for the crew members / a win for the impostors.\n"
	msg += "/stats                               See your statistics."
	msg += "```\n"
	msg += "Info commands:\n"
	msg += "```\n"
	msg += "/vpn   Show the guide for how to connect to the private Among Us server."
	msg += "```"
	discordSend(m.ChannelID, msg)
}
