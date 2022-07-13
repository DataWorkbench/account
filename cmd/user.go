/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmds

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/account/handler/user"
	"github.com/DataWorkbench/account/options"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run the add, init_admin, delete, and reset subcommands")
	},
}

var init_admin = &cobra.Command{
	Use:   "init_admin",
	Short: "init admin",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置文件
		err := options.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 创建用户
		userId, err := options.IdGeneratorUser.Take()
		if err != nil {
			fmt.Println(err)
			return
		}
		hash := sha256.New()
		_, err = hash.Write([]byte("admin"))
		if err != nil {
			return
		}
		passwordWithSHA256 := hex.EncodeToString(hash.Sum(nil))
		err = user.CreateAdminUser(options.DBConn, userId, "admin", passwordWithSHA256, "admin@yunify.com")
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var add = &cobra.Command{
	Use:   "add",
	Short: "add user",
	Long:  `add user,parameters: username password email`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置文件
		err := options.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 创建用户
		userId, err := options.IdGeneratorUser.Take()
		if err != nil {
			fmt.Println(err)
			return
		}
		hash := sha256.New()
		_, err = hash.Write([]byte(args[1]))
		if err != nil {
			return
		}
		passwordWithSHA256 := hex.EncodeToString(hash.Sum(nil))
		if len(args[0]) == 0 || len(args[1]) == 0 {
			fmt.Println("username or password is empty")
		}
		err = user.CreateUser(options.DBConn, userId, args[0], passwordWithSHA256, args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var reset = &cobra.Command{
	Use:   "reset",
	Short: "reset password",
	Long:  `reset pasword ,parameters: username password`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置文件
		err := options.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 创建用户
		hash := sha256.New()
		_, err = hash.Write([]byte(args[1]))
		if err != nil {
			return
		}
		passwordWithSHA256 := hex.EncodeToString(hash.Sum(nil))
		if len(args[0]) == 0 || len(args[1]) == 0 {
			fmt.Println("username or password is empty")
		}
		err = user.ResetPasswordByName(options.DBConn, args[0], passwordWithSHA256)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var delete = &cobra.Command{
	Use:   "delete",
	Short: "delete user",
	Long:  `delete user,parameters: username`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置文件
		err := options.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = user.DeleteUserByNames(options.DBConn, args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	root.AddCommand(userCmd)

	userCmd.AddCommand(init_admin)

	userCmd.AddCommand(reset)

	userCmd.AddCommand(add)

	userCmd.AddCommand(delete)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	userCmd.PersistentFlags().StringVarP(&config.FilePath, "file", "f", "", "config file path")

}
