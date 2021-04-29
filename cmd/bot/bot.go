package main

import (
	"github.com/bwmarrin/discordgo"
	rpc "github.com/plally/foxbot_gateway/rpc_server"
	"log"
	"os"
)

func main() {
	session, err := discordgo.New("Bot "+os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("error creating session", err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMembers
	session.AddHandler(onMemberAdd)

	err = session.Open()
	if err != nil {
		log.Fatal("error opening session", err)
	}

	rpc.StartRPCServer(session)
}
