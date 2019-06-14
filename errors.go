package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// UnmatchedKeysError errors are returned by the Load function when
// ErrorOnUnmatchedKeys is set to true and there are unmatched keys in the input
// toml config file. The string returned by Error() contains the names of the
// missing keys.
type UnmatchedKeysError struct {
	Keys []toml.Key
}

func (e *UnmatchedKeysError) Error() string {
	return fmt.Sprintf("There are keys in the config file that do not match any field in the given struct: %v", e.Keys)
}
