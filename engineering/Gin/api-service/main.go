package main

import (
	"api-service/pkg/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"api-service/config"
	"api-service/model"
	v "api-service/pkg/version"
	"api-service/router"
	"api-service/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "", "api-service config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title Example API
// @version 1.0
// @description api-service demo

// @contact.name yipwinghong
// @contact.url http://www.swagger.io/support
// @contact.email yipwinghong

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()
	// API version info
	if *version {
		value := v.Get()
		marshalled, err := json.MarshalIndent(&value, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// Set gin mode.
	gin.SetMode(viper.GetString("mode"))

	// Create the Gin engine.
	g := gin.New()

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlewares.
		middleware.Logging(),
		middleware.RequestId(),
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Error("The router has no response, or it might took too long to start up.", err)
		}
		log.Infof("The router has been deployed successfully.")
	}()

	// Start to listening the incoming requests.
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Infof(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}
	log.Infof("Start to listening the incoming requests on http address: %s\n", viper.GetString("addr"))
	log.Infof(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {

		// Ping the server by sending a GET request to `/health`.

		resp, err := http.Get(fmt.Sprintf("%s/%s/monitor/health", viper.GetString("url"), viper.GetString("version")))
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Infof("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
