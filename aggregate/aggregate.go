package aggregate

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Luna-devv/nekostic/config"
	"github.com/redis/go-redis/v9"
)

type CommandEvent struct {
	Event  string `json:"event"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type AggregatedEvent struct {
	Event string `json:"event"`
	Name  string `json:"name"`
	Uses  int    `json:"uses"`
	Users int    `json:"users"`
}

func Run() {
	start := makeTimestamp()

	conf := config.Get()
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		Username: conf.Redis.Username,
		DB:       conf.Redis.Db,
	})

	var commandEvents []CommandEvent

	nameMap := make(map[string]map[string]bool)
	aggregatedEvents := make(map[string]*AggregatedEvent)

	keys, err := client.Keys(ctx, "command-event:*").Result()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, key := range keys {
		hashExists, _ := client.Type(ctx, key).Result()
		if hashExists != "hash" {
			continue
		}

		fieldsVals, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error retrieving hash values for key %s: %s\n", key, err)
			continue
		}

		commandEvent := CommandEvent{
			Event:  fieldsVals["event"],
			Name:   fieldsVals["name"],
			UserId: fieldsVals["userId"],
		}
		commandEvents = append(commandEvents, commandEvent)

		uid := fieldsVals["event"] + "-" + fieldsVals["name"]

		if _, ok := nameMap[uid]; !ok {
			nameMap[uid] = make(map[string]bool)
		}
		nameMap[uid][fieldsVals["userId"]] = true

		err = client.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("Error deleting key:", err)
		}

	}

	for _, ce := range commandEvents {
		uid := ce.Event + "-" + ce.Name

		if _, ok := aggregatedEvents[uid]; !ok {
			aggregatedEvents[uid] = &AggregatedEvent{
				Event: ce.Event,
				Name:  ce.Name,
			}
		}

		aggregatedEvents[uid].Uses++
	}

	for _, ae := range aggregatedEvents {
		uid := ae.Event + "-" + ae.Name
		aggregatedEvents[uid].Users = len(nameMap[uid])
	}

	for _, ae := range aggregatedEvents {
		uid := ae.Event + "-" + ae.Name

		err := client.HSet(ctx, "aggregated-command:"+make6DigitDay()+":"+uid, "event", ae.Event, "name", ae.Name, "uses", strconv.Itoa(ae.Uses), "users", strconv.Itoa(ae.Users)).Err()
		if err != nil {
			fmt.Println("Error storing struct in Redis hash:", err)
			continue
		}

	}

	client.Close()
	fmt.Printf("Cron job ran: %s with %d events in %d ms\n", make6DigitDay(), len(aggregatedEvents), makeTimestamp()-start)
}

func make6DigitDay() string {
	now := time.Now()

	return fmt.Sprintf("%02d", int64(now.Day())) + strconv.FormatInt(int64(now.Month()), 10) + strconv.FormatInt(int64(now.Year()), 10)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
