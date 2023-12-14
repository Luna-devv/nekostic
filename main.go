package main

import (
	"fmt"
	"time"

	"github.com/Luna-devv/nekostic/aggregate"
	"github.com/Luna-devv/nekostic/api"
	"github.com/Luna-devv/nekostic/config"
	"github.com/go-co-op/gocron"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conf := config.Get()

	fmt.Println("Starting cron job...")

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Cron("0 23 * * *").Do(aggregate.Run)

	scheduler.StartAsync()
	api.New(conf)
}
