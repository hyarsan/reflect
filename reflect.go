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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
	// externals
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/superwhiskers/harmony"
)

var (
	config       configuration
	handler      *harmony.CommandHandler
	mentionRegex *regexp.Regexp
	escapeRegex  *regexp.Regexp
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})

}

// the main function for this bot
func main() {

	runtime.GOMAXPROCS(1000)

	file, err := os.OpenFile("reflect.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		log.Warnf("unable to open logfile. falling back to stdout-only. error: %v", err)

	} else {

		defer file.Close()
		log.SetOutput(io.MultiWriter(os.Stdout, file))

	}

	cfgByte, err := ioutil.ReadFile("config.json")
	if err != nil {

		log.Fatalf("unable to read config file. error: %v", err)

	}

	err = json.Unmarshal(cfgByte, &config)
	if err != nil {

		log.Fatalf("unable to parse config file as json. error: %v", err)

	}

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", config.Token))
	if err != nil {

		log.Fatalf("unable to make a new discordgo session object. error: %v", err)

	}

	mentionRegex = regexp.MustCompile("\\@everyone|\\@here")

	escapeRegex = regexp.MustCompile("\\`|\\*|\\_|\\\\")

	handler = harmony.New("r~", true)

	registerUtilityCommands()
	registerModCommands()
	registerEvtHandlers()

	dg.AddHandler(handler.OnMessage)

	dg.AddHandler(onReady)

	err = dg.Open()
	if err != nil {

		log.Fatalf("[err]: unable to initiate a websocket session. error: %v", err)

	}

	log.Printf("press ctrl-c to stop the bot...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

	cfgByte, err = json.MarshalIndent(config, "", "	")
	if err != nil {

		log.Fatalf("[err]: unable to convert the config back to json. error: %v", err)

	}

	err = ioutil.WriteFile("config.json", cfgByte, 0644)
	if err != nil {

		log.Fatalf("[err]: unable to output config back to file. error: %v", err)

	}

}
