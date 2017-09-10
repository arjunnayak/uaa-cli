package cmd

import (
	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"code.cloudfoundry.org/uaa-cli/utils"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func addAuthcodeTokenToContext(clientId string, tokenResponse uaa.TokenResponse, log *utils.Logger) {
	ctx := uaa.UaaContext{
		GrantType:     uaa.AUTHCODE,
		ClientId:      clientId,
		TokenResponse: tokenResponse,
	}

	SaveContext(ctx, log)
}

func AuthcodeTokenArgumentValidation(args []string, port int, cmd *cobra.Command) error {
	if len(args) < 1 {
		MissingArgument("client_id", cmd)
	}
	if port == 0 {
		MissingArgumentWithExplanation("port", cmd, `The port number must correspond to a localhost redirect_uri specified in the client configuration.`)
	}
	if clientSecret == "" {
		MissingArgument("client_secret", cmd)
	}
	validateTokenFormat(cmd, tokenFormat)
	return nil
}

func AuthcodeTokenCommandRun(doneRunning chan bool, clientId string, authcodeImp cli.ClientImpersonator, log *utils.Logger) {
	authcodeImp.Start()
	authcodeImp.Authorize()
	tokenResponse := <-authcodeImp.Done()
	addAuthcodeTokenToContext(clientId, tokenResponse, log)
	doneRunning <- true
}

var getAuthcodeToken = &cobra.Command{
	Use:   "get-authcode-token CLIENT_ID -s CLIENT_SECRET --port REDIRECT_URI_PORT",
	Short: "Obtain a token using the authorization_code grant type",
	PreRun: func(cmd *cobra.Command, args []string) {
		EnsureTarget()
	},
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool)
		authcodeImp := cli.NewAuthcodeClientImpersonator(GetHttpClient(), GetSavedConfig(), args[0], clientSecret, tokenFormat, scope, port, log, open.Run)
		go AuthcodeTokenCommandRun(done, args[0], authcodeImp, GetLogger())
		<-done
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return AuthcodeTokenArgumentValidation(args, port, cmd)
	},
}

func init() {
	getAuthcodeToken.Flags().IntVarP(&port, "port", "", 0, "port on which to run local callback server")
	getAuthcodeToken.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getAuthcodeToken.Flags().StringVarP(&scope, "scope", "", "openid", "comma-separated scopes to request in token")
	getAuthcodeToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
	getAuthcodeToken.Annotations = make(map[string]string)
	getAuthcodeToken.Annotations[TOKEN_CATEGORY] = "true"
	RootCmd.AddCommand(getAuthcodeToken)
}
