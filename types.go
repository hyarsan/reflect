/*

types.go -
the types used in the bot

credits:
  - @hyarsan#3653 - original bot creator

license: gnu agplv3

*/

package main

// structure for the config
type configuration struct {
	Token  string   `json:"token"`
	Prefix string   `json:"prefix"`
	Bans   []string `json:"bans"`
	Owners []string `json:"owners"`
}
