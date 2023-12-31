package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dreamscached/minequery/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discordToken := os.Getenv("DISCORD_TOKEN")
	minecraftServerHost := os.Getenv("MINECRAFT_HOST")
	minecraftServerPort, err := strconv.Atoi(os.Getenv("MINECRAFT_PORT"))
	if err != nil {
		panic(err)
	}
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		panic(err)
	}
	pinger := minequery.NewPinger(
		minequery.WithTimeout(5*time.Second),
		minequery.WithUseStrict(true),
		minequery.WithProtocolVersion16(minequery.Ping16ProtocolVersion162),
		minequery.WithProtocolVersion17(minequery.Ping17ProtocolVersion172),
	)

	err = discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer discord.Close()
	countFails := 0
	go func() {
		for {
			time.Sleep(2 * time.Second)

			data, err := pinger.PingBeta18(minecraftServerHost, minecraftServerPort)
			if err != nil {
				log.Printf("Error querying server: %v", err)
				if countFails > 30 {
					err = discord.UpdateGameStatus(123, fmt.Sprintf("Offline"))
					if err != nil {
						log.Printf("Error updating game status Err: %s", err)
					}
					continue
				}
				countFails++
				continue
			} else {
				countFails = 0
			}
			err = discord.UpdateGameStatus(123, fmt.Sprintf("Игроков: %d/%d", data.OnlinePlayers, data.MaxPlayers))
			if err != nil {
				log.Printf("Error updating game status Err: %s", err)
			}

		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}
