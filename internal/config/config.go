package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	HTTP     HTTPConfig     `env-prefix:"HTTP_"`
	Postgres PostgresConfig `env-prefix:"POSTGRES_"`
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

func Load() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	return cfg, err
}

// MustLoad panics if some error occurred
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
