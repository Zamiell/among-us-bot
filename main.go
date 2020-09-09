package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	ProjectName = "discord-waiting-list"
	ErrorMsg    = "Something went wrong. Please contact an administrator."
)

var (
	logger      *zap.SugaredLogger
	projectPath string
	models      *Models
)

func main() {
	// Initialize logging using the Zap library
	var zapLogger *zap.Logger
	if v, err := zap.NewDevelopment(); err != nil {
		log.Fatal("Failed to initialize logging:", err)
	} else {
		zapLogger = v
	}
	logger = zapLogger.Sugar()

	// Welcome message
	startText := "| Starting " + ProjectName + " |"
	borderText := "+" + strings.Repeat("-", len(startText)-2) + "+"
	logger.Info(borderText)
	logger.Info(startText)
	logger.Info(borderText)

	// Get the project path
	// https://stackoverflow.com/questions/18537257/
	if v, err := os.Executable(); err != nil {
		logger.Fatal("Failed to get the path of the currently running executable:", err)
	} else {
		projectPath = filepath.Dir(v)
	}

	// Check to see if the ".env" file exists
	envPath := path.Join(projectPath, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		logger.Fatal("The \"" + envPath + "\" file does not exist. Copy the \".env_template\" file to \".env\".")
		return
	} else if err != nil {
		logger.Fatal("Failed to check if the \""+envPath+"\" file exists:", err)
		return
	}

	// Load the ".env" file which contains environment variables with secret values
	if err := godotenv.Load(envPath); err != nil {
		logger.Fatal("Failed to load the \".env\" file:", err)
		return
	}

	// Initialize the database model (in "models.go")
	if v, err := modelsInit(); err != nil {
		logger.Fatal("Failed to open the database:", err)
		return
	} else {
		models = v
	}
	defer models.Close()

	// Get the players from the database (in "player_list.go")
	playerListInit()

	// Initialize the command map (in "command.go")
	commandInit()

	// Initialize the Discord connection (in "discord.go")
	discordInit()
	defer discord.Close()

	// Block until a terminal signal is received
	logger.Info(ProjectName + " is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := logger.Sync(); err != nil {
		log.Fatal("Failed to flush buffered log entries:", err)
	}
}
