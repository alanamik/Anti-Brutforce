package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	server "github.com/alanamik/Anti-Brutforce/internal/server/http"
	"github.com/stretchr/testify/require"
)

const (
	limitLogin    = 6
	limitPassword = 8
	limitIP       = 10
)

type ResponseSuccess struct {
	Status     string `json:"status"`
	StatusCode int    `json:"code"`
}

type ResponseError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"code"`
}

//nolint:gocognit
func TestABFServer(t *testing.T) {
	inputCheckReqsIn := []server.CheckRequestIn{
		{IP: "23.44.135.90", Login: "user1", Password: "ghfgh"},
		{IP: "100.44.100.50", Login: "user2", Password: "dfvdfv"},
		{IP: "45.44.43.56", Login: "user3", Password: "rtgtgrebc"},
		{IP: "45.44.43.54", Login: "user4", Password: "vcbfg"},
	}
	t.Run("ClearBucketByLoginAndIPTest", func(t *testing.T) {
		client := &http.Client{}
		for _, r := range inputCheckReqsIn {
			d := server.ClearBucketIn{
				IP:    r.IP,
				Login: r.Login,
			}
			data, err := json.Marshal(d)
			require.NoError(t, err)
			buf := bytes.NewBuffer(data)
			req, err := http.NewRequestWithContext(context.Background(), "DELETE", "http://0.0.0.0:8000/clearBucket", buf)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
		}
		t.Run("MiddlewareTest", func(t *testing.T) {
			// Check successfully passed request
			for i := 0; i <= limitLogin; i++ {
				resp, err := RequestAuth(inputCheckReqsIn[0], client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				require.Equal(t, 200, resp.StatusCode)
				time.Sleep(10 * time.Second)
			}
			// Check did not passed login
			for i := 0; i <= limitLogin+3; i++ {
				resp, err := RequestAuth(inputCheckReqsIn[1], client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				if i <= limitLogin {
					require.Equal(t, 200, resp.StatusCode)
					time.Sleep(5 * time.Second)
				} else {
					require.Equal(t, 400, resp.StatusCode)
				}
			}
			// Check did not passed password
			for i := 0; i <= limitPassword+3; i++ {
				if i == limitLogin {
					inputCheckReqsIn[2].Login = "xcvxv"
				}
				resp, err := RequestAuth(inputCheckReqsIn[2], client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				if i <= limitPassword {
					require.Equal(t, 200, resp.StatusCode)
					time.Sleep(5 * time.Second)
				} else {
					require.Equal(t, 400, resp.StatusCode)
				}
			}
		})
		t.Run("ListFunctionsTest", func(t *testing.T) {
			// Whitelist add
			whiteCidr := "33.87.120.175/20"
			resp, err := RequestAddInList(whiteCidr, client, true)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
			whiteReq := server.CheckRequestIn{
				IP:       "33.87.120.175",
				Login:    "goodHuman",
				Password: "fghgngn",
			}
			for i := 0; i <= limitIP+3; i++ {
				resp, err := RequestAuth(whiteReq, client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				require.Equal(t, 200, resp.StatusCode)
			}
			// Whitelist delete
			resp, err = RequestDeleteFromList(whiteCidr, client, true)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
			for i := 0; i <= limitIP+3; i++ {
				if i == limitLogin {
					whiteReq.Login = "zero"
					whiteReq.Password = "0000"
				}
				resp, err := RequestAuth(whiteReq, client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				if i <= limitIP {
					require.Equal(t, 200, resp.StatusCode)
				} else {
					require.Equal(t, 400, resp.StatusCode)
				}
			}
			// Blacklist add
			blackCidr := "145.67.100.175/20"
			resp, err = RequestAddInList(blackCidr, client, false)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
			blackReq := server.CheckRequestIn{
				IP:       "145.67.100.175",
				Login:    "badHuman",
				Password: "okiokl'p",
			}
			// Should be blocked in first time
			resp, err = RequestAuth(blackReq, client)
			if err != nil {
				t.Fatalf("Failed Request(): %s", err)
			}
			defer resp.Body.Close()
			require.Equal(t, 400, resp.StatusCode)
			// Whitelist delete
			resp, err = RequestDeleteFromList(blackCidr, client, false)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)

			for i := 0; i <= limitIP+3; i++ {
				if i == limitLogin {
					blackReq.Login = "zero00000"
					blackReq.Password = "000000000"
				}
				resp, err := RequestAuth(blackReq, client)
				if err != nil {
					t.Fatalf("Failed Request(): %s", err)
				}
				defer resp.Body.Close()
				if i <= limitIP {
					require.Equal(t, 200, resp.StatusCode)
				} else {
					require.Equal(t, 400, resp.StatusCode)
				}
			}
		})
	})
	t.SkipNow()
}

func RequestAuth(val server.CheckRequestIn, client *http.Client) (*http.Response, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://0.0.0.0:8000/hello", buf)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

func RequestAddInList(cidr string, client *http.Client, passed bool) (*http.Response, error) {
	val := server.AddIPIn{
		Cidr: cidr,
	}
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	var resp *http.Response
	var req *http.Request
	if passed {
		req, err = http.NewRequestWithContext(context.Background(), "POST", "http://0.0.0.0:8000/addWhiteIp", buf)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequestWithContext(context.Background(), "POST", "http://0.0.0.0:8000/addBlackIp", buf)
		if err != nil {
			return nil, err
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

func RequestDeleteFromList(cidr string, client *http.Client, passed bool) (*http.Response, error) {
	val := server.DeleteIPIn{
		Cidr: cidr,
	}
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	var resp *http.Response
	var req *http.Request
	if passed {
		req, err = http.NewRequestWithContext(context.Background(), "DELETE", "http://0.0.0.0:8000/deleteWhiteIP", buf)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequestWithContext(context.Background(), "DELETE", "http://0.0.0.0:8000/deleteBlackIP", buf)
		if err != nil {
			return nil, err
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
