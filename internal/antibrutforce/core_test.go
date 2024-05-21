package antibrutforce

import (
	"context"
	"os"
	"testing"

	"OTUS_hws/Anti-BruteForce/internal/config"
	"OTUS_hws/Anti-BruteForce/internal/redisdb"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	os.Chdir("..")
	os.Chdir("..")
	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "[config.New()]")
		panic(err)
	}

	ctx := context.Background()
	redisClient := redisdb.NewClient(*conf)
	//	redisClient.Client.FlushDB(ctx)

	// check adding to lists
	inputBlockedIPs := []string{"126.10.10.11/10", "120.1.5.7/10", "128.10.10.11/10"}
	for _, ip := range inputBlockedIPs {
		err = redisClient.AddToList(ctx, ip, redisdb.Blacklist)
		require.NoError(t, err)
	}

	inputPassedIPs := []string{"56.10.10.11/10", "45.12.50.7/10", "125.1.14.10/10"}
	for _, ip := range inputPassedIPs {
		err = redisClient.AddToList(ctx, ip, redisdb.Whitelist)
		require.NoError(t, err)
	}

	blockedIPs, err := redisClient.GetAllIPFromList(ctx, redisdb.Blacklist)
	require.NoError(t, err)
	for i, j := 0, len(blockedIPs)-1; i < j; i, j = i+1, j-1 {
		blockedIPs[i], blockedIPs[j] = blockedIPs[j], blockedIPs[i]
	}
	require.Equal(t, inputBlockedIPs, blockedIPs)

	passedIPs, err := redisClient.GetAllIPFromList(ctx, redisdb.Whitelist)
	require.NoError(t, err)
	for i, j := 0, len(passedIPs)-1; i < j; i, j = i+1, j-1 {
		passedIPs[i], passedIPs[j] = passedIPs[j], passedIPs[i]
	}

	require.EqualValues(t, inputPassedIPs, passedIPs)

	// check adding IP, that is in a list
	err = redisClient.AddToList(ctx, inputBlockedIPs[2], redisdb.Blacklist)
	require.ErrorIs(t, redisdb.ErrIPInListYet, err)

	err = redisClient.AddToList(ctx, inputPassedIPs[2], redisdb.Whitelist)
	require.ErrorIs(t, redisdb.ErrIPInListYet, err)

	// check adding IP in a non-existent list
	err = redisClient.AddToList(ctx, "testIP", "nonExistentList")
	require.ErrorIs(t, redisdb.ErrListDoesNotExist, err)

	// check deleting from lists
	err = redisClient.DeleteFromList(ctx, inputBlockedIPs[2], redisdb.Blacklist)
	require.NoError(t, err)
	isInList, err := redisClient.CheckInList(ctx, inputBlockedIPs[2], redisdb.Blacklist)
	require.NoError(t, err)
	require.Equal(t, false, isInList)

	err = redisClient.DeleteFromList(ctx, inputPassedIPs[2], redisdb.Whitelist)
	require.NoError(t, err)
	isInList, err = redisClient.CheckInList(ctx, inputPassedIPs[2], redisdb.Whitelist)
	require.NoError(t, err)
	require.Equal(t, false, isInList)
}
