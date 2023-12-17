package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Luna-devv/nekostic/aggregate"
	"github.com/Luna-devv/nekostic/api"
	"github.com/Luna-devv/nekostic/config"
	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func task() {
	conf := config.Get()

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		Username: conf.Redis.Username,
		DB:       conf.Redis.Db,
	})

	aggregate.Commands(client)
	aggregate.Tags(client)

	client.Close()
}

func main() {
	conf := config.Get()

	if runtime.GOOS == "windows" {
		task()
		return
	}

	fmt.Println("Starting cron job...")

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Cron("0 23 * * *").Do(task)

	scheduler.StartAsync()
	api.New(conf)
}
