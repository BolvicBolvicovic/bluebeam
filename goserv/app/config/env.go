package config

import (
	"log"
	"github.com/spf13/viper"
)

type Env struct {
	// server
	GoMode     string `mapstructure:"GO_MODE"`
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort uint16 `mapstructure:"SERVER_PORT"`
	// database
	DBHost         string `mapstructure:"DB_HOST"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         uint16 `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBUserPwd      string `mapstructure:"DB_USER_PWD"`
	DBMinPoolSize  uint16 `mapstructure:"DB_MIN_POOL_SIZE"`
	DBMaxPoolSize  uint16 `mapstructure:"DB_MAX_POOL_SIZE"`
	DBQueryTimeout uint16 `mapstructure:"DB_QUERY_TIMEOUT_SEC"`
	// keys
	RSAPrivateKeyPath string `mapstructure:"RSA_PRIVATE_KEY_PATH"`
	RSAPublicKeyPath  string `mapstructure:"RSA_PUBLIC_KEY_PATH"`
}

func NewEnv(filename string, override bool) *Env {
	env := Env{}

	viper.SetConfigFile(filename)

	if override { viper.AutomaticEnv() }
	
	err := viper.ReadInConfig()
	if err != nil { log.Fatal("Error reading env file", err) }
	
	err = viper.Unmarshal(&env)
	if err != nil { log.Fatal("Error loading env file", err) }

	return &env
}
