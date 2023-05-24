package main

import (
	"fmt"
	"time"

	"github.com/h3ll0kitt1/tg_bot_for_indecisive/config"
	"github.com/h3ll0kitt1/tg_bot_for_indecisive/https/telegram"
	"github.com/h3ll0kitt1/tg_bot_for_indecisive/interaction/inter"
	"github.com/h3ll0kitt1/tg_bot_for_indecisive/storage/redis"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load configuration: ", err)
	}
	token := cfg.AccessToken

	var c telegram.Counter
	redisDb, err := redis.New("localhost", "6379")

	if err != nil {
		fmt.Println("new database: ", err)
	}

	for {

		ch := make(chan *telegram.UpdatesResponse, 1)

		update, err := c.NextUpdate(1, token)
		if err != nil {
			fmt.Println("request update: ", err)
			return
		}

		if update == nil {
			time.Sleep(time.Second)
			continue
		}

		ch <- update

		updateToProcess := <-ch
		updateToSent, err := inter.Process(redisDb, *updateToProcess)
		if updateToSent == nil {
			fmt.Println("process message: ", err)
			return
		}

		_, err = telegram.SendMessageToChat(*updateToSent, token)
		if err != nil {
			fmt.Println("send message to chat: ", err)
			return
		}
	}
}
