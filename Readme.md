# Discord Bot for Minecraft Whitelist

A simple Discord Bot designed to allow your server members to add themselves to the whitelist of a configured Minecraft server. Each member assigns one MC username to their discord username. Currently, there is no verification for MC usernames so deploy only on trusted servers with trustable members.

# Discord Usage

The bot adds 3 slash functions:

`/whitelist-add <MC username>` assigns MC username to the calling Discord member

`/whitelist-remove` removes the assignment from the calling Discord member

`/whitelist-show` prints whether an assignment to the calling Discord member exists and its value

# Installation

## Prerequisites

Create a new Discord App using the developer portal and invite the bot to your server with at least the following scopes:
- bot
- application.commands
and the following permissions:
- Send Messages

## Installation

1. Download `discord-minecraft-whitelist-bot` binary and `template.config.json` from from [releases](https://git.tronnet.net/alu/discord-minecraft-whitelist-bot/releases)
2. Rename `template.config.json` to `config.json` and modify:
    - app-id: Discord App ID from the developer portal
    - guild-id: Server ID which can be obtained by right clicking the server in Discord and clicking `Copy Server ID`
    - token: Discord App Token from the developer portal
    - mc-rcon: `IP:port` or `hostname:port` formatted ip/hostanme and port number of the MC server rcon
    - mc-rcon-password: Password for the MC server rcon
3. Run the binary