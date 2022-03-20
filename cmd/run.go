/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"opencmd/command"
	"opencmd/utils"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a command",
	Long:  `run a command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("args:", args[0])
		// fixme: set working directory
		// fixme: 同.opencmd目录下有两个同名命令的处理

		commandPath, err := utils.FindCommandByName(
			utils.GetCurrentDir(), args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		cmdInfo := &command.Command{
			Name: args[0],
			Path: commandPath,
		}
		output, runErr := cmdInfo.Run()
		if runErr != nil {
			fmt.Println(runErr)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}