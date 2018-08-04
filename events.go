/*

events.go -
misc. event handlers

credits:
  - @hyarsan#3653 - original bot creator

license: gnu agplv3

*/

package main

import (
	// internals
	"time"
	// externals
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	reflectUser *discordgo.User
	userCache   = map[string][]*discordgo.User{}
)

// set up some harmony event handlers
func registerEvtHandlers() {

	handler.OnMessageHandler = func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Type != discordgo.MessageTypeDefault {

			return

		}

		msgChannel, err := s.Channel(m.ChannelID)
		if err != nil {

			log.Errorf("unable to fetch channel object. error: %v", err)
			return

		}

		if msgChannel.Name != "megachat" {

			return

		}

		isHook, err := isWebhook(s, msgChannel.GuildID, m.Author)
		if err != nil {

			log.Errorf("unable to check if the message sender was a webhook. error: %v", err)
			return

		}

		if isHook == true {

			return

		}

		for _, bannedID := range config.Bans {

			if m.Author.ID == bannedID {

				s.MessageReactionAdd(m.ChannelID, m.ID, "🔨")
				return

			}

		}

		member, err := s.GuildMember(msgChannel.GuildID, m.Author.ID)
		if err != nil {

			log.Errorf("unable to get guild member. error: %v", err)
			return

		}

		var username string
		if member.Nick == "" {

			username = m.Author.Username

		} else {

			username = member.Nick

		}

		foundMatch := false
		for _, user := range userCache[username] {

			if m.Author.ID == user.ID {

				foundMatch = true
				break

			}

		}

		if !foundMatch {

			userCache[username] = append(userCache[username], m.Author)

		}

		var webhookParams *discordgo.WebhookParams
		if len(m.Attachments) > 0 {

			if m.Content == "" {

				webhookParams = &discordgo.WebhookParams{
					Username:  username,
					AvatarURL: m.Author.AvatarURL(""),
					TTS:       false,
					File:      m.Attachments[0].URL,
					Embeds:    m.Embeds,
				}

			} else {

				webhookParams = &discordgo.WebhookParams{
					Content:   m.ContentWithMentionsReplaced(),
					Username:  username,
					AvatarURL: m.Author.AvatarURL(""),
					TTS:       false,
					File:      m.Attachments[0].URL,
					Embeds:    m.Embeds,
				}

			}

		} else {

			if m.Content == "" {

				webhookParams = &discordgo.WebhookParams{
					Username:  username,
					AvatarURL: m.Author.AvatarURL(""),
					TTS:       false,
					Embeds:    m.Embeds,
				}

			} else {

				webhookParams = &discordgo.WebhookParams{
					Content:   m.ContentWithMentionsReplaced(),
					Username:  username,
					AvatarURL: m.Author.AvatarURL(""),
					TTS:       false,
					Embeds:    m.Embeds,
				}

			}

		}

		for _, guild := range userGuilds {

			go backgroundWebhookExec(s, guild, msgChannel, webhookParams)

		}

	}

}

// the ready handler that is called when we are authenticated and such
func onReady(s *discordgo.Session, r *discordgo.Ready) {

	go backgroundGuildUpdater(s)
	go backgroundStatusUpdater(s)
	go backgroundConfigSave()

	time.Sleep(500 * time.Millisecond)

	log.Printf("logged in as %s on %d servers...", r.User.String(), len(userGuilds))

	reflectUser = r.User

}
