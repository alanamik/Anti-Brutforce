package antibrutforce

import (
	"OTUS_hws/Anti-BruteForce/internal/config"
	"net"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestListFunctions(t *testing.T) {
	os.Chdir("..")
	os.Chdir("..")
	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "[config.New()]")
		panic(err)
	}
	abf, _ := New(conf)

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
}
