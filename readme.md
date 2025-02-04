# hsd

Home security daemon:

- uses door/window sensor information from zigbee2mqtt and sends alerts to telegram chat.
- has an API to enable/disable telegram alerts.

## Supported sensors

I only own Aqara Door and Window Sensor T1, so this is  the only supported sensor. Feel free to send PRs to support more devices.

## Note on zigbee2mqtt version compitability

Zigbee2mqtt 2.0 broke backward compitability, so, for now, only 1.x zigbee2mqtt branch is supported.

## How to run locally (for development)

docker-compose.yaml and Taskfile.yaml already have most of the stuff you need to start the app, just add your Telegram information to either `hsd.env` file, or just into Taskfile.yaml.

`hsd.env` file example:

```dotenv
TELEGRAM_BOT_TOKEN=VALUE
TELEGRAM_CHAT_ID=VALUE
```

It's ok for telegram chat id to be a negative integer (don't ask me why, ask telegram developers).

Then run with

```sh
docker compose up -d
task run
```

If everything is configured correctly, you should see `app started` in logs. Now we can do short end-to-end test by running:

```sh
mosquitto_pub --username hsd --pw hsd --topic 'zigbee2mqtt/door1/availability'  -m 'offline' -h localhost
```

Bot will send a telegram notification to your chat after this (or crashes with some error in logs that could help you debug things).

## How to deploy

Package into a docker image and ship to your server. Or use one of the tagged images from this repo artifacts.

## TODO 

- proper documentation
- zigbee2mqtt button support to enable/disable alerts
- zigbee2mqtt v2 support
