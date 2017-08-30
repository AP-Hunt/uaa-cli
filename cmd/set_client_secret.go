package cmd

import (
	"github.com/spf13/cobra"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"fmt"
	"os"
)

var setClientSecretCmd = &cobra.Command{
	Use:   "set-client-secret CLIENT_ID -s CLIENT_SECRET",
	Short: "Update secret for a client",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureContext()
	},
	Run: func(cmd *cobra.Command, args []string) {
		clientId := args[0]
		c := GetSavedConfig()
		cm := &uaa.ClientManager{GetHttpClient(), c}
		err := cm.ChangeSecret(clientId, clientSecret)
		if err != nil {
			fmt.Printf("The secret for client %v was not updated.\n", clientId)
			TraceRetryMsg(c)
			os.Exit(1)
		}
		fmt.Printf("The secret for client %v has been successfully updated.\n", clientId)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return MissingArgument("client_id")
		}
		if clientSecret == "" {
			return MissingArgument("client_secret")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(setClientSecretCmd)
	setClientSecretCmd.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "new client secret")
	setClientSecretCmd.Flags().StringVarP(&zoneSubdomain, "zone", "z", "", "the identity zone subdomain where the client resides")
}
