package config

type Config struct {
	OpencmdDir   string
	CommandDir   string
	DefaultShell string
}

var DefaultConfig = Config{
	OpencmdDir:   ".opencmd",
	CommandDir:   "commands",
	DefaultShell: "/bin/bash",
}

var Version = "0.2.0"
