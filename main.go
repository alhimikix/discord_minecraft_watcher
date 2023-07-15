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
	discordToken := os.Getenv("DISCORD_TOKEN")                            //"MTEyOTY2ODc4MTk2MzIzNTMyOA.GWPYJM.Hr402kKsciqQo27_NpZtaEXPTKAKVOyuFXcrC0" //"RYtbDdzl_-XViFyBahY1Ni6iVA9CqQvS" //flag.Bool("color", false, "display colorized output")
	minecraftServerHost := os.Getenv("MINECRAFT_HOST")                    //"141.101.24.182"
	minecraftServerPort, err := strconv.Atoi(os.Getenv("MINECRAFT_PORT")) //25765
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

	go func() {
		for {
			data, err := pinger.PingBeta18(minecraftServerHost, minecraftServerPort)
			if err != nil {
				log.Printf("Error querying server: %v", err)
				continue
			}
			err = discord.UpdateGameStatus(123, fmt.Sprintf("Игроков: %d/%d", data.OnlinePlayers, data.MaxPlayers))
			if err != nil {
				log.Printf("Error updating game status Err: %s", err)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}
