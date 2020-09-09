package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var mentionRegExp = regexp.MustCompile(`<@!*(\d+?)>`)

// /next - Add yourself to the waiting list
// /playing - Mark yourself as a player in the active game
// /remove - Remove yourself from the playing list / waiting list
// /notplaying - Move yourself from the playing list to the waiting list
// /stats - See your statistics
func commandModifyPlayerList(command string, args []string, m *discordgo.MessageCreate) {
	// Get the target player
	// If no arguments were provided, then the target is the person typing out the command
	// If an argument is provided, then the target is the person supplied as the argument
	var discordID string
	if len(args) == 0 {
		discordID = m.Author.ID
	} else {
		match := mentionRegExp.FindStringSubmatch(args[0])
		if match == nil || len(match) <= 1 {
			discordSend(m.ChannelID, "\""+args[0]+"\" is not a valid Discord user.")
			return
		}
		discordID = match[1]
	}

	var member *discordgo.Member
	if v, err := discord.GuildMember(discordGuildID, discordID); err != nil {
		discordSend(m.ChannelID, "\""+discordID+"\" is not a valid Discord user.")
		return
	} else {
		member = v
	}

	// Get this player from the database
	var exists bool
	var player *Player
	if v1, v2, err := models.Players.Get(discordID); err != nil {
		logger.Error("Failed to check to see if the Discord ID "+discordID+" exists:", err)
		discordSend(m.ChannelID, ErrorMsg)
		return
	} else {
		exists = v1
		player = v2
	}

	if !exists {
		nickname := discordGetNickname(member)
		logger.Info("Creating a new row for user:", nickname)
		player = &Player{
			Username:  nickname,
			DiscordID: discordID,
		}
		if id, err := models.Players.Insert(player); err != nil {
			logger.Error("Failed to insert \""+player.Username+"\" into the database:", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		} else {
			player.ID = id
		}
	}

	// If this player is already on the player list, use the player object already inside of the
	// list so that we can perform direct operations on it
	index := playerListGetIndex(player)
	if index != -1 {
		player = playerList[index]
	}

	var msg string
	if command == "next" || command == "add" || command == "waiting" {
		if index != -1 {
			discordSend(m.ChannelID, player.Username+" is already on the playing list / waiting list, so there is no need to perform this command.")
			return
		}

		if err := playerListAdd(player); err != nil {
			logger.Error("Failed to add \""+player.Username+"\" to the player list:", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		}
		msg = player.Username + " has been added to the waiting list.\n"
	} else if command == "playing" {
		if index != -1 && player.Playing {
			discordSend(m.ChannelID, player.Username+" is already marked as an active player, so there is no need to perform this command.")
			return
		}

		if index == -1 {
			if err := playerListAdd(player); err != nil {
				logger.Error("Failed to add \""+player.Username+"\" to the player list:", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
		}

		if err := player.SetPlaying(true); err != nil {
			logger.Error("Failed to update the status of \""+player.Username+"\" in the player list:", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		}
	} else if command == "remove" || command == "delete" || command == "unnext" {
		if index == -1 {
			discordSend(m.ChannelID, player.Username+" is not on the playing list / waiting list, so there is no need to perform this command.")
			return
		}

		if err := playerListDelete(player); err != nil {
			logger.Error("Failed to remove \""+player.Username+"\" from the player list:", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		}
		msg = player.Username + " has been removed from the player list / waiting list."
	} else if command == "notplaying" {
		if index != -1 && !player.Playing {
			discordSend(m.ChannelID, player.Username+" is already on the waiting list, so there is no need to perform this command.")
			return
		}

		if index == -1 {
			discordSend(m.ChannelID, player.Username+" is not on the playing list, so there is no need to perform this command.")
			return
		}

		if err := player.SetPlaying(false); err != nil {
			logger.Error("Failed to update the status of \""+player.Username+"\" in the player list:", err)
			discordSend(m.ChannelID, ErrorMsg)
			return
		}
	} else if command == "stats" {
		msg = "Stats for " + player.Username + ":\n"
		msg += "- Total games: **" + strconv.Itoa(player.Stats.TotalGames) + "**\n"

		crewWinRate := float64(player.Stats.CrewWins) / float64(player.Stats.NumCrewGames) * 100
		crewWinRateString := fmt.Sprintf("%.2f", crewWinRate)
		msg += "- Crew wins: " + strconv.Itoa(player.Stats.CrewWins) + " / " + strconv.Itoa(player.Stats.NumCrewGames) + " "
		msg += "(" + crewWinRateString + "%)\n"

		imposterWinRate := float64(player.Stats.ImposterWins) / float64(player.Stats.NumImposterGames) * 100
		imposterWinRateString := fmt.Sprintf("%.2f", imposterWinRate)
		msg += "- Imposter wins: " + strconv.Itoa(player.Stats.ImposterWins) + " / " + strconv.Itoa(player.Stats.NumImposterGames) + " "
		msg += "(" + imposterWinRateString + "%)\n"

		msg += "- Total games: " + strconv.Itoa(player.Stats.TotalGames) + "\n"
		discordSend(m.ChannelID, msg)
		return
	}

	msg += playerListGetSummary()
	discordSend(m.ChannelID, msg)
}
