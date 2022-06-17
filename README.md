# Phish.in Discord Bot

Discord bot that fetches setlist and song data from the phish.in API.

## Building

```sh
docker build -t mguterl/phishin-discord-bot .
```

## Running

```sh
docker run -e DISCORD_TOKEN=secret PHISHIN_TOKEN=secret mguterl/phishin-discord-bot
```

## Publishing

```sh
docker push mguterl/phishin-discord-bot
```
