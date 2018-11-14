/*

    reflect - link discord servers together like never before
    Copyright (C) 2018  superwhiskers <whiskerdev@protonmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package main

// structure for the config
type configuration struct {
	Token       string   `json:"token"`
	Prefix      string   `json:"prefix"`
	ChannelName string   `json:"channel_name"`
	Bans        []string `json:"bans"`
	Owners      []string `json:"owners"`
}
