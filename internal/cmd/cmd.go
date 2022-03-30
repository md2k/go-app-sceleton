package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"web-app-go/internal/cmd/run"
	"web-app-go/internal/cmd/version"
	"web-app-go/internal/config"
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "Location of configuration file.")
	// Adding Commands
	rootCmd.AddCommand(version.Command)
	rootCmd.AddCommand(run.Command)
}

var rootCmd = &cobra.Command{
	Use:   "web-app-go",
	Short: "Simple Web App",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		//utils.Dump(config.GetConfig())
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
}

func Execute(version string) {
	rootCmd.SetVersionTemplate("{{with .Name}}{{printf \"%s: \" .}}{{end}}{{printf \"%s\" .Version}}\n")
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	err := config.Loader(rootCmd.Root().PersistentFlags().Lookup("config").Value.String())
	if err != nil {
		log.Printf("[WARNING  ] Not able to load config: %s", err.Error())
		os.Exit(1)
	}
}
