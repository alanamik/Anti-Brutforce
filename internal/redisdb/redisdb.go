package redisdb

import (
	"context"
	"errors"
	"fmt"

	"OTUS_hws/Anti-BruteForce/internal/config"

	"github.com/redis/go-redis/v9"
)

const (
	Blacklist = "blacklist"
	Whitelist = "whitelist"
)

var (
	ErrListDoesNotExist       = errors.New("the list does not exist")
	ErrIPInListYet            = errors.New("the IP already added in a list")
	ErrInternalServiceRedisDB = errors.New("failed request to RedisDB")
)

type RedisClient struct {
	Client *redis.Client
}

func NewClient(con config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     con.Redis.Address,
		Password: "password",
		DB:       con.Redis.DB,
	})

	return &RedisClient{
		Client: client,
	}
}

func checkListName(name string) error {
	if name != Whitelist && name != Blacklist {
		fmt.Println(name)
		return ErrListDoesNotExist
	}
	return nil
}

func (r *RedisClient) AddToList(ctx context.Context, ip string, list string) error {
	isInList, err := r.CheckInList(ctx, ip, list)
	if err != nil {
		return err
	}
	if isInList {
		return ErrIPInListYet
	}

	addedIP, err := r.Client.LPush(ctx, list, ip).Result()
	if err != nil {
		fmt.Println(addedIP)
		return ErrInternalServiceRedisDB
	}

	fmt.Println(r.GetAllIPFromList(ctx, list))
	return nil
}

func (r *RedisClient) DeleteFromList(ctx context.Context, ip string, list string) error {
	err := checkListName(list)
	if err != nil {
		return err
	}

	deletedIP, err := r.Client.LRem(ctx, list, 1, ip).Result()
	if err != nil {
		fmt.Println(deletedIP)
		return ErrInternalServiceRedisDB
	}

	return nil
}

func (r *RedisClient) GetAllIPFromList(ctx context.Context, list string) ([]string, error) {
	err := checkListName(list)
	if err != nil {
		return nil, err
	}

	ips, err := r.Client.LRange(ctx, list, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// TODO: удалить вывод потом
	for i, val := range ips {
		fmt.Println(i, "-", val)
	}

	return ips, nil
}

func (r *RedisClient) CheckInList(ctx context.Context, ip string, list string) (bool, error) {
	err := checkListName(list)
	if err != nil {
		return false, err
	}

	ips, err := r.Client.LRange(ctx, list, 0, -1).Result()
	if err != nil {
		return false, err
	}
	isInList := false

	for _, val := range ips {
		if val == ip {
			isInList = true
		}
	}

	return isInList, nil
}

// Buckets functions
/*func (r *RedisClient) IncrementBucketValue(ctx context.Context, key string) {

}

func (r *RedisClient) SetBucketValue(ctx context.Context, key string, value int) {

}

func (r *RedisClient) GetBucketValue(ctx context.Context, key string) {

}
*/
