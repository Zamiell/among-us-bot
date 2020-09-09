package main

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	discord               *discordgo.Session
	discordBotID          string
	discordGuildID        string
	discordListenChannels []string
)

func discordInit() {
	// Read some configuration values from environment variables
	discordToken := os.Getenv("DISCORD_TOKEN")
	if len(discordToken) == 0 {
		logger.Fatal("The \"DISCORD_TOKEN\" environment variable is blank.")
		return
	}
	discordGuildID = os.Getenv("DISCORD_GUILD_ID")
	if len(discordGuildID) == 0 {
		logger.Fatal("The \"DISCORD_GUILD_ID\" environment variable is blank.")
		return
	}
	discordListenChannelsString := os.Getenv("DISCORD_LISTEN_CHANNEL_IDS")
	if len(discordListenChannelsString) == 0 {
		logger.Fatal("The \"DISCORD_LISTEN_CHANNEL_IDS\" environment variable is blank.")
		return
	}
	discordListenChannels = strings.Split(discordListenChannelsString, ",")

	// Bot accounts must be prefixed with "Bot"
	if v, err := discordgo.New("Bot " + discordToken); err != nil {
		logger.Error("Failed to create a Discord session:", err)
		return
	} else {
		discord = v
	}

	// Register function handlers for various events
	discord.AddHandler(discordReady)
	discord.AddHandler(discordMessageCreate)

	// Open the websocket and begin listening
	if err := discord.Open(); err != nil {
		logger.Fatal("Failed to open the Discord session:", err)
		return
	}
}

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	logger.Info("Discord bot connected with username: " + event.User.Username)
	discordBotID = event.User.ID
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == discordBotID {
		return
	}

	// Ignore all messages that are in channels outside of the ones specified in the ".env" file
	if !stringInSlice(m.ChannelID, discordListenChannels) {
		return
	}

	// Get the channel
	var channel *discordgo.Channel
	if v, err := discord.Channel(m.ChannelID); err != nil {
		logger.Error("Failed to get the Discord channel of \""+m.ChannelID+"\":", err)
		return
	} else {
		channel = v
	}

	// Log the message
	logger.Info("[#" + channel.Name + "] " +
		"<" + m.Author.Username + "#" + m.Author.Discriminator + "> " + m.Content)

	args := strings.Split(m.Content, " ")
	command := args[0]
	args = args[1:] // This will be an empty slice if there is nothing after the command

	// Commands will start with a "/", so we can ignore everything else
	if !strings.HasPrefix(command, "/") {
		return
	}
	command = strings.TrimPrefix(command, "/")
	command = strings.ToLower(command) // Commands are case-insensitive

	// Check to see if there is a command handler for this command
	commandFunction, ok := commandMap[command]
	if ok {
		commandFunction(command, args, m)
	} else {
		discordSend(m.ChannelID, "That is not a valid command.")
	}
}

func discordSend(to string, msg string) {
	if _, err := discord.ChannelMessageSend(to, msg); err != nil {
		// Occasionally, sending messages to Discord can time out; if this occurs,
		// do not bother retrying, since losing a single message is fairly meaningless
		logger.Info("Failed to send \""+msg+"\" to Discord:", err)
		return
	}
}

func discordGetNickname(member *discordgo.Member) string {
	if member.Nick != "" {
		return member.Nick
	}

	return member.User.Username
}
