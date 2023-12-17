package aggregate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Luna-devv/nekostic/utils"
	"github.com/redis/go-redis/v9"
)

type TagEvent struct {
	GuildId   string `json:"guildId"`
	TagId     string `json:"tagId"`
	ChannelId string `json:"channelId"`
	UserId    string `json:"userId"`
}

type AggregatedTagEvent struct {
	GuildId string `json:"guildId"`
	TagId   string `json:"tagId"`
	Uses    int    `json:"uses"`
	Users   int    `json:"users"`
}

func Tags(client *redis.Client) {
	start := utils.MakeTimestamp()
	ctx := context.Background()

	var commandEvents []TagEvent

	nameMap := make(map[string]map[string]bool)
	aggregatedEvents := make(map[string]*AggregatedTagEvent)

	keys, err := client.Keys(ctx, "tag-event:*").Result()

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

		commandEvent := TagEvent{
			GuildId:   fieldsVals["guildId"],
			TagId:     fieldsVals["tagId"],
			ChannelId: fieldsVals["shannelId"],
			UserId:    fieldsVals["userId"],
		}
		commandEvents = append(commandEvents, commandEvent)

		if _, ok := nameMap[commandEvent.TagId]; !ok {
			nameMap[commandEvent.TagId] = make(map[string]bool)
		}
		nameMap[commandEvent.TagId][fieldsVals["userId"]] = true

		err = client.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("Error deleting key:", err)
		}

	}

	for _, ce := range commandEvents {

		if _, ok := aggregatedEvents[ce.TagId]; !ok {
			aggregatedEvents[ce.TagId] = &AggregatedTagEvent{
				GuildId: ce.GuildId,
				TagId:   ce.TagId,
			}
		}

		aggregatedEvents[ce.TagId].Uses++
	}

	for _, ae := range aggregatedEvents {
		aggregatedEvents[ae.TagId].Users = len(nameMap[ae.TagId])
	}

	for _, ae := range aggregatedEvents {

		err := client.HSet(ctx, "aggregated-tag:"+utils.Make6DigitDay()+":"+ae.TagId, "guildId", ae.GuildId, "tagId", ae.TagId, "uses", strconv.Itoa(ae.Uses), "users", strconv.Itoa(ae.Users)).Err()
		if err != nil {
			fmt.Println("Error storing struct in Redis hash:", err)
			continue
		}

	}

	fmt.Printf("Tag cron ran: %s with %d events in %d ms\n", utils.Make6DigitDay(), len(aggregatedEvents), utils.MakeTimestamp()-start)
}
