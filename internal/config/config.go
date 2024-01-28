package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/shopspring/decimal"
)

type Config struct {
	HTTP     HTTPConfig     `env-prefix:"HTTP_"`
	Postgres PostgresConfig `env-prefix:"POSTGRES_"`
	Wallet   WalletConfig   `env-prefix:"WALLET_" env:"WALLET_INITIAL_BALANCE"`
}

type HTTPConfig struct {
	Host string `env:"HOST" env-default:""`
	Port string `env:"PORT" env-default:"8080"`
}

func (cfg HTTPConfig) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

type PostgresConfig struct {
	ConnectionString string `env:"CONNECTION_STRING" env-required:"true"`
}

type WalletConfig struct {
	InitialBalanceString string `env:"INITIAL_BALANCE" env-default:"100"`
	InitialBalance       decimal.Decimal
}

func Load() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	initialBalance, err := decimal.NewFromString(cfg.Wallet.InitialBalanceString)
	if err != nil {
		return cfg, err
	}
	cfg.Wallet.InitialBalance = initialBalance

	return cfg, nil
}

// MustLoad panics if some error occurred
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
