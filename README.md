# Karejone-Training-Bot

## Build

install [Go](https://golang.org/)

import [discordgo](github.com/bwmarrin/discordgo) package

build with version number

`go build -ldflags "-X main.version=1.0.0"`

Cross compile

`env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=1.0.0"`

## Run

`./<build> -t <bot_token>`
