package configparser

type ApiConfig struct {
	Name        string            `yaml:"name"`
	Url         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Header      map[string]string `yaml:"header"`
	RequestBody string            `yaml:"requestBody"`
}

type Config struct {
	Apis []ApiConfig `yaml:"apis"`
}
