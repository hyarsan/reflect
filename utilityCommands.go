/*

utilityCommands.go -
the utility command code for reflect

credits:
  - @hyarsan#3653 - original bot creator

license: gnu agplv3

*/

package main

import (
	// internals
	"fmt"
	"strings"
	"math"
	"time"
	// externals
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// function that registers the utility commands
func registerUtilityCommands() {

	handler.AddCommand("help", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		embed := discordgo.MessageEmbed{
			Title:       "Reflect",
			Description: "A bot that links servers together like never before.",
			Color:       0x89da72,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://u.catgirl.host/4seqjs.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "How do I use it?",
					Value: fmt.Sprintf("Type %ssetup to set it up.\nIt creates a channel named #%s where you can talk to other guilds.", config.Prefix, config.ChannelName),
				},
				&discordgo.MessageEmbedField{
					Name: "Commands",
					Value: `**help**: Shows this message.
**setup**: Sets up the bot environment in your server.
**unsetup**: Removes the bot environment from the server.
**ban** **<nick|id|mention>**: Bans a user from Reflect.
**unban** **<nick|id|mention>**: Unbans a user from Reflect.
**bans**: Shows the banlist.
**info**: Shows some stats about the bot.
**user** **<nick|id|mention>**: Shows information about a user.
**notify** **<message>**: Sends out a message using the notification system.`,
				},
				&discordgo.MessageEmbedField{
					Name:  "Invite",
					Value: fmt.Sprintf("[Link](https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8)", reflectUser.ID),
				},
			},
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {

			log.Errorf("unable to send help message. error: %v", err)
			return

		}

	})
	handler.AddCommand("info", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		guildNames := allGuildNames()

		embed := discordgo.MessageEmbed{
			Title:       "Stats",
			Description: "Check up on Reflect",
			Color:       0x89da72,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://u.catgirl.host/4seqjs.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Server Count",
					Value: fmt.Sprintf("Reflect is currently in %d servers...", len(userGuilds)),
				},
				&discordgo.MessageEmbedField{
					Name:  "Server List Excerpt",
					Value: strings.Join(shuffleStringSlice(guildNames[:10]), "\n"),
				},
				&discordgo.MessageEmbedField{
					Name: "System Info",
					Value: fmt.Sprintf(`**Uptime**: %dh
**Operating System**: %s %s %s
**User**: %s@%s`, int(math.Floor(time.Since(startTime).Hours())), strings.Title(hostInfo.PlatformFamily), strings.Title(hostInfo.OS), strings.Title(hostInfo.PlatformVersion), currentUser.Username, hostInfo.Hostname),
				},
			},
		}

		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {

			log.Errorf("unable to send status message. error %v", err)
			return

		}

	})
	
	handler.AddCommand("notify", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	
		if !inList(m.Author.ID, config.Owners) {

			_, err = s.ChannelMessageSend(m.ChannelID, "You aren't an owner!")
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}
		
		if len(args) < 1 {
		
			_, err := s.ChannelMessageSend(m.ChannelID, "You didn't provide a message!")
			if err != nil {
			
				log.Errorf("unable to send message. error: %v", err)
				
			}
			return
			
		}
		
		for _, userGuild := range userGuilds {
		
			guild, err := s.Guild(userGuild.ID)
			if err != nil {
			
				log.Errorf("unable to retrieve guild. error: %v", err)
				
			}
			
			for _, channel := range guild.Channels {
			
				if channel.Name == "megachat" {
				
					_, err := s.ChannelMessageSend(channel.ID, fmt.Sprintf("<@%s> %s", guild.OwnerID, strings.Join(args, " ")))
					if err != nil {
					
						log.Errorf("unable to send message. error: %v", err)
						continue
						
					}
					break
					
				}
				
			}
			
		}
		
		_, err := s.ChannelMessageSend(m.ChannelID, "Successfully sent your notification!")
		if err != nil {
		
			log.Errorf("unable to send message. error: %v", err)
			
		}
		
	})

	handler.AddCommand("setup", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		ogMessage, err := s.ChannelMessageSend(m.ChannelID, "Setting up Reflect in your server...")
		if err != nil {

			log.Errorf("unable to send a message. error: %v", err)
			return

		}

		channel, err := s.Channel(m.ChannelID)
		if err != nil {

			log.Errorf("unable to retrieve the channel. error: %v", err)
			return

		}

		channelList, err := guildChannelByName(s, channel.GuildID, config.ChannelName)
		if err != nil {

			log.Errorf("unable to check for the channel in a guild. error: %v", err)
			return

		}

		var reflectChannel *discordgo.Channel

		if len(channelList) > 0 {

			ogMessage, err = s.ChannelMessageEdit(ogMessage.ChannelID, ogMessage.ID, fmt.Sprintf("Detected existing mirror channel at <#%s>. Using that channel instead...", channelList[0].ID))
			reflectChannel = channelList[0]

		} else {

			reflectChannel, err = s.GuildChannelCreate(channel.GuildID, config.ChannelName, discordgo.ChannelTypeGuildText)
			if err != nil {

				ogMessage, err = s.ChannelMessageEdit(ogMessage.ChannelID, ogMessage.ID, "Unable to set up Reflect in your server... Please check if I have a role that can create channels...")
				if err != nil {

					log.Errorf("unable to edit a message. error: %v", err)

				}

				return

			}

		}

		reflectChannel, err = s.ChannelEditComplex(reflectChannel.ID, &discordgo.ChannelEdit{
			Topic: fmt.Sprintf("Channel that links together discord servers using <@%s>. Please avoid sending NSFW material in here and please do not spam.", reflectUser.ID),
		})
		if err != nil {

			ogMessage, err = s.ChannelMessageEdit(ogMessage.ChannelID, ogMessage.ID, "Unable to set up Reflect in your server... Please check if I have a role that can edit channels...")
			if err != nil {

				log.Errorf("unable to edit a message. error: %v", err)

			}

			return

		}

		ogMessage, err = s.ChannelMessageEdit(ogMessage.ChannelID, ogMessage.ID, fmt.Sprintf("Setup is now finished. Try talking in <#%s> and see what happens!", reflectChannel.ID))
		if err != nil {

			log.Errorf("unable to edit a message. error: %v", err)

		}
		return

	})

}
