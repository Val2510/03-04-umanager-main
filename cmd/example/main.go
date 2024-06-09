package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

var app string
var version string

func init() {
	flag.StringVar(&app, "a", "my app", "-a my_app_name")
	flag.StringVar(&version, "v", "v1.0.0", "-v v.1.0.0")
	flag.Parse()
}

type CConfig struct {
	Name string `env:"TELEGRAM_TOKEN"`
}

type Config struct {
	Debug int     `env:"DEBUG,default=10"`
	C     CConfig `env:",prefix=C_"`
}

func main() {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(app, version)
}
