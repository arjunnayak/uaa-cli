package cmd

import (
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
)

var userinfoCmd = cobra.Command{
	Use:     "userinfo",
	Short:   "See claims about the authenticated user",
	Aliases: []string{"me"},
	Long:    help.Userinfo(),
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		i, err := uaa.Me(GetHttpClient(), GetSavedConfig())
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		j, err := json.MarshalIndent(&i, "", "  ")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Robots(string(j))
	},
}

func init() {
	RootCmd.AddCommand(&userinfoCmd)
	userinfoCmd.Annotations = make(map[string]string)
	userinfoCmd.Annotations[MISC_CATEGORY] = "true"
}