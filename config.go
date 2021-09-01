package config

// Load will load configurations from files or load from env args
func Load(config interface{}, files ...string) error {
	configor := NewConfigor(Files(files))
	return configor.Load(config)
}
