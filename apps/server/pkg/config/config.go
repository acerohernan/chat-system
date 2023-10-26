package config

type Config struct {
	Redis  *RedisConfig
	Logger *LoggerConfig
}

// TODO: Add redis connection params
type RedisConfig struct {
	Host string
}

type LoggerConfig struct {
	Level string
}

// TODO: Use env variables
func NewConfig() *Config {
	return &Config{
		Redis: &RedisConfig{
			Host: "localhost:6379",
		},
		Logger: &LoggerConfig{
			Level: "debug",
		},
	}
}
