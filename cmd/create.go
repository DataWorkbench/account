/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmds

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an admin account",
	Long:  `create an admin account`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tx := options.DBConn
		adminId, err := options.IdGeneratorUser.Take()
		if err != nil {
			return
		}
		hash := sha256.New()
		hash.Write([]byte("zhu88jie"))
		password := hex.EncodeToString(hash.Sum(nil))
		err = user.CreateAdminUser(tx, adminId, "admin", password, "account@yunify.com")
		if err != nil {
			return
		}
	},
}

func init() {
	root.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
