package antibrutforce

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"OTUS_hws/Anti-BruteForce/internal/config"

	"gopkg.in/yaml.v2"
)

const (
	BucketRangeTime  time.Duration = time.Minute * 1
	BucketLivingTime time.Duration = time.Minute * 5
)

var (
	ErrNoSuchLogin = errors.New("no such login")
	ErrNoSuchIP    = errors.New("no such IP")
	ErrIPInListYet = errors.New("IP in the list yet")
)

type AntiBrutForce struct {
	LimitIP       int
	LimitLogin    int
	LimitPassword int

	ClientsLogins    map[string]Bucket
	ClientsPasswords map[string]Bucket
	ClientsIPs       map[string]Bucket

	CertainedIps map[string]IPNet // true - whitelist, false - blacklist
}

type Bucket struct {
	RequestsPerMinutes int
	Timer              time.Time // время первого запроса после обнуления количества попыток
}

type IPNetIn struct {
	Cidr   string `json:"cidr"`
	Passed bool   `json:"passed"`
}

type IPNet struct {
	Cidr   string
	Mask   *net.IPNet
	Passed bool
}

func New(conf *config.Config) (*AntiBrutForce, error) {
	abf := &AntiBrutForce{
		LimitIP:       conf.Parameters.LimitIP,
		LimitLogin:    conf.Parameters.LimitLogin,
		LimitPassword: conf.Parameters.LimitPassword,
	}
	abf.ClientsLogins = make(map[string]Bucket, 0)
	abf.ClientsPasswords = make(map[string]Bucket, 0)
	abf.ClientsIPs = make(map[string]Bucket, 0)

	abf.CertainedIps = make(map[string]IPNet, 0)
	err := abf.LoadCertainedIps(conf.IPs.Path)
	if err != nil {
		return nil, err
	}
	return abf, nil
}

func (abf *AntiBrutForce) LoadCertainedIps(filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	out := make([]IPNetIn, 0)
	err = yaml.Unmarshal(b, &out)
	if err != nil {
		return err
	}

	for _, addr := range out {
		abf.AddToList(addr.Cidr, addr.Passed)
	}

	return nil
}

func (abf *AntiBrutForce) CheckRequest(ctx context.Context, ip string, login string, password string) (bool, error) {
	// проверяем сначала IP, если есть в листах, то прерываем проверку
	addr := net.ParseIP(ip)
	val, isFound, err := abf.CheckIPInList(addr)
	if err != nil {
		return false, err
	}
	if isFound {
		return val, nil
	}
	passed, err := abf.CheckIP(ip)
	if err != nil || !passed {
		return false, err
	}
	passed, err = abf.CheckPassword(password)
	if err != nil || !passed {
		return false, err
	}
	passed, err = abf.CheckLogin(login)
	if err != nil || !passed {
		return false, err
	}

	return passed, nil
}

func (abf *AntiBrutForce) CheckLogin(login string) (bool, error) {
	if _, ok := abf.ClientsLogins[login]; !ok {
		fmt.Println("NO L IN MAP")
		client := Bucket{
			RequestsPerMinutes: 1,
			Timer:              time.Now(),
		}
		abf.ClientsLogins[login] = client
		return true, nil
	}
	client := abf.ClientsLogins[login]
	if time.Since(client.Timer) > BucketRangeTime {
		fmt.Println("L ZERO")
		client.Timer = time.Now()
		client.RequestsPerMinutes = 1
		abf.ClientsLogins[login] = client
		return true, nil
	}
	if (time.Since(client.Timer) < BucketRangeTime) && client.RequestsPerMinutes <= abf.LimitLogin {
		fmt.Println("L PASSED - ", client.RequestsPerMinutes)
		client.RequestsPerMinutes++
		abf.ClientsLogins[login] = client
		return true, nil
	}

	fmt.Println("L END")
	return false, nil
}

func (abf *AntiBrutForce) CheckPassword(password string) (bool, error) {
	if _, ok := abf.ClientsPasswords[password]; !ok {
		fmt.Println("NO P IN MAP")
		client := Bucket{
			RequestsPerMinutes: 1,
			Timer:              time.Now(),
		}
		abf.ClientsPasswords[password] = client
		return true, nil
	}

	client := abf.ClientsPasswords[password]
	if time.Since(client.Timer) > BucketRangeTime {
		fmt.Println("ZERO P")
		client.Timer = time.Now()
		client.RequestsPerMinutes = 1
		abf.ClientsPasswords[password] = client
		return true, nil
	}
	if (time.Since(client.Timer) < BucketRangeTime) && client.RequestsPerMinutes <= abf.LimitPassword {
		fmt.Println("PASSED P- ", client.RequestsPerMinutes)
		client.RequestsPerMinutes++
		abf.ClientsPasswords[password] = client
		return true, nil
	}

	fmt.Println("P END")
	return false, nil
}

func (abf *AntiBrutForce) CheckIP(ip string) (bool, error) {
	if _, ok := abf.ClientsIPs[ip]; !ok {
		fmt.Println("NO IP IN MAP")
		client := Bucket{
			RequestsPerMinutes: 1,
			Timer:              time.Now(),
		}
		abf.ClientsIPs[ip] = client
		return true, nil
	}
	client := abf.ClientsIPs[ip]
	if time.Since(client.Timer) > BucketRangeTime {
		fmt.Println("ZERO IP")
		client.Timer = time.Now()
		client.RequestsPerMinutes = 1
		abf.ClientsIPs[ip] = client
		return true, nil
	}
	if (time.Since(client.Timer) < BucketRangeTime) && client.RequestsPerMinutes <= abf.LimitIP {
		fmt.Println("PASSED IP- ", client.RequestsPerMinutes)
		client.RequestsPerMinutes++
		abf.ClientsIPs[ip] = client
		return true, nil
	}

	fmt.Println("IP END")
	return false, nil
}

func (abf *AntiBrutForce) CheckIPInList(ip net.IP) (passed bool, isFound bool, err error) {
	if _, ok := abf.CertainedIps[ip.String()]; !ok {
		return false, false, nil
	}
	addr := abf.CertainedIps[ip.String()]
	return addr.Passed, true, nil
}

func (abf *AntiBrutForce) AddToList(cidr string, passed bool) error {
	addr, ipNet, err := net.ParseCIDR(cidr)
	fmt.Println(addr, ipNet, err)
	_, isFound, err := abf.CheckIPInList(addr)
	if err != nil {
		return err
	}
	if isFound {
		return ErrIPInListYet
	}

	address := IPNet{
		Cidr:   cidr,
		Mask:   ipNet,
		Passed: passed,
	}
	abf.CertainedIps[addr.String()] = address
	return nil
}

func (abf *AntiBrutForce) DeleteFromList(cidr string) error {
	addr, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	_, isFound, err := abf.CheckIPInList(addr)
	if err != nil {
		return err
	}
	if !isFound {
		return ErrNoSuchIP
	}

	delete(abf.CertainedIps, addr.String())
	return nil
}

func (abf *AntiBrutForce) ClearOldBuckets() {
	for c, b := range abf.ClientsLogins {
		if time.Since(b.Timer) > BucketLivingTime {
			delete(abf.ClientsLogins, c)
		}
	}
	for c, b := range abf.ClientsIPs {
		if time.Since(b.Timer) > BucketLivingTime {
			delete(abf.ClientsIPs, c)
		}
	}
	for c, b := range abf.ClientsPasswords {
		if time.Since(b.Timer) > BucketLivingTime {
			delete(abf.ClientsPasswords, c)
		}
	}
}

func (abf *AntiBrutForce) ClearLoginBuckets(login string) error {
	if _, ok := abf.ClientsLogins[login]; !ok {
		return ErrNoSuchLogin
	}
	delete(abf.ClientsLogins, login)
	return nil
}

func (abf *AntiBrutForce) ClearIPBuckets(ip string) error {
	if _, ok := abf.ClientsIPs[ip]; !ok {
		return ErrNoSuchIP
	}
	delete(abf.ClientsIPs, ip)
	return nil
}

func (abf *AntiBrutForce) ClearAllBuckets() {
	for c := range abf.ClientsLogins {
		delete(abf.ClientsLogins, c)
	}
	for c := range abf.ClientsIPs {
		delete(abf.ClientsIPs, c)
	}
	for c := range abf.ClientsPasswords {
		delete(abf.ClientsPasswords, c)
	}
}
