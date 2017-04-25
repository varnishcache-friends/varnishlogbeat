// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
)

type Config struct {
	Path    string        `config:"path"`
	Timeout time.Duration `config:"timeout"`
}

var DefaultConfig = Config{}
