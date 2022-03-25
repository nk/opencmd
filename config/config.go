package config

import "github.com/spf13/cobra"

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

var Version = "0.3.1"

var CmdArgs *cobra.Command

func GetCobraCmd() *cobra.Command {
	return CmdArgs
}
