/*

modCommands.go -
the moderation command handlers for the reflect bot

credits:
  - @hyarsan#3653 - original bot creator

license: gnu agplv3

*/

package main

import (
	// internals
	"fmt"
	"strings"
	"strconv"
	// externals
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var err error

// function that registers all of the commands
func registerModCommands() {

	handler.AddCommand("ban", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		if !inList(m.Author.ID, config.Owners) {

			_, err = s.ChannelMessageSend(m.ChannelID, "You aren't an owner!")
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}

		escapedUserMentionThing := strings.Join(strings.Split(m.ContentWithMentionsReplaced(), " ")[1:], " ")
		escapedUserMentionThing = string(mentionRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("<at>"), in[1:]...)

		}))

		escapedUserMentionThing = string(escapeRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("\\"), in...)

		}))

		var userToBan *discordgo.User
		if len(m.Mentions) > 0 {

			userToBan = m.Mentions[0]

		} else if userSet, ok := userCache[strings.Join(args[0:], " ")]; ok {

			if len(userSet) == 1 {

				userToBan = userSet[0]

			} else {

				_, err = s.ChannelMessageSend(m.ChannelID, "Multiple users in usercache not supported yet.")
				if err != nil {

					log.Errorf("unable to send message. error: %v", err)

				}
				return

			}

		} else {

			userToBan, err = s.User(args[0])
			if err != nil {

				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid user ID, nick, or mention", escapedUserMentionThing))
				if err != nil {

					log.Errorf("unable to send message. error: %v", err)

				}
				return

			}

		}

		if inList(userToBan.ID, config.Owners) {

			_, err = s.ChannelMessageSend(m.ChannelID, "You can't ban an owner!")
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}

		if inList(userToBan.ID, config.Bans) {

			_, err = s.ChannelMessageSend(m.ChannelID, "That user is already banned!")
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}

		config.Bans = append(config.Bans, userToBan.ID)

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has been banned", userToBan.Username))
		if err != nil {

			log.Errorf("unable to send message. error: %v", err)
			return

		}

	})

	handler.AddCommand("unban", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		if !inList(m.Author.ID, config.Owners) {

			_, err = s.ChannelMessageSend(m.ChannelID, "You aren't an owner!")
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}
		
		escapedUserMentionThing := strings.Join(strings.Split(m.ContentWithMentionsReplaced(), " ")[1:], " ")
		escapedUserMentionThing = string(mentionRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("<at>"), in[1:]...)

		}))
		
		escapedUserMentionThing = string(escapeRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("\\"), in...)

		}))

		var userToUnban *discordgo.User
		if len(m.Mentions) > 0 {

			userToUnban = m.Mentions[0]

		} else {

			userToUnban, err = s.User(args[0])
			if err != nil {

				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid user ID", escapedUserMentionThing))
				if err != nil {

					log.Errorf("unable to send message. error: %v", err)

				}
				return

			}

		}

		if !inList(userToUnban.ID, config.Bans) {

			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not banned!", userToUnban.Username))
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}

		var index int
		for i, bannedUser := range config.Bans {

			if bannedUser == userToUnban.ID {

				index = i

			}

		}

		config.Bans[len(config.Bans)-1], config.Bans[index] = config.Bans[index], config.Bans[len(config.Bans)-1]
		config.Bans = config.Bans[:len(config.Bans)-1]

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has been unbanned", userToUnban.Username))
		if err != nil {

			log.Errorf("unable to send message. error: %v", err)
			return

		}

	})

	handler.AddCommand("bans", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		var banlistFields []*discordgo.MessageEmbedField

		if len(config.Bans) == 0 {

			banlistFields = []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "There are no bans yet. Impressive!",
					Value:  "...",
					Inline: true,
				},
			}

		} else {

			banlistFields = []*discordgo.MessageEmbedField{}
			for i, banID := range config.Bans {

				banlistFields = append(banlistFields, &discordgo.MessageEmbedField{
					Name:   strconv.Itoa(i+1),
					Value:  banID,
					Inline: true,
				})

			}

		}

		embed := discordgo.MessageEmbed{
			Title:       "Bans",
			Description: fmt.Sprintf("There are currently %d user(s) banned from Reflect.", len(config.Bans)),
			Color:       0x89da72,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://u.catgirl.host/4seqjs.png",
			},
			Fields: banlistFields,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
			},
		}

		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {

			log.Errorf("unable to send banlist. error: %v", err)
			return

		}

	})

	handler.AddCommand("user", false, func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		escapedUserMentionThing := strings.Join(strings.Split(m.ContentWithMentionsReplaced(), " ")[1:], " ")
		escapedUserMentionThing = string(mentionRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("<at>"), in[1:]...)

		}))
		
		escapedUserMentionThing = string(escapeRegex.ReplaceAllFunc([]byte(escapedUserMentionThing), func(in []byte) []byte {

			return append([]byte("\\"), in...)

		}))

		var targetUser *discordgo.User
		if len(args) == 0 {

			targetUser = m.Author

		} else if len(m.Mentions) > 0 {

			targetUser = m.Mentions[0]

		} else if userSet, ok := userCache[strings.Join(args[0:], " ")]; ok {

			if len(userSet) == 1 {

				targetUser = userSet[0]

			} else {

				embedFields := []*discordgo.MessageEmbedField{}
				for _, user := range userSet {

					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name: user.Username,
						Value: user.ID,
					})

				}

				embed := discordgo.MessageEmbed{
					Title:      "Usercache Lookup Results",
					Description: fmt.Sprintf("Showing every known user with the name %s", strings.Join(args[0:], " ")),
					Color:       0x89da72,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://u.catgirl.host/4seqjs.png",
					},
					Fields: embedFields,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
					},
				}

				_, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				if err != nil {

					log.Errorf("unable to send usercache lookup results. error: %v", err)
					return

				}
				return

			}

		} else {

			targetUser, err = s.User(args[0])
			if err != nil {

				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid user ID, nick, or mention", escapedUserMentionThing))
				if err != nil {

					log.Errorf("unable to send message. error: %v", err)

				}
				return

			}

		}

		if targetUser.ID == reflectUser.ID {

			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("What a hack. Just use %sinfo next time.", config.Prefix))
			if err != nil {

				log.Errorf("unable to send message. error: %v", err)

			}
			return

		}

		embed := discordgo.MessageEmbed{
			Title:       targetUser.String(),
			Description: fmt.Sprintf("Showing user information for %s", targetUser.Username),
			Color:       0x89da72,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: targetUser.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "Discord Information",
					Value: fmt.Sprintf(`**Username:** %s
**ID:** %s
**Bot:** %t`, targetUser.Username, targetUser.ID, targetUser.Bot),
				},
				&discordgo.MessageEmbedField{
					Name:  "Reflect Servers",
					Value: "Loading shared server list...",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
			},
		}

		ogMsg, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {

			log.Errorf("unable to send user information. error: %v", err)
			return

		}

		var targetUserSharedGuilds = []string{}
		for _, guild := range userGuilds {

			_, err := s.GuildMember(guild.ID, targetUser.ID)
			if err != nil {

				continue

			}

			targetUserSharedGuilds = append(targetUserSharedGuilds, guild.Name)

		}

		embed = discordgo.MessageEmbed{
			Title:       targetUser.String(),
			Description: fmt.Sprintf("Showing user information for %s", targetUser.Username),
			Color:       0x89da72,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: targetUser.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "Discord Information",
					Value: fmt.Sprintf(`**Username:** %s
**ID:** %s
**Bot:** %t`, targetUser.Username, targetUser.ID, targetUser.Bot),
				},
				&discordgo.MessageEmbedField{
					Name:  "Reflect Servers",
					Value: strings.Join(targetUserSharedGuilds, "\n"),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Built with ❤ by superwhiskers#3210 & hyarsan#3653",
			},
		}

		_, err = s.ChannelMessageEditEmbed(m.ChannelID, ogMsg.ID, &embed)
		if err != nil {

			log.Errorf("unable to edit user information message. error: %v", err)
			return

		}

	})

}
