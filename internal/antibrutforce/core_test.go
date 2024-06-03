package antibrutforce

import (
	"fmt"
	"net"
	"os"
	"testing"

	"OTUS_hws/Anti-BruteForce/internal/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAntibrutforceCore(t *testing.T) {
	os.Chdir("..")
	os.Chdir("..")
	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "[config.New()]")
		panic(err)
	}
	abf, _ := New(conf)
	t.Run("ListsFunctions", func(t *testing.T) {
		// check adding to lists
		inputBlockedIPs := []string{"145.92.137.88/20", "162.198.0.44/20", "162.198.0.157/27"}
		for _, ip := range inputBlockedIPs {
			err = abf.AddToList(ip, false)
			require.NoError(t, err)
		}
		inputPassedIPs := []string{"12.34.56.78/24", "192.168.0.101/24", "10.8.248.131/23"}
		for _, ip := range inputPassedIPs {
			err = abf.AddToList(ip, true)
			require.NoError(t, err)
		}
		// check adding IP, that is in a list
		err = abf.AddToList(inputBlockedIPs[2], false)
		require.ErrorIs(t, err, ErrIPInListYet)

		err = abf.AddToList(inputPassedIPs[2], true)
		require.ErrorIs(t, err, ErrIPInListYet)

		err = abf.SaveCertainedIpsInFile()
		require.NoError(t, err)
		err = abf.LoadCertainedIps()
		require.NoError(t, err)
		for _, i := range abf.CertainedIps {
			fmt.Println(i)
		}
		// check deleting from lists
		err = abf.DeleteFromList(inputBlockedIPs[1])
		require.NoError(t, err)

		_, isFound, _ := abf.CheckIPInList(net.ParseIP(inputBlockedIPs[2]))
		require.Equal(t, false, isFound)

		err = abf.DeleteFromList("78.34.201.90/21")
		require.ErrorIs(t, err, ErrNoSuchIP)

		// check another ips that not in the list
		ip := net.ParseIP("192.67.45.3")
		pass, found, err := abf.CheckIPInList(ip)
		require.NoError(t, err)
		require.Equal(t, false, pass)
		require.Equal(t, false, found)
	})
	t.Run("BucketFunctions", func(t *testing.T) {
		// Checking login buckets
		inputClientLogins := []string{
			"user1", "user2", "user3",
		}
		inputClientPassword := []string{
			"tert", "retvg", "rtefvgr",
		}
		inputClientIPs := []string{
			"12.38.58.78", "192.168.10.101", "10.80.248.131",
		}
		for i := 0; i <= abf.LimitLogin; i++ {
			pass, err := abf.CheckLogin(inputClientLogins[0])
			require.NoError(t, err)
			require.Equal(t, true, pass)
			pass, err = abf.CheckLogin(inputClientLogins[1])
			require.NoError(t, err)
			require.Equal(t, true, pass)
		}
		pass, err := abf.CheckLogin(inputClientLogins[0])
		require.NoError(t, err)
		require.Equal(t, false, pass)

		// Checking password buckets
		for i := 0; i <= abf.LimitPassword; i++ {
			pass, err := abf.CheckPassword(inputClientPassword[0])
			require.NoError(t, err)
			require.Equal(t, true, pass)
		}
		pass, err = abf.CheckPassword(inputClientPassword[0])
		require.NoError(t, err)
		require.Equal(t, false, pass)

		// Check IP buckets
		for i := 0; i <= abf.LimitIP; i++ {
			pass, err = abf.CheckIP(inputClientIPs[0])
			require.NoError(t, err)
			require.Equal(t, true, pass)
			pass, err = abf.CheckIP(inputClientIPs[1])
			require.NoError(t, err)
			require.Equal(t, true, pass)
		}
		pass, err = abf.CheckIP(inputClientIPs[0])
		require.NoError(t, err)
		require.Equal(t, false, pass)

		// Check login and ip buckets clearing
		err = abf.ClearLoginBuckets(inputClientLogins[0])
		require.NoError(t, err)
		if _, ok := abf.ClientsLogins[inputClientLogins[0]]; ok {
			require.Fail(t, "ClearLoginBuckets is failed")
		}
		err = abf.ClearIPBuckets(inputClientIPs[0])
		require.NoError(t, err)
		if _, ok := abf.ClientsIPs[inputClientIPs[0]]; ok {
			require.Fail(t, "ClearIPBuckets is failed")
		}

		// Check common request
		// blocked yet
		pass, err = abf.CheckRequest(inputClientIPs[1], inputClientLogins[1], inputClientPassword[1])
		require.NoError(t, err)
		require.Equal(t, false, pass)
		// new request
		for i := 0; i <= abf.LimitIP; i++ {
			pass, err := abf.CheckRequest(inputClientIPs[2], inputClientLogins[2], inputClientPassword[2])
			if i <= abf.LimitLogin {
				require.NoError(t, err)
				require.Equal(t, true, pass)
			} else {
				require.NoError(t, err)
				require.Equal(t, false, pass)
			}
		}
		// Check clearing all buckets
		abf.ClearAllBuckets()
		require.Len(t, abf.ClientsIPs, 0)
		require.Len(t, abf.ClientsLogins, 0)
		require.Len(t, abf.ClientsPasswords, 0)
	})
}
