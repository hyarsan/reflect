/*

    reflect - link discord servers together like never before
    Copyright (C) 2018  superwhiskers <whiskerdev@protonmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	// internals
	"math/rand"
	"strings"
	// externals
	"github.com/bwmarrin/discordgo"
)

// check if a user is a webhook
func isWebhook(s *discordgo.Session, serverID string, user *discordgo.User) (bool, error) {

	hooks, err := s.GuildWebhooks(serverID)
	if err != nil {

		return false, err

	}

	for _, hook := range hooks {

		if hook.ID == user.ID {

			return true, nil

		}

	}

	return false, nil

}


// gets all of the guilds the client is in but instead it's only a list of the names
func allGuildNames() []string {

	var finalisedGuildNameList []string

	for _, guild := range userGuilds {

		finalisedGuildNameList = append(finalisedGuildNameList, guild.Name)

	}

	return finalisedGuildNameList

}

// shuffle a slice of strings
func shuffleStringSlice(slice []string) []string {

	retSlice := make([]string, len(slice))

	perm := rand.Perm(len(slice))

	for i, v := range perm {

		retSlice[v] = slice[i]

	}

	return retSlice

}

// check if something is a valid command
func command(prefix, name, message string) bool {

	return strings.Join([]string{prefix, name}, "") == message

}

// get a channel by a name
func guildChannelByName(s *discordgo.Session, guildID, name string) ([]*discordgo.Channel, error) {

	channels, err := s.GuildChannels(guildID)
	if err != nil {

		return []*discordgo.Channel{&discordgo.Channel{}}, err

	}

	retChannels := []*discordgo.Channel{}

	for _, channel := range channels {

		if channel.Name == name {

			retChannels = append(retChannels, channel)

		}

	}

	return retChannels, nil

}

// find all webhooks with a name in a channel
func channelWebhooksByName(s *discordgo.Session, channelID, name string) ([]*discordgo.Webhook, error){

	retHooks := []*discordgo.Webhook{}

	hooks, err := s.ChannelWebhooks(channelID)
	if err != nil {

		return []*discordgo.Webhook{}, err

	}

	for _, hook := range hooks {

		if hook.Name == name {

			retHooks = append(retHooks, hook)

		}

	}

	return retHooks, nil

}

// check if a string is in a list
func inList(obj string, list []string) bool {

	for _, val := range list {

		if val == obj {

			return true

		}

	}

	return false

}
