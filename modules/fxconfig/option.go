package fxconfig

type options struct {
	FileName  string
	FilePaths []string
}

var defaultConfigOptions = options{
	FileName: "config",
	FilePaths: []string{
		".",
		"./configs",
	},
}

type ConfigOption func(o *options)

func WithFileName(n string) ConfigOption {
	return func(o *options) {
		o.FileName = n
	}
}

func WithFilePaths(p ...string) ConfigOption {
	return func(o *options) {
		o.FilePaths = p
	}
}
