package errata

type WebUIConfig struct {
	Source   string
	BindAddr string
}

type CodeGenConfig struct {
	Source   string
	Template string
	Package  string
}
