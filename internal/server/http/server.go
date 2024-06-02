package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"OTUS_hws/Anti-BruteForce/internal/antibrutforce"
	"OTUS_hws/Anti-BruteForce/internal/config"
)

var HandlersPaths = []string{"/addWhiteIp", "/addBlackIp", "/deleteWhiteIP", "/deleteBlackIP", "/clearBucket"}

type AntiBrutForce interface {
	LoadCertainedIps(filePath string) error
	CheckRequest(ip string, login string, password string) (bool, error)
	CheckLogin(login string) (bool, error)
	CheckPassword(password string) (bool, error)
	CheckIP(ip string) (bool, error)
	CheckIPInList(ip net.IP) (passed bool, isFound bool, err error)

	AddToList(cidr string, passed bool) error
	DeleteFromList(cidr string) error

	ClearLoginBuckets(login string) error
	ClearIPBuckets(ip string) error
	ClearOldBuckets()
	ClearAllBuckets()
}
type Server struct {
	Serv *http.Server
	Abf  AntiBrutForce
}

func New(abf *antibrutforce.AntiBrutForce, conf *config.Config) *Server {
	server := Server{
		Abf: abf,
	}
	server.Serv = &http.Server{
		Addr:              conf.Service.Host + ":" + fmt.Sprint(conf.Service.Port),
		Handler:           server.routes(),
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	return &server
}

func (s *Server) Start() error {
	middleware := s.CheckRequest(s.Serv.Handler)
	err := http.ListenAndServe(s.Serv.Addr, middleware)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server start error: ", err.Error())
			return err
		}
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.Abf.ClearAllBuckets()
	err := s.Serv.Shutdown(ctx)
	if err != nil {
		fmt.Println("server shutdown error: " + err.Error())
		return err
	}
	fmt.Println("Server has stopped")

	return nil
}
