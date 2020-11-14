package main

import (
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
