package main

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

func renderAccountBalance(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)
	blockID := getBlockID(c)

	balanceResp, rosettaErr, err := client.AccountAPI.AccountBalance(
		context.Background(),
		&types.AccountBalanceRequest{
			NetworkIdentifier: netID,
			BlockIdentifier:   blockID,
			AccountIdentifier: &types.AccountIdentifier{
				Address: c.Param("address"),
			},
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "account.html", gin.H{
		"network": netID,
		"block":   blockID,
		"balance": balanceResp,
	})
}
