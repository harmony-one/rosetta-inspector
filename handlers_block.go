package main

import (
	"context"
	"fmt"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

func renderBlock(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)
	blockID := &types.PartialBlockIdentifier{}

	// Parse out the block identifier
	id := c.Param("id")
	if len(id) > 16 {
		blockID.Hash = &id
	} else {
		var index int64
		fmt.Sscanf(id, "%d", &index)
		blockID.Index = &index
	}

	blockResp, rosettaErr, err := client.BlockAPI.Block(
		context.Background(),
		&types.BlockRequest{
			NetworkIdentifier: netID,
			BlockIdentifier:   blockID,
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "block.html", gin.H{
		"network":           netID,
		"block":             blockResp.Block,
		"otherTransactions": blockResp.OtherTransactions,
	})
}

func renderBlockTransaction(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)

	blockID := &types.BlockIdentifier{}
	txID := &types.TransactionIdentifier{
		Hash: c.Param("hash"),
	}

	// Parse out the block identifier
	id := c.Param("id")
	if len(id) > 16 {
		blockID.Hash = id
	} else {
		var index int64
		fmt.Sscanf(id, "%d", &index)
		blockID.Index = index
	}

	txResp, rosettaErr, err := client.BlockAPI.BlockTransaction(
		context.Background(),
		&types.BlockTransactionRequest{
			NetworkIdentifier:     netID,
			BlockIdentifier:       blockID,
			TransactionIdentifier: txID,
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "transaction.html", gin.H{
		"network": netID,
		"block":   blockID,
		"tx":      txResp.Transaction,
	})
}
