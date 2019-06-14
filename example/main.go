package main

import (
	"fmt"
	"os"
	"time"

	"github.com/x-punch/go-config"
)

// Config represents config
type Config struct {
	Timeout config.Duration `toml:"timeout"`
	Shop    shop            `toml:"shopping"`
	Network network         `toml:"net"`
}

type shop struct {
	Discount float64   `toml:"discount"`
	Time     time.Time `toml:"discount_time"`
}

type network struct {
	Host  string  `toml:"host"`
	Port  uint16  `toml:"port"`
	Log   log     `toml:"log"`
	Nodes []int32 `toml:"nodes"`
}

type log struct {
	Enable bool     `toml:"enable"`
	Level  []string `toml:"level"`
}

func main() {
	os.Setenv("CONFIG_SHOW_LOG", "true")
	os.Setenv("CONFIG_PREFIX", "GO")
	os.Setenv("GO_NET_PORT", "1809")
	os.Setenv("GO_SHOPPING_DISCOUNT", "0.75")
	os.Setenv("GO_SHOPPING_DISCOUNT_TIME", "1979-05-27T07:32:00Z")
	os.Setenv("GO_TIMEOUT", "2h3m4s")
	var cfg Config
	if err := config.Load(&cfg, "example/config.toml"); err != nil {
		panic(err)
	}
}

func env() {
	os.Setenv("CONFIG_FILES", "example/config.toml")
	os.Setenv("CONFIG_SHOW_LOG", "true")
	os.Setenv("CONFIG_PREFIX", "GO")
	os.Setenv("GO_NET_PORT", "1809")
	os.Setenv("GO_SHOPPING_DISCOUNT", "0.75")
	os.Setenv("GO_SHOPPING_DISCOUNT_TIME", "1979-05-27T07:32:00Z")
	os.Setenv("GO_TIMEOUT", "2h3m4s")
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		panic(err)
	}
}

func options() {
	os.Setenv("GO_NET_PORT", "1809")
	os.Setenv("GO_SHOPPING_DISCOUNT", "0.75")
	os.Setenv("GO_SHOPPING_DISCOUNT_TIME", "1979-05-27T07:32:00Z")
	os.Setenv("GO_TIMEOUT", "2h3m4s")
	configor := config.NewConfigor(config.Files([]string{"example/config.toml"}), config.ShowLog(false), config.Prefix("GO"))
	var cfg Config
	if err := configor.Load(&cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
