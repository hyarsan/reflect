/*

backgroundProcesses.go -
threads that run in the background

credits:
  - @hyarsan#3653 - original bot creator

license: gnu agplv3

*/

package main

import (
	// internals
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	// externals
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var userGuilds []*discordgo.UserGuild

// updates a slice of all of the guilds the client is in periodically
func backgroundGuildUpdater(s *discordgo.Session) {

	var (
		guilds    = []*discordgo.UserGuild{}
		guildList = []*discordgo.UserGuild{}
		err       error
	)

	for {

		var afterID string
		if len(guilds) == 0 {

			afterID = ""

		} else {

			afterID = guilds[len(guilds)-1].ID

		}

		guildList, err = s.UserGuilds(100, "", afterID)
		if err != nil {

			log.Printf("unable to update guild list. error: %v", err)

		}

		guilds = append(guilds, guildList...)

		if len(guildList) < 100 {

			userGuilds = guilds
			guilds = []*discordgo.UserGuild{}
			time.Sleep(10 * time.Second)

		}

	}

}

// updates the bot's status every 10 seconds
func backgroundStatusUpdater(s *discordgo.Session) {

	idleTime := 0

	for {

		err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			IdleSince: &idleTime,
			Game: &discordgo.Game{
				Name: fmt.Sprintf("#%s on %d servers!", config.ChannelName, len(userGuilds)),
				Type: 2,
			},
			AFK:    true,
			Status: fmt.Sprintf("#%s on %d servers!", config.ChannelName, len(userGuilds)),
		})
		if err != nil {

			log.Printf("unable to set the bot status. error: %v", err)

		}

		time.Sleep(10 * time.Second)

	}

}

// mirrors the messages
func backgroundMessageSend(s *discordgo.Session, g *discordgo.UserGuild, mc *discordgo.Channel, messageData *discordgo.MessageSend) {

	channels, err := guildChannelByName(s, g.ID, config.ChannelName)
	if err != nil {

		log.Printf("unable to grab reflect channel. error: %v", err)
		return

	}

	if g.ID == mc.GuildID {

		return

	}

	if len(channels) == 0 {

		return

	}

	_, err = s.ChannelMessageSendComplex(channels[0].ID, messageData)
	if err != nil {

		log.Errorf("unable to send message. error: %v\n", err)
		return

	}

}

// outputs the config to the file every 5 minutes
func backgroundConfigSave() {

	for {

		cfgByte, err := json.MarshalIndent(config, "", "	")
		if err != nil {

			log.Fatalf("[err]: unable to convert the config back to json. error: %v", err)

		}

		err = ioutil.WriteFile("config.json", cfgByte, 0644)
		if err != nil {

			log.Fatalf("[err]: unable to output config back to file. error: %v", err)

		}

		time.Sleep(5 * time.Minute)

	}

}
