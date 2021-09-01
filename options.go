package config

// Options represents options of configor
type Options struct {
	Files       []string
	Prefix      string
	ShowLog     bool
	LoadFromEnv bool
}

// Option represents func to set options
type Option func(*Options)

// Files represents configuration files path
func Files(files []string) Option {
	return func(o *Options) {
		o.Files = append(o.Files, files...)
	}
}

// Prefix represents env prefix load from env args
func Prefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// ShowLog represents wether show logs of configuration
func ShowLog(show bool) Option {
	return func(o *Options) {
		o.ShowLog = show
	}
}

// LoadFromEnv represents wether load configuration from env args
func LoadFromEnv(load bool) Option {
	return func(o *Options) {
		o.LoadFromEnv = load
	}
}

func newOptions(options ...Option) *Options {
	opts := &Options{LoadFromEnv: true}
	for _, option := range options {
		option(opts)
	}
	return opts
}
