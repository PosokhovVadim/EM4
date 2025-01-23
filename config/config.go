package config

type Config struct {
	Storage Storage `env:"STORAGE"`
}

type Storage struct {
	Path string `env:"PATH" required:"true"`
}
