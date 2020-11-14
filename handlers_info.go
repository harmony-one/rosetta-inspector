package main

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

func renderIndex(c *gin.Context) {
	resp, rosettaErr, err := getClient(c).NetworkAPI.NetworkList(
		context.Background(),
		&types.MetadataRequest{},
	)

	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "index.html", gin.H{
		"networks": resp.NetworkIdentifiers,
	})
}

func renderNetwork(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)

	networkStatus, rosettaErr, err := client.NetworkAPI.NetworkStatus(
		context.Background(),
		&types.NetworkRequest{
			NetworkIdentifier: getNetworkID(c),
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	networkOptions, rosettaErr, err := client.NetworkAPI.NetworkOptions(
		context.Background(),
		&types.NetworkRequest{
			NetworkIdentifier: getNetworkID(c),
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "network.html", gin.H{
		"network": netID,
		"status":  networkStatus,
		"options": networkOptions,
	})
}

func renderPeers(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)

	networkStatus, rosettaErr, err := client.NetworkAPI.NetworkStatus(
		context.Background(),
		&types.NetworkRequest{
			NetworkIdentifier: getNetworkID(c),
		},
	)
	if shouldAbort(c, rosettaErr, err) {
		return
	}

	c.HTML(200, "peers.html", gin.H{
		"network": netID,
		"peers":   networkStatus.Peers,
	})
}
