package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/vault"
)

var namespacePath string

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Parse and inspect a namespaced Vault secret path",
	RunE:  runNamespace,
}

func init() {
	namespaceCmd.Flags().StringVarP(&namespacePath, "path", "p", "", "Namespaced path (e.g. myns/secret/myapp)")
	_ = namespaceCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(namespaceCmd)
}

func runNamespace(cmd *cobra.Command, args []string) error {
	nc, err := vault.ParseNamespacedPath(namespacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	vaultPath := vault.BuildVaultPath(nc)
	headers := vault.ApplyNamespaceHeader(map[string]string{}, nc.Namespace)

	fmt.Printf("Namespace  : %s\n", nc.Namespace)
	fmt.Printf("Mount Path : %s\n", nc.MountPath)
	fmt.Printf("Secret Path: %s\n", nc.SecretPath)
	fmt.Printf("Vault Path : %s\n", vaultPath)
	fmt.Printf("Header     : X-Vault-Namespace=%s\n", headers["X-Vault-Namespace"])
	return nil
}
