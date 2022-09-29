package main

import (
	"fmt"

	"github.com/night-codes/easyconfig"
)

type (
	Config struct {
		PostgresUser     string   `json:"postgresUser" yaml:"postgresUser" dir:"postgres-user" env:"APP_POSTGRES_USER"`
		PostgresPassword string   `json:"postgresPassword" yaml:"postgresPassword"`
		PostgresHost     string   `json:"postgresHost" yaml:"postgresHost"`
		PostgresPort     uint64   `json:"postgresPort" yaml:"postgresPort"`
		PostgresDBName   string   `json:"postgresDBName" yaml:"postgresDBName"`
		PostgresSSLMode  string   `json:"postgresSSLMode" yaml:"postgresSSLMode"`
		Slice            []string `json:"-" yaml:"-" toml:"-" env:"slice,:"`   // env list separator is ":"
		Version          string   `json:"-" yaml:"-" dir:"-" toml:"-" env:"-"` // ignore field for all sources
	}
)

func main() {
	config := &Config{
		PostgresHost: "localhost", // default value
	}

	// Simple configuration load
	jsonLoader := easyconfig.JSONSource{Path: "../../tests/config.json"}
	if err := jsonLoader.Load(config); err != nil {
		fmt.Println(err)
		// use another loader, for example
		// ....
	}

	// Multi sources loader
	loader := easyconfig.NewLoader(
		[]easyconfig.Source{
			// all sources are sorted in the order in which data loading from each tip will be attempted
			easyconfig.YAMLSource{Path: "../../tests/config.yaml"},
			easyconfig.JSONSource{Path: "../../tests/config.json"},
			easyconfig.TOMLSource{Path: "../../tests/config.toml"},
			easyconfig.EnvFileSource{Prefix: "APP", Path: "../../tests/config.env"},
			easyconfig.DirSource{Path: "../../tests/k8s-secret"}, // loads configuration from the mounted Secrets or ConfigMap kubernetes directory (file=value)
			easyconfig.EnvSource{Prefix: "APP"},                  // loads configuration from environment variables
			easyconfig.FlagsSource{},                             // loads configuration from flags
		},
	)

	loader.Load(config) // collect data from each source
	fmt.Printf("%v\n", config)
}
