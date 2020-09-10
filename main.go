package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
)

var opts struct {
	serverURL  string
	agent      string
	listenAddr string
}

func init() {
	flag.StringVar(&opts.serverURL, "url", "", "Rosetta server URL")
	flag.StringVar(&opts.listenAddr, "listen", "0.0.0.0:5555", "Listen address")
	flag.Parse()

	if opts.serverURL == "" {
		log.Fatal("please provide rosetta server url")
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*.html")

	router.Use(func(c *gin.Context) {
		c.Set("client", initClient())
	})

	router.GET("/", renderHome)
	router.GET("/:blockchain/:network", renderNetwork)
	router.GET("/:blockchain/:network/block/:id", renderBlock)
	router.GET("/:blockchain/:network/account/:address", renderAccountBalance)

	log.Fatal(router.Run(opts.listenAddr))
}

func initClient() *client.APIClient {
	clientCfg := client.NewConfiguration(
		opts.serverURL,
		"rosetta-go-sdk",
		&http.Client{
			Timeout: time.Second * 10,
		},
	)

	return client.NewAPIClient(clientCfg)
}

func renderHome(c *gin.Context) {
	client := c.MustGet("client").(*client.APIClient)

	networkList, rosettaErr, err := client.NetworkAPI.NetworkList(
		context.Background(),
		&types.MetadataRequest{},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	if len(networkList.NetworkIdentifiers) == 0 {
		log.Fatal("no available networks")
	}

	c.HTML(200, "index.html", gin.H{
		"networks": networkList.NetworkIdentifiers,
	})
}

func renderNetwork(c *gin.Context) {
	client := c.MustGet("client").(*client.APIClient)

	identifier := &types.NetworkIdentifier{
		Blockchain: c.Param("blockchain"),
		Network:    c.Param("network"),
	}

	networkStatus, rosettaErr, err := client.NetworkAPI.NetworkStatus(
		context.Background(),
		&types.NetworkRequest{
			NetworkIdentifier: identifier,
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(200, "network.html", gin.H{
		"identifier": identifier,
		"status":     networkStatus,
	})
}

func renderBlock(c *gin.Context) {
	client := c.MustGet("client").(*client.APIClient)

	netId := &types.NetworkIdentifier{
		Blockchain: c.Param("blockchain"),
		Network:    c.Param("network"),
	}

	blockID := &types.PartialBlockIdentifier{}

	id := c.Param("id")
	if strings.Contains(id, "0x") {
		blockID.Hash = &id
	} else {
		var index int64
		fmt.Sscanf(id, "%d", &index)
		blockID.Index = &index
	}

	block, _, err := client.BlockAPI.Block(
		context.Background(),
		&types.BlockRequest{
			NetworkIdentifier: netId,
			BlockIdentifier:   blockID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(200, "block.html", gin.H{
		"network": netId,
		"block":   block.Block,
	})
}

func renderAccountBalance(c *gin.Context) {
	client := c.MustGet("client").(*client.APIClient)

	netId := &types.NetworkIdentifier{
		Blockchain: c.Param("blockchain"),
		Network:    c.Param("network"),
	}

	balance, _, err := client.AccountAPI.AccountBalance(
		context.Background(),
		&types.AccountBalanceRequest{
			NetworkIdentifier: netId,
			AccountIdentifier: &types.AccountIdentifier{
				Address: c.Param("address"),
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	data := gin.H{
		"network": netId,
		"balance": balance,
	}

	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		c.HTML(200, "account.html", data)
	case gin.MIMEJSON:
		c.JSON(200, data)
	}
}
