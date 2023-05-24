package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func New(host string, port string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // Addr: "localhost:6379"
		Password: "",
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("Redis ping: %w", err)
	}
	return &Redis{client}, nil
}

func (r Redis) Save(chatId int, book string) (bool, error) {
	ok, err := r.client.SAdd(strconv.Itoa(chatId), book).Result()
	if err != nil {
		return false, fmt.Errorf("Redis SAdd: %w", err)
	}

	if ok == 0 {
		return false, nil
	}
	return true, nil
}

func (r Redis) Delete(chatId int, book string) (bool, error) {
	ok, err := r.client.SRem(strconv.Itoa(chatId), book).Result()
	if err != nil {
		return false, fmt.Errorf("Redis SRem: %w", err)
	}

	if ok == 0 {
		return false, nil
	}
	return true, nil
}

func (r Redis) Exists(chatId int, book string) (bool, error) {
	ok, err := r.client.SIsMember(strconv.Itoa(chatId), book).Result()
	if err != nil {
		return false, fmt.Errorf("Redis sIsMember: %w", err)
	}
	return ok, nil
}

func (r Redis) LenNotZero(chatId int) (bool, error) {
	len, err := r.client.SCard(strconv.Itoa(chatId)).Result()
	if err != nil {
		return false, fmt.Errorf("Redis sCard: %w", err)
	}

	if len == 0 {
		return false, nil
	}
	return true, nil
}

func (r Redis) Print(chatId int) ([]string, error) {
	list, err := r.client.SMembers(strconv.Itoa(chatId)).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis SMembers: %w", err)
	}
	return list, nil
}

func (r Redis) Rand(chatId int) (string, error) {
	random, err := r.client.SRandMember(strconv.Itoa(chatId)).Result()
	if err != nil {
		return "", fmt.Errorf("Redis SRandMember: %w", err)
	}
	return random, nil
}
