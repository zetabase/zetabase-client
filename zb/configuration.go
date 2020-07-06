package main

import (
	"crypto/ecdsa"
	"github.com/mitchellh/go-homedir"
	"github.com/zetabase/zetabase-client"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"strings"
)

const (

	ConfigKeyUserEmail     = "user.email"
	ConfigKeyUserMobile    = "user.mobile"
	ConfigKeyUserName      = "user.name"
	ConfigKeyPubKeyFile    = "user.pubkey-file"
	ConfigKeyPrivKeyFile   = "user.privkey-file"
	ConfigKeyAdminTask     = "admin.task"
	ConfigKeyAdminPassword = "admin.pass"

	ConfigKeySubIdSignupCode = "admin.signupcode"
	ConfigKeyGroupId = "admin.groupid"

	AdminTaskNewUser       = "newuser"
	AdminTaskNewSubUser    = "newsubuser"
	AdminTaskTestingGround = "test"
	AdminTaskListSubUsers = "listusers"
	AdminTaskTestClient    = "client"

	ConfigKeyIdentityFile   = "identity"
	ConfigKeyTableId        = "table"
	ConfigKeyTableOwnerId   = "owner"
	ConfigKeyTableKey       = "key"
	ConfigKeyValue          = "value"
	ConfigKeyOutputDataType = "dtype"
	ConfigKeyParentUid      = "parent"
	ConfigKeyKeyPattern     = "pattern"

	ConfigKeyZbHostPort          = "host"
	ConfigKeyVerbose             = "verbose"
	ConfigKeyConnectInsecure     = "insecure"
	ConfigKeyConnectNoCertVerify = "nocertverify"

	ConfigKeyCreatePermissions = "permissions"
	ConfigKeyCreateAllowTokens = "allowjwt"
	ConfigKeyPutOverwrite      = "overwrite"
	ConfigKeyPutBinInFile      = "inputfile"

	ConfigKeyIdPassword    = "password"
	ConfigKeyLoginId       = "loginid"
	ConfigKeyLoginParentId = "loginparentid"

	ConfigKeyExportDataMode = "mode.export"
)

var (
	zbHostPort          = ""
	cfgFile             = ""
	identityFile        = ""
	idPassword          = ""
	loginId             = ""
	loginParentId       = ""
	subIdSignupCode     = ""
	subIdGroupId        = ""
	binInFile           = ""
	keyPattern          = ""
	tableId             = ""
	tableOwnerId        = ""
	tableKey            = ""
	outputDataType      = ""
	tableValue          = ""
	userEmail           = ""
	userMobile          = ""
	userName            = ""
	userPubKeyFile      = ""
	userPrivKeyFile     = ""
	adminTask           = ""
	adminPass           = ""
	createPermissions   = ""
	allowJwt            = false
	connectInsecure     = false
	connectNoCertVerify = false
	verbose             = false
	parentUid           = ""
	putOverwrite        = false
	exportDataMode        = ""
)

type IdentityDefinition struct {
	Id         string `json:"id"`
	ParentId   string `json:"parent_id,omitempty"`
	PubKeyEnc  string `json:"pub_key"`
	PrivKeyEnc string `json:"priv_key"`
}

func (d *IdentityDefinition) ToUserIdentity() (*UserIdentity, error) {
	pub, err := zetabase.DecodeEcdsaPublicKey(d.PubKeyEnc)
	if err != nil {
		return nil, err
	}
	priv, err := zetabase.DecodeEcdsaPrivateKey(d.PrivKeyEnc)
	if err != nil {
		return nil, err
	}
	var parentPtr *string
	if len(d.ParentId) > 0 {
		parentPtr = &d.ParentId
	}
	return &UserIdentity{
		Id:       d.Id,
		ParentId: parentPtr,
		PubKey:   pub,
		PrivKey:  priv,
	}, nil
}

type UserIdentity struct {
	Id       string
	ParentId *string
	PubKey   *ecdsa.PublicKey
	PrivKey  *ecdsa.PrivateKey
}

func (u *UserIdentity) MakeSignature(nonce int64, relBytes []byte) *zbprotocol.EcdsaSignature {
	r, s := zetabase.MakeZetabaseSignature(u.Id, nonce, relBytes, u.PrivKey)
	return &zbprotocol.EcdsaSignature{
		R: r,
		S: s,
	}
}

func GetVersionString() string {
	return zetabase.ClientVersion
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return //er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.zetabase")
		viper.AddConfigPath("/etc/zetabase")
		viper.AddConfigPath(".")
		//viper.SetConfigName(".cobra")
	}

	viper.SetEnvPrefix("zb")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func isWindows() bool {
	s := runtime.GOOS
	if strings.Contains(strings.ToLower(s), "windows") {
		return true
	}
	return false
}


