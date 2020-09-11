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
	// Get the target player(s)
	// If no arguments were provided, then the target is the person typing out the command
	// If an argument is provided, then the target is the person supplied as the argument
	discordIDs := make([]string, 0)
	if len(args) == 0 {
		discordID := m.Author.ID
		discordIDs = append(discordIDs, discordID)
	} else {
		for _, arg := range args {
			match := mentionRegExp.FindStringSubmatch(arg)
			if match == nil || len(match) <= 1 {
				discordSend(m.ChannelID, "\""+arg+"\" is not a valid Discord user.")
				return
			}
			discordID := match[1]
			discordIDs = append(discordIDs, discordID)
		}
	}

	// Get this player from the database
	players := make([]*Player, 0)
	for _, discordID := range discordIDs {
		var player *Player
		index := playerListGetIndex(discordID)
		if index != -1 {
			// They are already on the list, so get the player object from the list
			player = playerList[index]
		} else {
			// They are not already on the list, so get the player object from the database
			// (and create it in the database if it does not exist)
			var exists bool
			if v1, v2, err := models.Players.Get(discordID); err != nil {
				logger.Error("Failed to check to see if the Discord ID "+discordID+" exists:", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			} else {
				exists = v1
				player = v2
			}

			if !exists {
				var member *discordgo.Member
				if v, err := discord.GuildMember(discordGuildID, discordID); err != nil {
					discordSend(m.ChannelID, "\""+discordID+"\" is not a valid Discord user.")
					return
				} else {
					member = v
				}

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
		}

		players = append(players, player)
	}

	for _, player := range players {
		index := playerListGetIndex(player.DiscordID)

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
		} else if command == "remove" ||
			command == "leave" ||
			command == "delete" ||
			command == "unnext" {

			if index == -1 {
				discordSend(m.ChannelID, player.Username+" is not on the playing list / waiting list, so there is no need to perform this command.")
				return
			}

			if err := playerListDelete(player); err != nil {
				logger.Error("Failed to remove \""+player.Username+"\" from the player list:", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
		} else if command == "notplaying" || command == "stopplaying" {
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
			msg := "Stats for **" + player.Username + "**:\n"
			msg += "```\n"
			msg += "- Total games:   " + strconv.Itoa(player.Stats.TotalGames) + "\n"

			crewWinRate := float64(player.Stats.CrewWins) / float64(player.Stats.NumCrewGames) * 100
			crewWinRateString := fmt.Sprintf("%.2f", crewWinRate)
			msg += "- Crew wins:     " + strconv.Itoa(player.Stats.CrewWins) + " / " + strconv.Itoa(player.Stats.NumCrewGames) + " "
			if player.Stats.NumCrewGames > 0 {
				msg += "(" + crewWinRateString + "%)"
			}
			msg += "\n"

			impostorWinRate := float64(player.Stats.ImpostorWins) / float64(player.Stats.NumImpostorGames) * 100
			impostorWinRateString := fmt.Sprintf("%.2f", impostorWinRate)
			msg += "- Impostor wins: " + strconv.Itoa(player.Stats.ImpostorWins) + " / " + strconv.Itoa(player.Stats.NumImpostorGames) + " "
			if player.Stats.NumImpostorGames > 0 {
				msg += "(" + impostorWinRateString + "%)"
			}
			msg += "\n"

			msg += "```"
			discordSend(m.ChannelID, msg)
			return
		} else if command == "pluscrew" || command == "plustown" {
			if err := player.AdjustWin(true, true); err != nil {
				logger.Error("Failed to plus a crew win for \""+player.Username+"\":", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
			discordSend(m.ChannelID, "Removed a crew win from: "+player.Username)
		} else if command == "plusimp" || command == "plusimpostor" || command == "plusimposter" || command == "plusmafia" {
			if err := player.AdjustWin(true, false); err != nil {
				logger.Error("Failed to plus a crew win for \""+player.Username+"\":", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
			discordSend(m.ChannelID, "Removed an impostor win from: "+player.Username)
		} else if command == "minuscrew" || command == "minustown" {
			if err := player.AdjustWin(false, true); err != nil {
				logger.Error("Failed to minus a crew win for \""+player.Username+"\":", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
			discordSend(m.ChannelID, "Removed a crew win from: "+player.Username)
		} else if command == "minusimp" || command == "minusimpostor" || command == "minusimposter" || command == "minusmafia" {
			if err := player.AdjustWin(false, false); err != nil {
				logger.Error("Failed to minus a crew win for \""+player.Username+"\":", err)
				discordSend(m.ChannelID, ErrorMsg)
				return
			}
			discordSend(m.ChannelID, "Removed an impostor win from: "+player.Username)
		}
	}

	msg := playerListGetSummary()
	discordSend(m.ChannelID, msg)
}
