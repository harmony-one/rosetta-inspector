package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"

	"github.com/figment-networks/rosetta-inspector/static"
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

	// make sure to remove trailing slashes
	opts.serverURL = strings.TrimSuffix(opts.serverURL, "/")
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	tpl, err := loadTemplate(static.Assets)
	if err != nil {
		log.Fatal(err)
	}
	router.SetHTMLTemplate(tpl)

	router.Use(func(c *gin.Context) {
		c.Set("client", initClient())
	})

	router.GET("/", renderHome)
	router.GET("/:blockchain/:network", renderNetwork)
	router.GET("/:blockchain/:network/block/:id", renderBlock)
	router.GET("/:blockchain/:network/block/:id/tx/:hash", renderBlockTransaction)
	router.GET("/:blockchain/:network/account/:address", renderAccountBalance)

	done := make(chan error)

	go func() {
		log.Println("starting server on", opts.listenAddr)
		done <- router.Run(opts.listenAddr)
	}()

	go func() {
		if err := exec.Command("open", "http://"+opts.listenAddr).Run(); err != nil {
			log.Println("cant open server url:", err)
		}
	}()

	<-done
}

func loadTemplate(fs *assets.FileSystem) (*template.Template, error) {
	funcmap := template.FuncMap{
		"time": func(input interface{}) string {
			switch input.(type) {
			case time.Time:
				return input.(time.Time).Format(time.RFC822)
			default:
				return "invalid-value"
			}
		},
	}

	t := template.New("").Funcs(funcmap)

	for name, file := range fs.Files {
		shortname := filepath.Base(name)
		if file.IsDir() || !strings.HasSuffix(shortname, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(shortname).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func initClient() *client.APIClient {
	clientCfg := client.NewConfiguration(
		opts.serverURL,
		"rosetta-inspector",
		&http.Client{
			Timeout: time.Second * 10,
		},
	)

	return client.NewAPIClient(clientCfg)
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
		"rosettaError": rosettaErr,
		"error":        err,
	})
}

func renderHome(c *gin.Context) {
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

func renderAccountBalance(c *gin.Context) {
	client := getClient(c)
	netID := getNetworkID(c)

	balanceResp, rosettaErr, err := client.AccountAPI.AccountBalance(
		context.Background(),
		&types.AccountBalanceRequest{
			NetworkIdentifier: netID,
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
		"balance": balanceResp,
	})
}
