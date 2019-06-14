# Go Config

Golang Configuration tool, used to parse toml files and load env args.

# Usage

```go
package main

import (
    "fmt"

	"github.com/x-punch/go-config"
)

func main() {
	var cfg Config
	if err := config.Load(&cfg, "config.toml"); err != nil {
		panic(err)
    }
    fmt.Println(cfg)
}
```

# Options
```go
package main

import (
    "fmt"

	"github.com/x-punch/go-config"
)

func main() {
	configor := config.NewConfigor(config.Files([]string{"example/config.toml"}), config.ShowLog(false), config.Prefix("GO"))
	var cfg Config
	if err := configor.Load(&cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
```

# Environment
```go
package main

import (
    "fmt"
    "os"

	"github.com/x-punch/go-config"
)

func main() {
    os.Setenv("CONFIG_FILES", "example/config.toml")
	os.Setenv("CONFIG_SHOW_LOG", "true")
	os.Setenv("CONFIG_PREFIX", "GO")
	os.Setenv("GO_TIMEOUT", "2h3m4s")
	var cfg Config
	if err := config.Load(&cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
```