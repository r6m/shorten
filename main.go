package main

import (
	"flag"
	"strings"

	"github.com/r6m/shorten/handlers"
	"github.com/r6m/shorten/store"
	"github.com/r6m/shorten/store/memory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	port = flag.String("port", "8080", "port to listen")
)

func init() {
	flag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	viper.SetDefault("port", "8080")
}

func main() {
	var store store.Store
	driver := viper.GetString("db.driver")
	switch driver {
	case "memory", "":
		store = memory.NewStore()
	case "mysql":
		// implement other drivers
		// store = mysql.NewStorage(viper.GetString("db.url"))
		fallthrough
	default:
		logrus.Fatalf("storage driver '%s' is not implemented", driver)
	}

	server := handlers.NewServer(store)
	addr := ":" + viper.GetString("port")
	server.ListenAndServe(addr)
}
