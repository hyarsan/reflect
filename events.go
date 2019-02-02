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
	"fmt"
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

	// on message handler
	handler.OnMessageHandler = func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Type != discordgo.MessageTypeDefault {

			return

		}

		if m.Author.ID == reflectUser.ID {

			return

		}

		msgChannel, err := s.Channel(m.ChannelID)
		if err != nil {

			log.Errorf("unable to fetch channel object. error: %v", err)
			return

		}

		if msgChannel.Name != config.ChannelName {

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

		username = string(escapeRegex.ReplaceAllFunc([]byte(username), func(in []byte) []byte {

			return append([]byte("\\"), in...)

		}))
		username = string(mentionRegex.ReplaceAllFunc([]byte(username), func(in []byte) []byte {

			return append([]byte("<at>"), in[1:]...)

		}))

		content := m.ContentWithMentionsReplaced()
		content = string(mentionRegex.ReplaceAllFunc([]byte(content), func(in []byte) []byte {

			return append([]byte("<at>"), in[1:]...)

		}))

		if !inList(m.Author.ID, config.Owners) {

			content = fmt.Sprintf("**%s**: %s", username, content)

		} else {

			content = fmt.Sprintf("**%s** **__(admin)__**: %s", username, content)

		}

		if len(content) > 2000 {

			return

		}

		var messageData *discordgo.MessageSend
		if len(m.Embeds) > 0 {

			messageData = &discordgo.MessageSend{
				Content: content,
				Embed:   m.Embeds[0],
			}

		} else {

			messageData = &discordgo.MessageSend{
				Content: content,
			}

		}

		for _, guild := range userGuilds {

			go backgroundMessageSend(s, guild, msgChannel, messageData)

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
