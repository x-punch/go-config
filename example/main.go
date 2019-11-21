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
	Tags   []tag    `toml:"tags"`
}

type tag struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

func main() {
	os.Setenv("CONFIG_SHOW_LOG", "true")
	os.Setenv("CONFIG_PREFIX", "GO")
	os.Setenv("GO_NET_PORT", "1809")
	os.Setenv("GO_NET_LOG_LEVEL", "h,k")
	os.Setenv("GO_NET_LOG_LEVEL_0", "err")
	os.Setenv("GO_NET_LOG_LEVEL_1", "info")
	os.Setenv("GO_NET_LOG_TAGS_LEN", "2")
	os.Setenv("GO_NET_LOG_TAGS_0_NAME", "err")
	os.Setenv("GO_NET_LOG_TAGS_1_Version", "2.0.1")
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
	os.Setenv("CONFIG_FILES", "example/config.toml")
	os.Setenv("GO_NET_PORT", "1809")
	os.Setenv("GO_SHOPPING_DISCOUNT", "0.75")
	os.Setenv("GO_SHOPPING_DISCOUNT_TIME", "1979-05-27T07:32:00Z")
	os.Setenv("GO_TIMEOUT", "2h3m4s")
	configor := config.NewConfigor(config.Files([]string{"example.toml"}), config.ShowLog(false), config.Prefix("GO"))
	var cfg Config
	if err := configor.Load(&cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
