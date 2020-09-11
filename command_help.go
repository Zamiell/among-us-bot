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
	msg += "/remove                               Remove yourself from the playing list / waiting list.\n"
	msg += "/playing                              Mark yourself as a player in the active game.\n"
	msg += "/notplaying                           Move yourself from the playing list to the waiting list.\n"
	msg += "/ping                                 Send a Discord ping to the next person on the waiting list.\n"
	msg += "/clear                                Remove all players from the playing list and waiting list.\n"
	msg += "```\n"
	msg += "Stats commands:\n"
	msg += "```\n"
	msg += "/crew [impostor1] [impostor2]   Report a win for the crew members / a loss for the impostors.\n"
	msg += "/imp [impostor1] [impostor2]    Report a loss for the crew members / a win for the impostors.\n"
	msg += "/stats                          See your statistics."
	msg += "/minuscrew                      Remove a town win"
	msg += "/minusimp                       Remove an impostor win"
	msg += "```\n"
	msg += "Info commands:\n"
	msg += "```\n"
	msg += "/vpn   Show the guide for how to connect to the private Among Us server."
	msg += "```"
	discordSend(m.ChannelID, msg)
}
