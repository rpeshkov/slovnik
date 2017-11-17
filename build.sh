#!/bin/sh

# Building api server
docker build -t slovnik_api -f cmd/api-server/Dockerfile .

# Building telegram bot
docker build -t telegram_bot -f cmd/telegram-bot/Dockerfile .

