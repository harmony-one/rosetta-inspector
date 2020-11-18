package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"

	"github.com/figment-networks/rosetta-inspector/static"
)

const version = "0.2.0"

var opts struct {
	serverURL  string
	agent      string
	listenAddr string
	version    bool
}

func init() {
	gin.SetMode(gin.ReleaseMode)

	flag.StringVar(&opts.serverURL, "url", "", "Rosetta server URL")
	flag.StringVar(&opts.listenAddr, "listen", "0.0.0.0:5555", "Listen address")
	flag.BoolVar(&opts.version, "v", false, "Show version")
	flag.Parse()

	// make sure to remove trailing slashes
	opts.serverURL = strings.TrimSuffix(opts.serverURL, "/")
}

func main() {
	if opts.version {
		fmt.Println(version)
		return
	}

	if opts.serverURL == "" {
		log.Fatal("please provide rosetta server url")
	}

	tpl, err := loadTemplate(static.Assets)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.SetHTMLTemplate(tpl)
	router.Use(setClient(newClient(opts.serverURL)))

	router.GET("/", renderIndex)
	router.GET("/:blockchain/:network", renderNetwork)
	router.GET("/:blockchain/:network/peers", renderPeers)
	router.GET("/:blockchain/:network/mempool", renderMempool)
	router.GET("/:blockchain/:network/block/:id", renderBlock)
	router.GET("/:blockchain/:network/block/:id/tx/:hash", renderBlockTransaction)
	router.GET("/:blockchain/:network/account/:address", renderAccountBalance)

	go func() {
		time.Sleep(time.Millisecond * 250)
		openURL("http://" + opts.listenAddr)
	}()

	log.Println("starting server on", opts.listenAddr)
	if err := router.Run(opts.listenAddr); err != nil {
		log.Fatal(err)
	}
}

func loadTemplate(fs *assets.FileSystem) (*template.Template, error) {
	funcmap := template.FuncMap{
		"time": func(input interface{}) string {
			switch val := input.(type) {
			case time.Time:
				return val.Format(time.RFC822)
			case int64:
				return time.Unix(0, val*1000000).Format(time.RFC822)
			default:
				return fmt.Sprintf("%v", val)
			}
		},
		"amount": func(input *types.Amount) string {
			return formatAmount(input)
		},
		"json": func(input interface{}) string {
			data, err := json.MarshalIndent(input, "", "  ")
			if err != nil {
				return "invalid json:" + err.Error()
			}
			return string(data)
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

func newClient(endpoint string) *client.APIClient {
	return client.NewAPIClient(client.NewConfiguration(
		endpoint,
		"rosetta-inspector",
		&http.Client{
			Timeout: time.Second * 10,
		},
	))
}

func openURL(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Run()
	}
}
