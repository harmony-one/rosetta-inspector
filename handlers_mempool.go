package main

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

func renderMempool(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)

	mempool, rosettaErr, err := client.MempoolAPI.Mempool(
		context.Background(),
		&types.NetworkRequest{
			NetworkIdentifier: getNetworkID(c),
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "mempool.html", gin.H{
		"network": netID,
		"mempool": mempool.TransactionIdentifiers,
	})
}
