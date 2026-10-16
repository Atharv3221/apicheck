package configparser

type ApiConfig struct {
	Name        string            `yaml:"name"`
	Url         string            `yaml:"url"`
	Method      string            `yaml:"method"`
	Header      map[string]string `yaml:"header,omitempty"`
	RequestBody string            `yaml:"requestBody,omitempty"`
}

type Config struct {
	Apis []ApiConfig `yaml:"apis"`
}
