package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

func main() {

	cobra.OnInitialize(initConfig)
	rand.Seed(time.Now().UTC().UnixNano())

	var rootCmd = &cobra.Command{Use: "zb"}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zetabase/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVarP(&zbHostPort, ConfigKeyZbHostPort, "A", "api.zetabase.io:443", "hostname:port, e.g. grpc.zetabase.io:443")
	viper.BindPFlag(ConfigKeyZbHostPort, rootCmd.PersistentFlags().Lookup(ConfigKeyZbHostPort))

	rootCmd.PersistentFlags().BoolVarP(&connectInsecure, ConfigKeyConnectInsecure, "3", false, "toggle insecure connection")
	viper.BindPFlag(ConfigKeyConnectInsecure, rootCmd.PersistentFlags().Lookup(ConfigKeyConnectInsecure))

	rootCmd.PersistentFlags().BoolVarP(&connectNoCertVerify, ConfigKeyConnectNoCertVerify, "4", true, "toggle no certificate verify")
	viper.BindPFlag(ConfigKeyConnectNoCertVerify, rootCmd.PersistentFlags().Lookup(ConfigKeyConnectNoCertVerify))

	rootCmd.PersistentFlags().StringVarP(&identityFile, ConfigKeyIdentityFile, "i", "", "zetabase.xxx.identity (path to identity file)")
	viper.BindPFlag(ConfigKeyIdentityFile, rootCmd.PersistentFlags().Lookup(ConfigKeyIdentityFile))

	rootCmd.PersistentFlags().StringVarP(&idPassword, ConfigKeyIdPassword, "P", "", "p@ssw0rd (or `prompt` to be prompted)")
	viper.BindPFlag(ConfigKeyIdPassword, rootCmd.PersistentFlags().Lookup(ConfigKeyIdPassword))

	rootCmd.PersistentFlags().StringVarP(&loginId, ConfigKeyLoginId, "I", "", "your handle")
	viper.BindPFlag(ConfigKeyLoginId, rootCmd.PersistentFlags().Lookup(ConfigKeyLoginId))

	rootCmd.PersistentFlags().StringVarP(&loginParentId, ConfigKeyLoginParentId, "R", "", "123f...")
	viper.BindPFlag(ConfigKeyLoginParentId, rootCmd.PersistentFlags().Lookup(ConfigKeyLoginParentId))

	rootCmd.PersistentFlags().BoolVarP(&verbose, ConfigKeyVerbose, "v", false, "turn on/off verbose logging")
	viper.BindPFlag(ConfigKeyVerbose, rootCmd.PersistentFlags().Lookup(ConfigKeyVerbose))

	// View flags
	cmdView.Flags().StringVarP(&tableId, ConfigKeyTableId, "t", "", "mytable")
	viper.BindPFlag(ConfigKeyTableId, cmdView.Flags().Lookup(ConfigKeyTableId))

	cmdView.Flags().StringVarP(&tableOwnerId, ConfigKeyTableOwnerId, "o", "", "123f-...")
	viper.BindPFlag(ConfigKeyTableOwnerId, cmdView.Flags().Lookup(ConfigKeyTableOwnerId))

	cmdView.Flags().StringVarP(&outputDataType, ConfigKeyOutputDataType, "X", "", "one of: binary, json, text")
	viper.BindPFlag(ConfigKeyOutputDataType, cmdView.Flags().Lookup(ConfigKeyOutputDataType))

	cmdView.Flags().StringVarP(&keyPattern, ConfigKeyKeyPattern, "K", "", "prefix/%")
	viper.BindPFlag(ConfigKeyKeyPattern, cmdView.Flags().Lookup(ConfigKeyKeyPattern))

	cmdView.Flags().StringVarP(&exportDataMode, ConfigKeyExportDataMode, "x", "", "flush data to STDOUT in export mode: one of json, jsonwithkey, base64")
	viper.BindPFlag(ConfigKeyExportDataMode, cmdView.Flags().Lookup(ConfigKeyExportDataMode))


	// Create flags
	cmdCreate.Flags().StringVarP(&tableId, ConfigKeyTableId, "t", "", "mytable")
	viper.BindPFlag(ConfigKeyTableId, cmdCreate.Flags().Lookup(ConfigKeyTableId))

	cmdCreate.Flags().BoolVarP(&allowJwt, ConfigKeyCreateAllowTokens, "J", false, "allow JWT token authentication")
	viper.BindPFlag(ConfigKeyCreateAllowTokens, cmdCreate.Flags().Lookup(ConfigKeyCreateAllowTokens))

	cmdCreate.Flags().StringVarP(&createPermissions, ConfigKeyCreatePermissions, "p", "", "see docs for usage")
	viper.BindPFlag(ConfigKeyCreatePermissions, cmdCreate.Flags().Lookup(ConfigKeyCreatePermissions))

	// List flags
	cmdList.Flags().StringVarP(&tableOwnerId, ConfigKeyTableOwnerId, "o", "", "123f-...")
	viper.BindPFlag(ConfigKeyTableOwnerId, cmdList.Flags().Lookup(ConfigKeyTableOwnerId))

	cmdList.Flags().StringVarP(&tableId, ConfigKeyTableId, "t", "", "mytable")
	viper.BindPFlag(ConfigKeyTableId, cmdList.Flags().Lookup(ConfigKeyTableId))

	// Delete flags
	cmdDelete.Flags().StringVarP(&tableOwnerId, ConfigKeyTableOwnerId, "o", "", "123f-...")
	viper.BindPFlag(ConfigKeyTableOwnerId, cmdDelete.Flags().Lookup(ConfigKeyTableOwnerId))

	cmdDelete.Flags().StringVarP(&tableId, ConfigKeyTableId, "t", "", "mytable")
	viper.BindPFlag(ConfigKeyTableId, cmdDelete.Flags().Lookup(ConfigKeyTableId))

	cmdDelete.Flags().StringVarP(&tableKey, ConfigKeyTableKey, "k", "", "path/to/key/0")
	viper.BindPFlag(ConfigKeyTableKey, cmdDelete.Flags().Lookup(ConfigKeyTableKey))


	// Put flags
	cmdPut.Flags().StringVarP(&tableId, ConfigKeyTableId, "t", "", "mytable")
	viper.BindPFlag(ConfigKeyTableId, cmdPut.Flags().Lookup(ConfigKeyTableId))

	cmdPut.Flags().StringVarP(&tableKey, ConfigKeyTableKey, "k", "", "path/to/key/0")
	viper.BindPFlag(ConfigKeyTableKey, cmdPut.Flags().Lookup(ConfigKeyTableKey))

	cmdPut.Flags().StringVarP(&tableValue, ConfigKeyValue, "V", "", "stringdata")
	viper.BindPFlag(ConfigKeyValue, cmdPut.Flags().Lookup(ConfigKeyValue))

	cmdPut.Flags().StringVarP(&tableOwnerId, ConfigKeyTableOwnerId, "o", "", "123f-...")
	viper.BindPFlag(ConfigKeyTableOwnerId, cmdPut.Flags().Lookup(ConfigKeyTableOwnerId))

	cmdPut.Flags().BoolVarP(&putOverwrite, ConfigKeyPutOverwrite, "O", false, "enable overwrite on put")
	viper.BindPFlag(ConfigKeyPutOverwrite, cmdPut.Flags().Lookup(ConfigKeyPutOverwrite))

	cmdPut.Flags().StringVarP(&binInFile, ConfigKeyPutBinInFile, "f", "", "somefile.bin")
	viper.BindPFlag(ConfigKeyPutBinInFile, cmdPut.Flags().Lookup(ConfigKeyPutBinInFile))

	// Manage flags
	cmdManage.Flags().StringVarP(&userEmail, ConfigKeyUserEmail, "", "", "user@domain.com")
	viper.BindPFlag(ConfigKeyUserEmail, cmdManage.Flags().Lookup(ConfigKeyUserEmail))

	cmdManage.Flags().StringVarP(&adminPass, ConfigKeyAdminPassword, "p", "", "p@$$w0rd")
	viper.BindPFlag(ConfigKeyAdminPassword, cmdManage.Flags().Lookup(ConfigKeyAdminPassword))

	cmdManage.Flags().StringVarP(&parentUid, ConfigKeyParentUid, "F", "", "12df...")
	viper.BindPFlag(ConfigKeyParentUid, cmdManage.Flags().Lookup(ConfigKeyParentUid))

	cmdManage.Flags().StringVarP(&userMobile, ConfigKeyUserMobile, "", "", "+19175551212")
	viper.BindPFlag(ConfigKeyUserMobile, cmdManage.Flags().Lookup(ConfigKeyUserMobile))

	cmdManage.Flags().StringVarP(&userName, ConfigKeyUserName, "", "", "User Name")
	viper.BindPFlag(ConfigKeyUserName, cmdManage.Flags().Lookup(ConfigKeyUserName))

	cmdManage.Flags().StringVarP(&userPubKeyFile, ConfigKeyPubKeyFile, "", "", "zetabase.xxx.pub (empty to generate new)")
	viper.BindPFlag(ConfigKeyPubKeyFile, cmdManage.Flags().Lookup(ConfigKeyPubKeyFile))

	cmdManage.Flags().StringVarP(&userPrivKeyFile, ConfigKeyPrivKeyFile, "", "", "zetabase.xxx.priv (empty to generate new)")
	viper.BindPFlag(ConfigKeyPrivKeyFile, cmdManage.Flags().Lookup(ConfigKeyPrivKeyFile))

	cmdManage.Flags().StringVarP(&adminTask, ConfigKeyAdminTask, "t", "", "One of: newuser, newsubuser")
	viper.BindPFlag(ConfigKeyAdminTask, cmdManage.Flags().Lookup(ConfigKeyAdminTask))

	cmdManage.Flags().StringVarP(&subIdSignupCode, ConfigKeySubIdSignupCode, "S", "", "Signup code for new subuser identity (if required)")
	viper.BindPFlag(ConfigKeySubIdSignupCode, cmdManage.Flags().Lookup(ConfigKeySubIdSignupCode))

	cmdManage.Flags().StringVarP(&subIdGroupId, ConfigKeyGroupId, "G", "", "Signup group for new subuser identity (optional)")
	viper.BindPFlag(ConfigKeyGroupId, cmdManage.Flags().Lookup(ConfigKeyGroupId))

	// Shell flags
	// Currently none

	rootCmd.AddCommand(cmdManage)
	rootCmd.AddCommand(cmdView)
	rootCmd.AddCommand(cmdPut)
	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdDelete)
	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdShell)
	rootCmd.Execute()

}
