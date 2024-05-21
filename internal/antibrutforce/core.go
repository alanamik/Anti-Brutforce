package antibrutforce

import (
	"context"
	"fmt"
	"time"

	"OTUS_hws/Anti-BruteForce/internal/config"
	"OTUS_hws/Anti-BruteForce/internal/redisdb"
)

const BucketRangeTime time.Duration = time.Minute * 1

const BucketLivingTime time.Duration = time.Minute * 10

type AntiBrutForce struct {
	RedisServer *redisdb.RedisClient

	LimitIP       int
	LimitLogin    int
	LimitPassword int

	ClientsLogins    map[string]Bucket
	ClientsPasswords map[string]Bucket
	ClientsIPs       map[string]Bucket
}

type Bucket struct {
	RequestsPerMinutes int
	Timer              time.Time // время первого запроса после обнуления количества попыток
}

func New(r *redisdb.RedisClient, conf *config.Config) *AntiBrutForce {
	abf := &AntiBrutForce{
		RedisServer:   r,
		LimitIP:       conf.Parameters.LimitIP,
		LimitLogin:    conf.Parameters.LimitLogin,
		LimitPassword: conf.Parameters.LimitPassword,
	}
	abf.ClientsLogins = make(map[string]Bucket, 0)
	fmt.Println(abf)
	return abf
}

func (abf *AntiBrutForce) CheckRequest(ctx context.Context, ip string, login string, password string) (bool, error) {
	// проверяем сначала IP, если есть в листах, то прерываем проверку
	isInList, err := abf.checkInBlackList(ctx, ip)
	if err != nil {
		return false, err
	}
	if isInList {
		fmt.Println("BLACK LIST")
		return false, nil
	}

	isInList, err = abf.checkInWhiteList(ctx, ip)
	if err != nil {
		return false, err
	}
	if isInList {

		fmt.Println("WHITE LIST")
		return true, nil
	}

	passed, err := abf.CheckLogin(login)
	if err != nil {
		return false, err
	}

	return passed, nil
}

func (abf *AntiBrutForce) CheckLogin(login string) (bool, error) {
	// no in the map
	if _, ok := abf.ClientsLogins[login]; !ok {
		fmt.Println("NO IN MAP")
		client := Bucket{
			RequestsPerMinutes: 1,
			Timer:              time.Now(),
		}
		abf.ClientsLogins[login] = client
		return true, nil
	}
	// in the map
	client := abf.ClientsLogins[login]
	// Если с времени первого запроса после обнуления прошло больше минуты, то обнуляем время для клиента
	if time.Since(client.Timer) > BucketRangeTime {
		fmt.Println("ZERO")
		client.Timer = time.Now()
		client.RequestsPerMinutes = 1
		abf.ClientsLogins[login] = client
		return true, nil
	}
	// Если с времени первого запроса не прошло больше минуты И лимит попыток не превышен, то пропускаем
	if (time.Since(client.Timer) < BucketRangeTime) && client.RequestsPerMinutes <= abf.LimitLogin {
		fmt.Println("PASSED - ", client.RequestsPerMinutes)
		client.RequestsPerMinutes++
		abf.ClientsLogins[login] = client
		return true, nil
	}

	fmt.Println("END")
	return false, nil
}

func (abf *AntiBrutForce) CheckPassword(password string) {
}

func (abf *AntiBrutForce) checkInBlackList(ctx context.Context, ip string) (bool, error) {
	isInList, err := abf.RedisServer.CheckInList(ctx, ip, redisdb.Blacklist)
	if err != nil {
		return false, err
	}

	return isInList, nil
}

func (abf *AntiBrutForce) checkInWhiteList(ctx context.Context, ip string) (bool, error) {
	isInList, err := abf.RedisServer.CheckInList(ctx, ip, redisdb.Whitelist)
	if err != nil {
		return false, err
	}

	return isInList, nil
}

func (abf *AntiBrutForce) ClearOldLoginBuckets() {
	for c, b := range abf.ClientsLogins {
		if time.Since(b.Timer) > BucketLivingTime {
			delete(abf.ClientsLogins, c)
		}
	}
}

func (abf *AntiBrutForce) ClearAllBuckets() {
	for c := range abf.ClientsLogins {
		delete(abf.ClientsLogins, c)
	}
	for c := range abf.ClientsPasswords {
		delete(abf.ClientsPasswords, c)
	}

	for c := range abf.ClientsPasswords {
		delete(abf.ClientsPasswords, c)
	}
}
