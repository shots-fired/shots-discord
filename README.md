# Shots Discord

[![Build Status](https://travis-ci.org/shots-fired/shots-discord.svg?branch=master&service=github)](https://travis-ci.org/shots-fired/shots-discord)
[![Coverage Status](https://coveralls.io/repos/github/shots-fired/shots-discord/badge.svg?branch=master&service=github)](https://coveralls.io/github/shots-fired/shots-discord?branch=master)

Shots is a Discord bot. This project houses all the code responsible for interacting with Discord.

## Contributing

1. Install Go 1.11 or higher
2. `go get -u github.com/shots-fired/shots-discord`

## Running

The easiest way to run Shots is by cloning the `shots-deploy` project and using Docker Compose.

1. `git clone github.com/shots-fired/shots-deploy`
2. `cd shots-deploy`
3. `docker-compose build`
4. `docker-compose up`

## Environment variables

* SERVER_ADDRESS
* BOT_KEY
