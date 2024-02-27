/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wiggers/goexpert/desafio-tecnico/2-stress-test/internal"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		internal.Execute(url, requests, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringP("url", "u", "", "Informe a url a ser consultada")
	executeCmd.MarkFlagRequired("url")

	executeCmd.Flags().Int("requests", 100, "Informe a quantidade de requisições que serão feitas")
	executeCmd.Flags().Int("concurrency", 10, "Informe a qtd de requisições simultaneas")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
