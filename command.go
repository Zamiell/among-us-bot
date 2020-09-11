package main

import (
	"github.com/bwmarrin/discordgo"
)

// Used to store all of the functions that handle each command
var commandMap = make(map[string]func(string, []string, *discordgo.MessageCreate))

func commandInit() {
	commandMap["help"] = commandHelp
	commandMap["list"] = commandList
	commandMap["next"] = commandModifyPlayerList
	commandMap["add"] = commandModifyPlayerList     // Synonym for "/next"
	commandMap["waiting"] = commandModifyPlayerList // Synonym for "/next"
	commandMap["playing"] = commandModifyPlayerList
	commandMap["remove"] = commandModifyPlayerList
	commandMap["leave"] = commandModifyPlayerList  // Synonym for /remove"
	commandMap["delete"] = commandModifyPlayerList // Synonym for /remove"
	commandMap["unnext"] = commandModifyPlayerList // Synonym for /remove"
	commandMap["notplaying"] = commandModifyPlayerList
	commandMap["stopplaying"] = commandModifyPlayerList // Synonym for /notplaying"
	commandMap["ping"] = commandPing
	commandMap["pingnext"] = commandPing
	commandMap["clear"] = commandClear
	commandMap["clearall"] = commandClear
	commandMap["deleteall"] = commandClear
	commandMap["crew"] = commandWin
	commandMap["town"] = commandWin
	commandMap["impostor"] = commandWin
	commandMap["impostors"] = commandWin
	commandMap["imposter"] = commandWin
	commandMap["imposters"] = commandWin
	commandMap["mafia"] = commandWin
	commandMap["imp"] = commandWin
	commandMap["mafia"] = commandWin
	commandMap["stats"] = commandModifyPlayerList
	commandMap["minuscrew"] = commandModifyPlayerList
	commandMap["minustown"] = commandModifyPlayerList // Synonym for /minuscrew"
	commandMap["minusimp"] = commandModifyPlayerList
	commandMap["minusimpostor"] = commandModifyPlayerList // Synonym for /minusimp"
	commandMap["minusimposter"] = commandModifyPlayerList // Synonym for /minusimp"
	commandMap["minusmafia"] = commandModifyPlayerList    // Synonym for /minusimp"
	commandMap["vpn"] = commandVPN
}
