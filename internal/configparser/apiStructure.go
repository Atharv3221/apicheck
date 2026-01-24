package configparser

type ApiConfig struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Method string `yaml:"method"`
}

type Config struct {
	Apis []ApiConfig `yaml:"apis"`
}
