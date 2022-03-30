package version

import (
	"html/template"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		t := template.New("version")
		template.Must(t.Parse(cmd.Root().VersionTemplate()))
		err := t.Execute(cmd.Root().OutOrStdout(), cmd.Parent())
		if err != nil {
			cmd.Root().Println(err)
		}
	},
}
