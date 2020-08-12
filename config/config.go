package config

type Config struct {
	DataPath string
}

func Get() Config {
	return Config{DataPath: "/tmp"}
}
