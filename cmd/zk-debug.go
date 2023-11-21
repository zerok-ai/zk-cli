package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zkctl/cmd/internal"
	"zkctl/cmd/internal/shell"
)

var (
	debugInstallCmd = &cobra.Command{
		Use:   "debug",
		Short: "Install ZeroK",
		RunE:  RunDebugInstallCmd,
	}

	debugDeleteCmd = &cobra.Command{
		Use:   "debug",
		Short: "Install ZeroK",
		RunE:  RunDebugDeleteCmd,
	}

	debugShellFile = "scripts/debug.sh"

	zksSpinnerInstallationText      = "enabling access to debug data"
	zksDebugInstallationSuccessText = "debug data access enabled successfully"
	zksDebugInstallationFailureText = "failed to enable debug data access"

	zksSpinnerDeletionText      = "disabling debug data access"
	zksDebugDeletionSuccessText = "debug data access disabled successfully"
	zksDebugDeletionFailureText = "failed to disable debug data access"
)

const (
	postgresUser = "postgres"
	postgresHost = "postgres.zk-client.svc.cluster.local"
	postgresPort = "5432"
	postgresDB   = "pl"
)

func RunDebugInstallCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("Running debug install command")

	redisPwdCmd := "kubectl get secret --namespace zk-client redis -o jsonpath=\"{.data.redis-password}\" | base64 -d"
	redisPwd, err := shell.Shellout(redisPwdCmd)
	if err != nil {
		internal.DumpErrorAndPrintLocation(fmt.Sprintf("Error in getting credentials 1 : %v", err))
		return err
	}

	postgresPwdCmd := "kubectl get secret --namespace zk-client postgres -o jsonpath=\"{.data.postgres-password}\" | base64 -d"
	postgresPwd, err := shell.Shellout(postgresPwdCmd)
	if err != nil {
		internal.DumpErrorAndPrintLocation(fmt.Sprintf("Error in getting credentials 2 : %v", err))
		return err
	}

	// Get base64 encoded postgres url
	urlPostgres := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgresUser, postgresPwd, postgresHost, postgresPort, postgresDB)
	urlPostgresBase64 := base64.StdEncoding.EncodeToString([]byte(urlPostgres))

	inputToShellFile := fmt.Sprintf("install=%s pgurl=%s redisAuth=%s", "true", urlPostgresBase64, redisPwd)
	if viper.Get(internal.EmbedKeyFlag) == false {
		return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+debugShellFile, inputToShellFile, zksSpinnerInstallationText, zksDebugInstallationSuccessText, zksDebugInstallationFailureText)
	} else {
		return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, debugShellFile, inputToShellFile, zksSpinnerInstallationText, zksDebugInstallationSuccessText, zksDebugInstallationFailureText)
	}

}

func RunDebugDeleteCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("Running debug delete command")
	inputToShellFile := "delete=true"
	if viper.Get(internal.EmbedKeyFlag) == false {
		return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+debugShellFile, inputToShellFile, zksSpinnerDeletionText, zksDebugDeletionSuccessText, zksDebugDeletionFailureText)
	} else {
		return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, debugShellFile, inputToShellFile, zksSpinnerDeletionText, zksDebugDeletionSuccessText, zksDebugDeletionFailureText)
	}
}

func init() {
	installCmd.AddCommand(debugInstallCmd)
	deleteCmd.AddCommand(debugDeleteCmd)
}
