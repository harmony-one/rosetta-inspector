package main

import (
	"fmt"
	"html/template"
	"math"
	"math/big"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

func setClient(client *client.APIClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("client", client)
	}
}

func getClient(c *gin.Context) *client.APIClient {
	return c.MustGet("client").(*client.APIClient)
}

func getNetworkID(c *gin.Context) *types.NetworkIdentifier {
	return &types.NetworkIdentifier{
		Blockchain: c.Param("blockchain"),
		Network:    c.Param("network"),
	}
}

func getBlockID(c *gin.Context) *types.PartialBlockIdentifier {
	var blockID *types.PartialBlockIdentifier

	if number := c.Query("block_number"); number != "" {
		blockID = &types.PartialBlockIdentifier{}
		var index int64
		fmt.Sscanf(number, "%d", &index)
		blockID.Index = &index
		return blockID
	}

	if hash := c.Query("block_hash"); hash != "" {
		blockID = &types.PartialBlockIdentifier{}
		blockID.Hash = &hash
	}

	return blockID
}

func shouldAbort(c *gin.Context, rosettaErr *types.Error, err error) bool {
	if rosettaErr != nil || err != nil {
		renderError(c, rosettaErr, err)
		return true
	}
	return false
}

func renderError(c *gin.Context, rosettaErr *types.Error, err error) {
	c.HTML(400, "error.html", gin.H{
		"network":      getNetworkID(c),
		"rosettaError": rosettaErr,
		"error":        err,
	})
}

func formatAmount(amount *types.Amount) template.HTML {
	if amount == nil {
		return "N/A"
	}

	n := big.Float{}
	val, ok := n.SetString(amount.Value)
	if !ok {
		return "invalid value"
	}

	d := big.NewFloat(math.Pow10(int(amount.Currency.Decimals)))
	result := new(big.Float).Quo(val, d)

	content := fmt.Sprintf(`<span title="%v">%.8f %v</span>`,
		amount.Value,
		result,
		amount.Currency.Symbol,
	)

	return template.HTML(content)
}
