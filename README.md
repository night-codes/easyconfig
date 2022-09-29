# Easyconfig [![Go Reference](https://pkg.go.dev/badge/github.com/night-codes/easyconfig.svg)](https://pkg.go.dev/github.com/night-codes/easyconfig)

Easy-to-use multi source configuration library in Go.

## Documentation

[Docs on pkg.go.dev](https://pkg.go.dev/github.com/night-codes/easyconfig)

## How To Install

```bash
go get github.com/night-codes/easyconfig
```

## Example

The usage looks like this:

```go
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
		Slice            []string `json:"-" yaml:"-" env:"slice,:"` // env list separator is ":"
		Version          string   `json:"-" yaml:"-" dir:"-" env:"-"` // ignore field
	}
)

func main() {
	config := &Config{
		PostgresHost: "localhost", // default value
	}

	// Simple configuration load
	jsonLoader := easyconfig.JSONSource{Path: "./config.json"}
	if err := jsonLoader.Load(config); err != nil {
		fmt.Println(err)
		// use another loader, for example
		// ....
	}

	// Multi sources loader
	loader := easyconfig.NewLoader(
		[]easyconfig.Source{
			// all sources are sorted in the order in which data loading from each tip will be attempted
			easyconfig.YAMLSource{Path: "./config.yaml"},
			easyconfig.JSONSource{Path: "./config.json"},
			easyconfig.TOMLSource{Path: "./config.toml"},
			easyconfig.EnvFileSource{Prefix: "APP", Path: "./config.env"},
			easyconfig.DirSource{Path: "./k8s-secret"}, // loads configuration from the mounted Secrets or ConfigMap kubernetes directory (file=value)
			easyconfig.EnvSource{Prefix: "APP"},        // loads configuration from environment variables
			easyconfig.FlagsSource{},                   // loads configuration from flags
		},
	)

	loader.Load(config) // collect data from each source
	fmt.Printf("%v\n", config)
}
```

You can use -help flag in your application:

```bash
❯ go run . -help

Usage:
    app [arguments]

The commands are:

    -postgresUser
        Set value of PostgresUser
    -postgresPassword
        Set value of PostgresPassword
    -postgresHost
        Set value of PostgresHost. Default: "localhost"
    -postgresPort
        Set number value of PostgresPort
    -postgresDBName
        Set value of PostgresDBName
    -postgresSSLMode
        Set value of PostgresSSLMode
    -slice
        Set value of Slice
    -version
        Set value of Version

Environment variables to use:

    APP_POSTGRES_USER
    APP_POSTGRES_PASSWORD
    APP_POSTGRES_HOST
    APP_POSTGRES_PORT
    APP_POSTGRES_DB_NAME
    APP_POSTGRES_SSL_MODE
    APP_SLICE

Configuration directory files to use:

    postgres-user
    postgres-password
    postgres-host
    postgres-port
    postgres-db-name
    postgres-ssl-mode
    slice
```

## License

MIT License

Copyright (C) 2022 Oleksiy Chechel (alex.mirrr@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
