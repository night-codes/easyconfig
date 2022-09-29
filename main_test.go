package easyconfig

import (
	"fmt"
	"os"
	"testing"
)

type (
	Config struct {
		PostgresUser     string   `json:"postgresUser" yaml:"postgresUser" dir:"postgres-user"`
		PostgresPassword string   `json:"postgresPassword" yaml:"postgresPassword"`
		PostgresHost     string   `json:"postgresHost" yaml:"postgresHost"`
		PostgresPort     uint64   `json:"postgresPort" yaml:"postgresPort"`
		PostgresDBName   string   `json:"postgresDBName" yaml:"postgresDBName"`
		PostgresSSLMode  string   `json:"postgresSSLMode" yaml:"postgresSSLMode"`
		Slice            []string `json:"slice" yaml:"slice" env:"slice,:"`
	}
)

func TestEnvFileSourceLoader(t *testing.T) {
	t.Run("EnvFileSource.Load", func(t *testing.T) {
		config := new(Config)
		loader := NewLoader([]Source{
			EnvFileSource{"APP", "tests/config.env"},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestYAMLSourceLoader(t *testing.T) {
	t.Run("YAMLSource.Load", func(t *testing.T) {
		config := new(Config)
		loader := NewLoader([]Source{
			YAMLSource{"tests/config.yaml"},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestJSONSourceLoader(t *testing.T) {
	t.Run("JSONSource.Load", func(t *testing.T) {
		config := new(Config)
		loader := NewLoader([]Source{
			JSONSource{"tests/config.json"},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestTOMLSourceLoader(t *testing.T) {
	t.Run("TOMLSource.Load", func(t *testing.T) {
		config := new(Config)
		loader := NewLoader([]Source{
			TOMLSource{"tests/config.toml"},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestDirSourceLoader(t *testing.T) {
	t.Run("DirSource.Load", func(t *testing.T) {
		config := new(Config)
		loader := NewLoader([]Source{
			DirSource{"tests/config"},
			// EnvSource{Prefix: "APP"},
			// FlagsSource{},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestEnvSourceLoader(t *testing.T) {
	t.Run("EnvSource.Load", func(t *testing.T) {
		env := map[string]string{
			"APP_POSTGRES_USER":     "postgres",
			"APP_POSTGRES_PASSWORD": "password",
			"APP_POSTGRES_HOST":     "localhost",
			"APP_POSTGRES_PORT":     "5432",
			"APP_POSTGRES_DB_NAME":  "db-name",
			"APP_POSTGRES_SSL_MODE": "disable",
			"APP_SLICE":             "a1:a2:a3:a4:a5",
		}
		for key, val := range env {
			if err := os.Setenv(key, val); err != nil {
				t.Fatal(err)
			}
		}
		config := new(Config)
		loader := NewLoader([]Source{
			EnvSource{Prefix: "APP"},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}

func TestFlagsSourceLoader(t *testing.T) {
	t.Run("FlagsSource.Load", func(t *testing.T) {

		os.Args = []string{
			"easyconfig",
			"-postgresUser=postgres",
			"-postgresPassword=password",
			"-postgresHost=localhost",
			"-postgresPort=5432",
			"-postgresDBName=db-name",
			"-postgresSSLMode=disable",
			"-slice=a1,a2,a3,a4,a5",
		}
		config := new(Config)
		loader := NewLoader([]Source{
			FlagsSource{},
		})
		if err := loader.Load(config); err != nil {
			t.Errorf("Error = %s, want %s", err.Error(), "nil")
		}
		if config.PostgresUser != "postgres" {
			t.Errorf("PostgresUser = %s, want %s", config.PostgresUser, "postgres")
		}
		if config.PostgresPassword != "password" {
			t.Errorf("PostgresPassword = %s, want %s", config.PostgresPassword, "password")
		}
		if config.PostgresHost != "localhost" {
			t.Errorf("PostgresHost = %s, want %s", config.PostgresHost, "localhost")
		}
		if config.PostgresDBName != "db-name" {
			t.Errorf("PostgresDBName = %s, want %s", config.PostgresDBName, "db-name")
		}
		if config.PostgresPort != 5432 {
			t.Errorf("PostgresPort = %d, want %d", config.PostgresPort, 5432)
		}
		if config.PostgresSSLMode != "disable" {
			t.Errorf("PostgresSSLMode = %s, want %s", config.PostgresSSLMode, "disable")
		}
		if fmt.Sprintf("%v", config.Slice) != "[a1 a2 a3 a4 a5]" {
			t.Errorf("PostgresSSLMode = %v, want %s", config.Slice, "[a1 a2 a3 a4 a5]")
		}
	})
}
