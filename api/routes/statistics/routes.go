package statistics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

type AggregatedEvent struct {
	Event    string `json:"event"`
	Name     string `json:"name"`
	Uses     int    `json:"uses"`
	Users    int    `json:"users"`
	Snapshot string `json:"snapshot"`
}

var httpClient = http.Client{Timeout: time.Second * 2}

func (i route) Statistics(c echo.Context) error {

	req, err := http.NewRequest(http.MethodGet, "http://10.0.0.50:7001/v1/users/@me", nil)
	if err != nil {
		fmt.Println("Api Error:", err)
		return err
	}

	req.Header.Set("Authorization", c.Request().Header.Get("authorization"))

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Api Error:", err)
		return err
	}

	if resp.StatusCode != 200 {
		return c.NoContent(resp.StatusCode)
	}

	ctx := context.Background()
	keys, err := i.Client.Keys(ctx, "aggregated-command:*").Result()

	if err != nil {
		fmt.Println("Api Error:", err)
		return err
	}

	var commands []AggregatedEvent

	for _, key := range keys {
		hashExists, _ := i.Client.Type(ctx, key).Result()
		if hashExists != "hash" {
			continue
		}

		fieldsVals, err := i.Client.HGetAll(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error retrieving hash values for key %s: %s\n", key, err)
			continue
		}

		uses, err := strconv.Atoi(fieldsVals["uses"])
		if err != nil {
			uses = 0
		}

		users, err := strconv.Atoi(fieldsVals["users"])
		if err != nil {
			users = 0
		}

		parts := strings.Split(key, ":")
		digits := parts[1]

		command := AggregatedEvent{
			Event:    fieldsVals["event"],
			Name:     fieldsVals["name"],
			Uses:     uses,
			Users:    users,
			Snapshot: fmt.Sprintf("%s-%s-%s", digits[4:8], digits[2:4], digits[0:2]),
		}
		commands = append(commands, command)

	}

	return c.JSON(200, commands)
}
