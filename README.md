# Phish.in Discord Bot

Discord bot that fetches setlist and song data from the phish.in API.

## Building

```sh
docker build -t phishin-discord-bot .
```

## Running

```sh
docker run -e DISCORD_TOKEN=secret PHISHIN_TOKEN=secret phishin-discord-bot
```
