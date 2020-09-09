package main

import (
	"errors"
	"strconv"
)

// For storing the players who are waiting for the next game to start
var playerList = make([]*Player, 0)

func playerListInit() {
	if v, err := models.PlayerList.GetAll(); err != nil {
		logger.Fatal("Failed to get the player list from the database:", err)
	} else {
		playerList = v
	}

	if len(playerList) != 0 {
		logger.Info("Restored " + strconv.Itoa(len(playerList)) + " player list entries from the database.")
	}
}

func playerListAdd(player *Player) error {
	playerList = append(playerList, player)
	err := models.PlayerList.Insert(player)
	return err
}

func playerListDelete(player *Player) error {
	index := playerListGetIndex(player.DiscordID)
	if index == -1 {
		return errors.New("player is not in the list")
	}

	playerList = append(playerList[:index], playerList[index+1:]...)
	err := models.PlayerList.Delete(player)
	return err
}

func playerListGetIndex(discordID string) int {
	for i, p := range playerList {
		if p.DiscordID == discordID {
			return i
		}
	}

	return -1
}

func playerListGetActivePlayers() []*Player {
	activePlayers := make([]*Player, 0)
	for _, p := range playerList {
		if p.Playing {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}

func playerListGetWaitingPlayers() []*Player {
	waitingPlayers := make([]*Player, 0)
	for _, p := range playerList {
		if !p.Playing {
			waitingPlayers = append(waitingPlayers, p)
		}
	}
	return waitingPlayers
}

func playerListGetSummary() string {
	msg := "**Active players:**\n"
	activePlayers := playerListGetActivePlayers()
	if len(activePlayers) == 0 {
		msg += "[no-one is currently playing]\n"
	} else {
		for i, p := range activePlayers {
			msg += strconv.Itoa(i+1) + ") " + p.Username + "\n"
		}
	}

	msg += "\n**Waiting players:**\n"
	waitingPlayers := playerListGetWaitingPlayers()
	if len(waitingPlayers) == 0 {
		msg += "[no-one is currently waiting]"
	} else {
		for i, p := range waitingPlayers {
			msg += strconv.Itoa(i+1) + ") " + p.Username + "\n"
		}
	}

	return msg
}

func playerListGetFirstWaiter() *Player {
	for _, p := range playerList {
		if !p.Playing {
			return p
		}
	}

	return nil
}
