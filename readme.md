[![](https://img.shields.io/discord/828676951023550495?color=5865F2&logo=discord&logoColor=white)](https://lunish.nl/support)
![](https://img.shields.io/github/repo-size/Luna-devv/nekostic?maxAge=3600)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I6AFVAP)

**⚠️ In development, breaking changes ⚠️**

## About
This is the analytics engine for the [wamellow.com](https://wamellow.com) discord bot, tracking command usage by command name, tracking and users.

If you need help using this, join **[our Discord Server](https://discord.com/invite/yYd6YKHQZH)**.

## Setup
Clone this repo with the following commands:

```bash
git clone https://github.com/Luna-devv/nekostic
```

Create a `.env` file and add the following values:
```env
REDIS_PW=""
REDIS_ADDR="127.0.0.1:6379"
REDIS_USR=""
PORT="3000"
```
Change the ports and/or add a user and/or password if your redis instance requires one. **This uses redis as a persistent database**, learn more on the [redis persistence documentation](https://redis.io/docs/management/persistence/) on how to make your deployment persistent.

## Deploy

Since docker is the best thing that exists for managing deployments, it's the thing we use. If you don't want to use docker, install the go programing language from [go.dev](https://go.dev) and run `go run .` or however you want to run it.

To build the docker container run
```bash
docker build -t nekostic .
```

To start the docker container (detached) run 
```bash
docker compose up -d
```