package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/johnsiilver/getcert"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zetabase/zetabase-client"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func cleanStringForFilename(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func validatePhoneNumber(s string) bool {

	return zetabase.ValidatePhoneNumber(s)
}

func parseIdentFromFile(fn string) *UserIdentity {
	var obj IdentityDefinition
	bs, e := ioutil.ReadFile(fn)
	if e != nil {
		Logf("Failed to read ID file: %s", e.Error())
		return nil
	}
	e = json.Unmarshal(bs, &obj)
	if e != nil {
		Logf("Failed to parse ID file: %s", e.Error())
		return nil
	}
	id, e := obj.ToUserIdentity()
	if e != nil {
		Logf("Failed to decode ID file: %s", e.Error())
		return nil
	}
	return id
}


func listTables(identity *UserIdentity, tblOwnerId string, nonce int64, poc *zbprotocol.ProofOfCredential, client zbprotocol.ZetabaseProviderClient) ([]*zbprotocol.TableCreate, error) {

	tlr := &zbprotocol.ListTablesRequest{
		Id:           identity.Id,
		TableOwnerId: tblOwnerId,
		Nonce:        nonce,
		Credential:   poc,
	}

	if isVerbose() {
		Logf("Listing tables owned by %s (as identity %s)", tblOwnerId, identity.Id)
		if poc.CredType == zbprotocol.CredentialProofType_JWT_TOKEN {
			Logf("\tUsing token: %s", poc.GetJwtToken())
		}
	}
	resp, err := client.ListTables(context.Background(), tlr)
	if err != nil {
		return nil, err
	}
	return resp.GetTableDefinitions(), nil
}

func deleteKey(identity *UserIdentity, tblOwnerId, tbl, key string, nonce int64, poc *zbprotocol.ProofOfCredential, client zbprotocol.ZetabaseProviderClient) error {
	req := &zbprotocol.DeleteSystemObjectRequest{
		Id:           identity.Id,
		ObjectType:   zbprotocol.SystemObjectType_KEY,
		TableOwnerId: tblOwnerId,
		TableId:      tbl,
		ObjectId:     key,
		Nonce:        nonce,
		Credential:   poc,
	}
	_, err := client.DeleteObject(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func deleteTable(identity *UserIdentity, tblOwnerId, tbl string, nonce int64, poc *zbprotocol.ProofOfCredential, client zbprotocol.ZetabaseProviderClient) error {
	req := &zbprotocol.DeleteSystemObjectRequest{
		Id:           identity.Id,
		ObjectType:   zbprotocol.SystemObjectType_TABLE,
		TableOwnerId: tblOwnerId,
		TableId:      tbl,
		ObjectId:     tbl,
		Nonce:        nonce,
		Credential:   poc,
	}
	_, err := client.DeleteObject(context.Background(), req)
	if err != nil {
		return err
	}
	return nil

}

var cmdDelete = &cobra.Command{
	Use:   "rm",
	Short: "Delete keys, users, and tables",
	Long:  `Delete any system object with rm key <key>, rm table <table>, rm subuser <userid>.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identity := loadIdentityFromConfigs()
		nonceMaker := zetabase.NewNonceMaker()
		nonce := nonceMaker.Get()
		tblOwnerId := chooseDefaultTableOwnerId(identity)
		tbl := viper.GetString(ConfigKeyTableId)
		conn := dialRemote()
		defer conn.Close()

		client := zbprotocol.NewZetabaseProviderClient(conn)

		var extraBytes []byte
		var req *zbprotocol.DeleteSystemObjectRequest

		switch strings.ToLower(args[0]) {
		case "key":
			key := viper.GetString(ConfigKeyTableKey)
			extraBytes = []byte(key)
			req = &zbprotocol.DeleteSystemObjectRequest{
				Id:           identity.Id,
				ObjectType:   zbprotocol.SystemObjectType_KEY,
				TableOwnerId: tblOwnerId,
				TableId:      tbl,
				ObjectId:     key,
				Nonce:        0,
				Credential:   nil,
			}
		case "table":
			extraBytes = []byte(tbl)
			req = &zbprotocol.DeleteSystemObjectRequest{
				Id:           identity.Id,
				ObjectType:   zbprotocol.SystemObjectType_TABLE,
				TableOwnerId: tblOwnerId,
				TableId:      tbl,
				ObjectId:     tbl,
				Nonce:        0,
				Credential:   nil,
			}
		default:
			Logf("Unknown object type for `rm`: '%s' -- should be table, key, or subuser", args[0])
			return
		}

		uid, poc, err := getUserCredential(identity, nonce, extraBytes, client)
		if err != nil {
			PrintErrorAndQuit(err)
		}

		// TODO: these clauses can be removed as they are redundant to `chooseDefaultTableOwnerId` (verify before removing)
		if len(identity.Id) == 0 || identity.PrivKey == nil {
			identity.Id = uid
			if len(tblOwnerId) == 0 {
				tblOwnerId = identity.Id
			}
		}

		if req != nil {
			// Inject credentials
			req.Nonce = nonce
			req.Credential = poc

			res, err := client.DeleteObject(context.Background(), req)
			if err != nil {
				PrintErrorAndQuit(err)
			} else {
				Logf("Request successful: code %d", res.GetCode())
			}
		} else {
			Logf("No request to make! Please check your syntax.")
		}

	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List owned tables",
	Long:  `List all tables created by the given identity.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		qtStr := ""
		if len(args) > 0 {
			qtStr = args[0]
		}
		if len(qtStr) == 0 || (qtStr != "keys" && qtStr != "tables") {
			Logf("Usage is: zb list (tables|keys) [prefix/%%] [-t tbl] [-o ownerId]")
			return
		}
		identity := loadIdentityFromConfigs()
		nonceMaker := zetabase.NewNonceMaker()
		nonce := nonceMaker.Get()
		tblOwnerId := chooseDefaultTableOwnerId(identity)
		conn := dialRemote()
		defer conn.Close()

		client := zbprotocol.NewZetabaseProviderClient(conn)

		//credential := zetabase.MakeCredentialEcdsa(nonce, identity.Id, nil, identity.PrivKey)

		uid, poc, err := getUserCredential(identity, nonce, nil, client)
		if err != nil {
			//Logf("Error 1.")
			PrintErrorAndQuit(err)
		}

		if len(identity.Id) == 0 || identity.PrivKey == nil {
			//Logf("No ID set -- using %s...", uid)
			identity.Id = uid
			if len(tblOwnerId) == 0 {
				tblOwnerId = identity.Id
			}
		}

		if qtStr == "tables" {
			//Logf("Listing tables for %s by %s...", identity.Id, tblOwnerId)
			tblDefs, err := listTables(identity, tblOwnerId, nonce, poc, client)
			/*tlr := &zetabase.ListTablesRequest{
				Id:           identity.Id,
				TableOwnerId: tblOwnerId,
				Nonce:        nonce,
				Credential:   poc,
			}

			resp, err := client.ListTables(context.Background(), tlr) */
			if isVerbose() {
				Logf("Listing tables owned by %s (as identity %s)", tblOwnerId, identity.Id)
			}
			if err != nil {
				PrintErrorAndQuit(err)
			} else {
				PrintTableDefinitions(tblDefs)
			}
		} else {
			// Keys
			pat := ""
			patExplanation := ""
			if len(args) > 1 {
				pat = args[1]
				patExplanation = " with prefix pattern " + pat
			}

			tblId := viper.GetString(ConfigKeyTableId)
			if isVerbose() {
				Logf("Looking for keys on table %s (owned by %s)%s (as %s)", tblId, tblOwnerId, patExplanation, identity.Id)
			}
			keyList, err := listKeys(identity, tblId, tblOwnerId, pat, nonce, poc, client)
			/*lkr := &zetabase.ListKeysRequest{
				Id:           identity.Id,
				TableId:      tblId,
				TableOwnerId: tblOwnerId,
				Pattern:      pat,
				Nonce:        nonce,
				Credential:   poc,
			}
			resp, err := client.ListKeys(context.Background(), lkr)*/
			if err != nil {
				PrintErrorAndQuit(err)
			} else {
				// print!
				//Logf("d", "%v", resp)
				PrintSingleColumn("Key", keyList)
				//PrintTableDefinitions(resp.TableDefinitions)
			}
		}
	},
}

var cmdShell = &cobra.Command{
	Use:   "shell",
	Short: "Start a shell session",
	Long:  `Enter an interactive (Zetabase shell) session.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		identity := loadIdentityFromConfigs()
		Logf("Welcome to Zetabase (client v%s %s).", GetVersionString(), zetabase.ClientVersionCode)
		cli := makeNewClient(identity.Id, identity.PrivKey, identity.PubKey)
		if identity != nil && len(identity.Id) > 0 {
			Logf("\tUsing identity: %s\n", identity.Id)
		} else {
			// Handle for JWT token (root account)
			usrHandl := viper.GetString(ConfigKeyLoginId)
			usrPass := viper.GetString(ConfigKeyIdPassword)
			if len(usrHandl) > 0 && len(usrPass) > 0 {
				// Do handle-based login
				cli.SetIdPassword(usrHandl, usrPass)
			} else {
				Logf("No identity found. Please supply an identity with the `-i` switch, via environment variable, or with -I and -P flags (username/password JWT auth).")
			}
		}

		// Check versions and compatibility
		versionOk, versionInfo, err := cli.CheckVersion()
		if err != nil {
			PrintErrorAndQuit(err)
		} else if !versionOk {
			Logf("Your client version %s is out of date (minimum version required: %s). Please go to zetabase.io and update.", zetabase.ClientVersion, versionInfo.GetMinClientVersion())
		} else {
			Logf("Connected to server v%s (client v%s)...", versionInfo.GetServerVersion(), zetabase.ClientVersion)
			if !zetabase.IsSemVerVersionAtLeast(zetabase.ClientVersion, versionInfo.GetClientVersion()) {
				Logf("\tNew client available: v%s (see zetabase.io to update)", versionInfo.GetClientVersion())
			}
		}

		// Check if Windows

		var e error
		if isWindows() {
			Logf("Detected Windows OS. Disabling auto-complete, command history, and advanced styles...")
			e = WindowsShellLoop("zb> ", "quit", identity, cli)
		} else {
			// Run normal Linux/Mac version
			e = ShellLoop("zb> ", "quit", identity, cli)
		}
		if e != nil {
			PrintErrorAndQuit(e)
		}
	},
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new table",
	Long:  `Create a new table from a template.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identity := loadIdentityFromConfigs()
		tbl := viper.GetString(ConfigKeyTableId)
		typ := strings.ToLower(args[0])
		nonceMaker := zetabase.NewNonceMaker()

		conn := dialRemote()
		defer conn.Close()
		client := zbprotocol.NewZetabaseProviderClient(conn)

		//credential := zetabase.MakeCredentialEcdsa(nonce, identity.Id, []byte(tableId), identity.PrivKey)
		nonce := nonceMaker.Get()
		uid, poc, err := getUserCredential(identity, nonce, []byte(tableId), client)
		if err != nil {
			PrintErrorAndQuit(err)
		}

		if len(identity.Id) == 0 || identity.PrivKey == nil {
			identity.Id = uid
		}

		permsRaw := viper.GetString(ConfigKeyCreatePermissions)
		//resp, err := createTable(permsRaw, identity, tbl, typ, args[1:], nonce, poc, client)
		ctr := getTableCreate(permsRaw, identity, tbl, typ, args[1:], nonce, poc)
		sigBytes := zetabase.TableCreateSigningBytes(ctr.GetTableId(), ctr.GetPermissions())
		_, poc, err = getUserCredential(identity, nonce, sigBytes, client)
		if err != nil {
			PrintErrorAndQuit(err)
		}
		ctr.Credential = poc
		resp, err := createTable(ctr, client)

		/*ctr := getTableCreate(permsRaw, identity, tbl, typ, args[1:], nonce, poc)

		//Logf("Created object: %v\n", ctr)

		resp, err := client.CreateTable(context.Background(), ctr)*/
		if err != nil {
			Logf("Error creating table: %s", err.Error())
		} else if resp != nil {
			if resp.Code == 0 {
				Logf("Successfully created table.")
			} else {
				Logf("Error creating table [2]: %s [%d]", resp.Message, resp.Code)
			}
		} else {
			Logf("Error creating table [3]: nil result")
		}
	},
}

func createTable(ctr *zbprotocol.TableCreate, client zbprotocol.ZetabaseProviderClient) (*zbprotocol.ZbError, error) {
	resp, err := client.CreateTable(context.Background(), ctr)
	return resp, err
}

func getTableCreate(permsRaw string, identity *UserIdentity, tblId string, typ string, idxArgs []string, nonce int64, poc *zbprotocol.ProofOfCredential) *zbprotocol.TableCreate {
	var perms []*zbprotocol.PermissionsEntry
	var idxFields []string
	var idxTypes []string
	var tdf zbprotocol.TableDataFormat
	if typ == "json" {
		tdf = zbprotocol.TableDataFormat_JSON
		if len(idxArgs)%2 != 0 {
			Logf("Error: incorrect index specification, should be:\n\tfield1 op1 field2 op2... fieldn opn\nwhere each `op` can be one of `lex`, `text`, `real`, or `natural`.")
			return nil
		}
		for i := 0; i < len(idxArgs); i += 2 {
			fld := strings.ToLower(idxArgs[i])
			op := strings.ToLower(idxArgs[i+1])
			idxFields = append(idxFields, fld)
			idxTypes = append(idxTypes, op)
		}
	} else if typ != "binary" && typ != "text" {
		Logf("Error: unknown type %s, should be one of: `json`, `binary`, or `text`.", typ)
		return nil
	} else {
		switch typ {
		case "binary":
			tdf = zbprotocol.TableDataFormat_BINARY
		default:
			tdf = zbprotocol.TableDataFormat_PLAIN_TEXT
		}
	}
	if isVerbose() {
		Logf("User %s creating table %s with index fields %v...", identity.Id, tblId, idxFields)
	}
	//store2998@gmail.com

	// Parse index fields input
	var tblIdxFields []*zbprotocol.TableIndexField
	for i := 0; i < len(idxFields); i++ {
		fld := idxFields[i]
		iTyp := idxTypes[i]
		var ordering zbprotocol.QueryOrdering
		var langCode string
		if strings.HasPrefix(iTyp, "real" ) || strings.HasPrefix(iTyp, "num") {
			ordering = zbprotocol.QueryOrdering_REAL_NUMBERS
		} else if strings.HasPrefix(iTyp, "natural") {
			ordering = zbprotocol.QueryOrdering_INTEGRAL_NUMBERS
		} else if strings.HasPrefix(iTyp, "text") {
			ordering = zbprotocol.QueryOrdering_FULL_TEXT
			arr := strings.Split(iTyp, ":")
			if len(arr) >= 2 {
				langCode = arr[1]
			}
		} else {
			ordering = zbprotocol.QueryOrdering_LEXICOGRAPHIC
		}
		tif := &zbprotocol.TableIndexField{
			Field:    fld,
			Ordering: ordering,
			LanguageCode: langCode,
		}
		tblIdxFields = append(tblIdxFields, tif)
	}

	// Wrap index fields data for zbprotocol
	tblIdxFieldsWrap := &zbprotocol.TableIndexFields{
		Fields: tblIdxFields,
	}
	permsSep := ","
	if isWindows() {
		permsSep = ";"
		Logf("Warning: on Windows, using separator ';' for permissions string (instead of usual separator ',').")
	}
	permsRawArr := strings.Split(permsRaw, permsSep)
	if isVerbose() {
		Logf("Parsing permissions: %v", permsRawArr)
	}
	for i := 0; i < len(permsRawArr) && len(permsRaw) > 0; i++ {
		arr := strings.Split(strings.TrimSpace(permsRawArr[i]), " ")
		if len(arr) < 2 {
			Logf("Error: incorrect permissions specification, should be `perm1,perm2,perm3` where each permi is:\n\ttype level <audience> [constraint field] [constraint value (e.g. @uid)].\n\ttype is one of public, user, single\n\tlevel is one of read, append, delete, admin\n\taudience is user ID if type single, otherwise `_`.")
			return nil
		}
		typ := strings.TrimSpace(arr[0])
		lvl := strings.TrimSpace(arr[1])
		aud := ""
		if len(arr) > 2 {
			aud = strings.TrimSpace(arr[2])
		}
		var at zbprotocol.PermissionAudienceType
		var lv zbprotocol.PermissionLevel

		switch typ {
		case "public":
			at = zbprotocol.PermissionAudienceType_PUBLIC
		case "single":
			at = zbprotocol.PermissionAudienceType_INDIVIDUAL
		case "user":
			at = zbprotocol.PermissionAudienceType_USER
		default:
			Logf("Error: incorrect permissions specification, unknown qualifier: %s.", typ)
			return nil
		}

		switch lvl {
		case "read":
			lv = zbprotocol.PermissionLevel_READ
		case "append":
			lv = zbprotocol.PermissionLevel_APPEND
		case "delete":
			lv = zbprotocol.PermissionLevel_DELETE
		case "admin":
			lv = zbprotocol.PermissionLevel_ADMINISTER
		default:
			Logf("Error: incorrect permissions specification, unknown qualifier: %s.", lvl)
			return nil
		}

		var permConstraints []*zbprotocol.PermissionConstraint

		for i := 3; i < len(arr); i += 2 {
			fld, valu := arr[i], arr[i+1]
			//Logf("Field, value (%d, %d) = %s, %s", i, i+1, fld, valu)
			// Check if value is a special identifier...
			cTyp := zbprotocol.PermissionConstraintType_KEY_PATTERN
			if fld != "@key" {
				cTyp = zbprotocol.PermissionConstraintType_FIELD
			}

			if valu == "@uid" {
				Logf("Adding permission constraint: field %s must contain user's ID...", fld)
				permConstraints = append(permConstraints, &zbprotocol.PermissionConstraint{
					ConstraintType: cTyp,
					FieldConstraint: &zbprotocol.FieldConstraint{
						ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
						FieldKey:       fld,
						ValueType:      zbprotocol.FieldConstraintValueType_UID,
						RequiredValue:  "",
					}, KeyConstraint: &zbprotocol.KeyPatternConstraint{
						ConstraintType:       zbprotocol.FieldConstraintType_EQUALS_VALUE,
						RequiredPrefix:       "",
						RequiredSuffix:       "",
						ValueType:            zbprotocol.FieldConstraintValueType_UID,
						RequiredValue:        "",
					},
				})
			} else if valu == "@order" {
				Logf("Natural order field: %s", fld)
				permConstraints = append(permConstraints, &zbprotocol.PermissionConstraint{
					ConstraintType: cTyp,
					FieldConstraint:      &zbprotocol.FieldConstraint{
						ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
						FieldKey:       fld,
						ValueType:      zbprotocol.FieldConstraintValueType_NATURAL_ORDER,
						RequiredValue:  "",
					}, KeyConstraint: &zbprotocol.KeyPatternConstraint{
						ConstraintType:       zbprotocol.FieldConstraintType_EQUALS_VALUE,
						RequiredPrefix:       "",
						RequiredSuffix:       "",
						ValueType:            zbprotocol.FieldConstraintValueType_NATURAL_ORDER,
						RequiredValue:        "",
					},
				})
			} else if valu == "@time" {
				Logf("Timestamp field: %s", fld)
				permConstraints = append(permConstraints, &zbprotocol.PermissionConstraint{
					ConstraintType: cTyp,
					FieldConstraint:      &zbprotocol.FieldConstraint{
						ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
						FieldKey:       fld,
						ValueType:      zbprotocol.FieldConstraintValueType_TIMESTAMP,
						RequiredValue:  "",
					}, KeyConstraint: &zbprotocol.KeyPatternConstraint{
						ConstraintType:       zbprotocol.FieldConstraintType_EQUALS_VALUE,
						RequiredPrefix:       "",
						RequiredSuffix:       "",
						ValueType:            zbprotocol.FieldConstraintValueType_TIMESTAMP,
						RequiredValue:        "",
					},
				})
			} else if valu == "@random" {
				Logf("Random field: %s", fld)
				permConstraints = append(permConstraints, &zbprotocol.PermissionConstraint{
					ConstraintType: cTyp,
					FieldConstraint:      &zbprotocol.FieldConstraint{
						ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
						FieldKey:       fld,
						ValueType:      zbprotocol.FieldConstraintValueType_RANDOM,
						RequiredValue:  "",
					}, KeyConstraint: &zbprotocol.KeyPatternConstraint{
						ConstraintType:       zbprotocol.FieldConstraintType_EQUALS_VALUE,
						RequiredPrefix:       "",
						RequiredSuffix:       "",
						ValueType:            zbprotocol.FieldConstraintValueType_RANDOM,
						RequiredValue:        "",
					},
				})
			} else {
				Logf("Adding permission constraint: field %s must contain constant value %s...", fld, valu)
				if len(valu) == 0 {
					PrintErrorStringAndQuit("Empty valued constraints not supported.")
				}
				permConstraints = append(permConstraints, &zbprotocol.PermissionConstraint{
					ConstraintType: cTyp,
					FieldConstraint: &zbprotocol.FieldConstraint{
						ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
						FieldKey:       fld,
						ValueType:      zbprotocol.FieldConstraintValueType_CONSTANT,
						RequiredValue:  valu,
					}, KeyConstraint: &zbprotocol.KeyPatternConstraint{
						ConstraintType:       zbprotocol.FieldConstraintType_EQUALS_VALUE,
						RequiredPrefix:       "",
						RequiredSuffix:       "",
						ValueType:            zbprotocol.FieldConstraintValueType_CONSTANT,
						RequiredValue:        "",
					},
				})
			}


		}

		p := &zbprotocol.PermissionsEntry{
			Id:           identity.Id,
			TableId:      tblId,
			AudienceType: at,
			AudienceId:   aud,
			Level:        lv,
			Nonce:        nonce,
			Credential:   poc,
			Constraints:  permConstraints,
		}
		perms = append(perms, p)
	}
	// Form TableCreate message
	ctr := &zbprotocol.TableCreate{
		Id:             identity.Id,
		TableId:        tblId,
		DataFormat:     tdf,
		Indices:        tblIdxFieldsWrap,
		AllowTokenAuth: true,
		Nonce:          nonce,
		Credential:     poc,
		Permissions:    perms,
	}
	return ctr
}

func getUserCredential(identity *UserIdentity, nonce int64, extraSigningBytes []byte, client zbprotocol.ZetabaseProviderClient) (string, *zbprotocol.ProofOfCredential, error) {
	loginUid := viper.GetString(ConfigKeyLoginId)
	loginPass := viper.GetString(ConfigKeyIdPassword)
	parentId := viper.GetString(ConfigKeyLoginParentId)
	validLoginPass := len(loginUid) > 0 && len(loginPass) > 0
	if validLoginPass {

		if isVerbose() {
			Logf("No keys provided: logging in as user `%s`...", loginUid)
		}
		resp, err := client.LoginUser(context.Background(), &zbprotocol.AuthenticateUser{
			ParentId:   parentId, // if empty string, will try to log in as the root user
			Handle:     loginUid,
			Password:   loginPass,
			Nonce:      nonce,
			Credential: zetabase.MakeEmptyCredentials(),
		})
		if err != nil {
			Logf("Error logging in user: %s", err.Error())
			return "", nil, err
		} else {
			//Logf("Got JWT response: %s / %s", resp.GetId(), resp.GetJwtToken())
			if len(resp.GetId()) > 0 && len(resp.GetJwtToken()) > 0 {
				tokCred := zetabase.MakeCredentialJwt(resp.GetJwtToken())
				return resp.GetId(), tokCred, nil
			} else {
				return "", nil, errors.New("Unknown")
			}
		}
	} else {
		if identity == nil {
			return "", nil, errors.New("NoIdentityProvided")
		} else if identity.PrivKey == nil {
			return "", nil, errors.New("NoKeyProvided")
		} else if len(identity.Id) == 0 {
			return "", nil, errors.New("WrongIdProvided")
		}
		sig := zetabase.MakeCredentialEcdsa(nonce, identity.Id, extraSigningBytes, identity.PrivKey)
		return identity.Id, sig, nil
	}
}

func chooseDefaultTableOwnerId(identity *UserIdentity) string {
	tblOwnerId := viper.GetString(ConfigKeyTableOwnerId)
	if len(tblOwnerId) == 0 {
		if identity.ParentId != nil {
			tblOwnerId = *identity.ParentId
		} else {
			tblOwnerId = identity.Id
		}
	}
	return tblOwnerId
}

var cmdPut = &cobra.Command{
	Use:   "put",
	Short: "Put a new piece of data",
	Long:  `Insert a new entry into a table.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		identity := loadIdentityFromConfigs()
		tbl := viper.GetString(ConfigKeyTableId)
		dKey := viper.GetString(ConfigKeyTableKey)
		dValuS := viper.GetString(ConfigKeyValue)
		dValu := []byte(dValuS)
		dValueFn := viper.GetString(ConfigKeyPutBinInFile)
		if len(dValueFn) > 0 {
			bs, err := ioutil.ReadFile(dValueFn)
			if err != nil {
				PrintErrorAndQuit(err)
			}
			dValu = bs
		}

		doOverwrite := viper.GetBool(ConfigKeyPutOverwrite)
		tblOwnerId := chooseDefaultTableOwnerId(identity)
		//tblOwnerId := viper.GetString(ConfigKeyTableOwnerId)
		//if len(tblOwnerId) == 0 {
		//	if identity.ParentId != nil {
		//		tblOwnerId = *identity.ParentId
		//	} else {
		//		tblOwnerId = identity.Id
		//	}
		//}
		conn := dialRemote()
		defer conn.Close()
		client := zbprotocol.NewZetabaseProviderClient(conn)
		nonceMaker := zetabase.NewNonceMaker()
		nonce := nonceMaker.Get()

		xBytes := zetabase.TablePutExtraSigningBytes(dKey, dValu)

		uid, poc, err := getUserCredential(identity, nonce, xBytes, client)
		if err != nil {
			PrintErrorAndQuit(err)
		}

		if len(identity.Id) == 0 || identity.PrivKey == nil {
			identity.Id = uid
			if len(tblOwnerId) == 0 {
				tblOwnerId = identity.Id
			}
		}

		_, err = client.PutData(context.Background(), &zbprotocol.TablePut{
			Id:           identity.Id,
			TableOwnerId: tblOwnerId,
			TableId:      tbl,
			Key:          dKey,
			Value:        []byte(dValu),
			Overwrite:    doOverwrite,
			Nonce:        nonce,
			//Credential:   zetabase.MakeCredentialEcdsa(nonce, identity.Id, []byte(dValu), identity.PrivKey),
			Credential: poc,
		})
		if err != nil {
			PrintErrorAndQuit(err)
		} else {
			Logf("Success.")
		}
	},
}

var cmdView = &cobra.Command{
	Use:   "view",
	Short: "View and export",
	Long:  `View and export data from databases, including queries.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		identity := loadIdentityFromConfigs()
		tbl := viper.GetString(ConfigKeyTableId)
		if len(tbl) == 0 {
			PrintErrorStringAndQuit("Please specify a table name (e.g. with `-t tablename`).")
		}
		tblOwner := chooseDefaultTableOwnerId(identity)
		if tblOwner != identity.Id {
			if isVerbose() {
				Logf("Using table owner: %s", tblOwner)
			}
		}

		cli := makeNewClient(identity.Id, identity.PrivKey, identity.PubKey)

		//conn := dialRemote()
		//defer conn.Close()
		//client := zetabase.NewZetabaseProviderClient(conn)
		client := cli.GrpcClient()
		nonceMaker := zetabase.NewNonceMaker()
		nonce := nonceMaker.Get()
		var searchKeys []string
		if len(args) > 0 {
			searchKeys = args
		}

		uid, poc, err := getUserCredential(identity, nonce, nil, client)
		if err != nil {
			PrintErrorAndQuit(err)
		}

		oldId := &UserIdentity{
			Id:       identity.Id,
			ParentId: identity.ParentId,
			PubKey:   identity.PubKey,
			PrivKey:  identity.PrivKey,
		} // save for later if needed:
		if len(identity.Id) == 0 || identity.PrivKey == nil {
			identity.Id = uid
			if len(tblOwner) == 0 {
				tblOwner = identity.Id
			}
		}

		keyPattern := viper.GetString(ConfigKeyKeyPattern)
		if len(keyPattern) > 0 {
			keys, err := listKeys(identity, tbl, tblOwner, keyPattern, nonce, poc, client)
			if err != nil {
				PrintErrorAndQuit(err)
			} else {
				searchKeys = append(searchKeys, keys...)
			}

			// provide new creds since we used the old ones
			nonce = nonceMaker.Get()
			_, poc, err = getUserCredential(oldId, nonce, nil, client)
			if err != nil {
				PrintErrorAndQuit(err)
			}
		}

		if isVerbose() {
			verbKeys := searchKeys
			/*if len(verbKeys) > 15 {
				verbKeys = searchKeys[:15]
			}*/
			Logf("Searching keys: \t%s", strings.Join(verbKeys, ", "))
		}

		pages := cli.Get(tblOwner, tbl, searchKeys)
		if err != nil {
			PrintErrorAndQuit(err)
		}
		datAll, err := pages.DataAll()

		//
		//valu, err := client.GetData(context.Background(), &zetabase.TableGet{
		//	Id:           identity.Id,
		//	TableId:      tbl,
		//	TableOwnerId: tblOwner,
		//	Nonce:        nonce,
		//	PageIndex:    0,
		//	Credential:   poc,
		//	Keys:         searchKeys,
		//})
		dTyp := viper.GetString(ConfigKeyOutputDataType)
		if err != nil {
			PrintErrorAndQuit(err)
		} else {
			exportMode := viper.GetString(ConfigKeyExportDataMode)

			switch exportMode {
			case "json":
				for _, v := range datAll {
					fmt.Printf("%s\n", string(v))
				}
			case "jsonwithkey":
				for k, v := range datAll {
					fmt.Printf("{\"%s\": %s}\n", k, string(v))
				}
			case "base64":
				for _, v := range datAll {
					s := base64.StdEncoding.EncodeToString(v)
					fmt.Printf("%s\n", s)
				}
			default:
				var pairs []*zbprotocol.DataPair
				var keyLst []string
				for k, v := range datAll {
					pairs = append(pairs, &zbprotocol.DataPair{
						Key:   k,
						Value: v,
					})
					keyLst = append(keyLst, k)
				}
				PrintKeyValuePairs(pairs, dTyp)
			}


		}
	},
}

func listKeys(identity *UserIdentity, tbl, tblOwner, keyPattern string, nonce int64, poc *zbprotocol.ProofOfCredential, client zbprotocol.ZetabaseProviderClient) ([]string, error) {
	Logf("debug - List keys: %s %s %s - pattern %s...", identity.Id, tbl, tblOwner, keyPattern)
	lkr := &zbprotocol.ListKeysRequest{
		Id:           identity.Id,
		TableId:      tbl,
		TableOwnerId: tblOwner,
		Pattern:      keyPattern,
		PageIndex:    0,
		Nonce:        nonce,
		Credential:   poc,
	}
	resp, err := client.ListKeys(context.Background(), lkr)
	if err != nil {
		return nil, err
	} else {
		return resp.GetKeys(), nil
	}
}

func loadIdentityFromConfigs() *UserIdentity {
	identFn := viper.GetString(ConfigKeyIdentityFile)
	identity := parseIdentFromFile(identFn)
	if identity == nil {
		lid := viper.GetString(ConfigKeyLoginId)
		parentId := viper.GetString(ConfigKeyLoginParentId)
		if isVerbose() {
			Logf("Using identity %s (parent: %s)...", lid, parentId)
		}
		if len(viper.GetString(ConfigKeyIdPassword)) > 0 && len(lid) > 0 {
			// JWT login scenario
			var pid *string
			if len(parentId) > 0 {
				pid = &parentId
			}
			return &UserIdentity{
				Id:       "",
				ParentId: pid,
				PubKey:   nil,
				PrivKey:  nil,
			}
		}
		PrintErrorAndQuit(errors.New("FailedToLoadIdentity"))
	}
	if isVerbose() {
		Logf("Using identity %s (from %s)...", identity.Id, identFn)
	}
	return identity
}

var cmdManage = &cobra.Command{
	Use:   "manage",
	Short: "Manage settings and databases",
	Long:  `Manage identities and perform actions on tables.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		task := viper.GetString(ConfigKeyAdminTask)
		conn := dialRemote()
		defer conn.Close()
		client := zbprotocol.NewZetabaseProviderClient(conn)
		if task == AdminTaskNewUser {
			email, name, mobile, adminPass, keyfn, pkfn := collectNewIdDetailsFromConfig()

			pubKeyEnc, e := ioutil.ReadFile(keyfn)
			if e != nil {
				PrintErrorStringAndQuit("Failed to read key file")
			}
			privKeyEnc, e := ioutil.ReadFile(pkfn)
			if e != nil {
				Logf("Cannot open private key file: %s", pkfn)
				PrintErrorStringAndQuit("Failed to read key file")
			}

			res, e := client.RegisterNewIdentity(context.Background(), &zbprotocol.NewIdentityRequest{
				Name:          name,
				Email:         email,
				AdminPassword: adminPass,
				Mobile:        mobile,
				PubKeyEncoded: string(pubKeyEnc),
			})

			if e != nil {
				Logf("Failed to register new identity: %s", e.Error())
			} else if res.GetError() != nil && res.GetError().GetCode() != 0 {
				Logf("Failed to register new identity: %s (%d)", res.GetError().GetMessage(), res.GetError().GetCode())
			} else {
				id := res.GetId()
				idDefn := IdentityDefinition{
					Id:         id,
					ParentId:   "",
					Handle:     name,
					PubKeyEnc:  string(pubKeyEnc),
					PrivKeyEnc: string(privKeyEnc),
				}
				dat, _ := json.MarshalIndent(idDefn, "", " ")
				idFn := fmt.Sprintf("zetabase.%s.identity", cleanStringForFilename(name))
				e := ioutil.WriteFile(idFn, dat, 0644)
				if e != nil {
					PrintErrorAndQuit(e)
				}
				irdr := NewInteractiveReader()
				verifCode := irdr.Prompt("You will receive a text message with a verification code. Please input it now:")
				err := zetabase.ConfirmUserIdentity(context.Background(), res.GetId(), nil, strings.Trim(verifCode, " \r\n"), client)
				for err != nil {
					verifCode = irdr.Prompt("Wrong code. Please re-input it now:")
					err = zetabase.ConfirmUserIdentity(context.Background(), res.GetId(), nil, strings.Trim(verifCode, " \r\n"), client)
				}
				Logf("Success!\n\nResult: saved identity to file: %s\n > ID %s", idFn, id)
			}
		} else if task == AdminTaskNewSubUser {
			parentId := viper.GetString(ConfigKeyParentUid)
			if len(parentId) < 1 {
				Logf("Invalid parent ID: %s", parentId)
				PrintErrorAndQuit(errors.New("InvalidParentId"))
			}
			email, name, mobile, adminPass, keyfn, pkfn := collectNewIdDetailsFromConfig()

			pubKeyEnc, e := ioutil.ReadFile(keyfn)
			if e != nil {
				PrintErrorStringAndQuit("Failed to read key file")
			}
			privKeyEnc, e := ioutil.ReadFile(pkfn)
			if e != nil {
				Logf("Cannot open private key file: %s", pkfn)
				return
			}

			signupCode := viper.GetString(ConfigKeySubIdSignupCode)
			if len(signupCode) == 0 {
				ird := NewInteractiveReader()
				signupCode = ird.Prompt("You didn't provide a signup code. Please enter a signup code if you were provided with one (or leave blank otherwise):")
			}

			groupId := viper.GetString(ConfigKeyGroupId)
			if len(groupId) == 0 {
				ird := NewInteractiveReader()
				groupId = ird.Prompt("You didn't provide a group ID. If you were provided one, please enter it now (or leave blank otherwise):")
				groupId = strings.TrimSpace(groupId)
			}

			res, e := client.CreateUser(context.Background(), &zbprotocol.NewSubIdentityRequest{
				Id:            parentId,
				Name:          name,
				Email:         email,
				Mobile:        mobile,
				LoginPassword: adminPass,
				PubKeyEncoded: string(pubKeyEnc),
				SignupCode:    signupCode,
				GroupId:       groupId,
				//AllowPasswordIdent: true,
			})

			if e != nil {
				PrintErrorAndQuit(e)
				//Logf("Failed to register new identity: %s", e.Error())
			} else if res.GetError() != nil && res.GetError().GetCode() != 0 {
				PrintErrorAndQuit(errors.New(res.GetError().GetMessage()))
				//Logf("Failed to register new identity: %s (%d)", res.GetError().GetMessage(), res.GetError().GetCode())
			} else {
				id := res.GetId()
				idDefn := IdentityDefinition{
					Id:         id,
					ParentId:   parentId,
					Handle:     name,
					PubKeyEnc:  string(pubKeyEnc),
					PrivKeyEnc: string(privKeyEnc),
				}
				dat, _ := json.MarshalIndent(idDefn, "", " ")
				idFn := fmt.Sprintf("zetabase.%s-%s.subidentity", cleanStringForFilename(name), parentId)
				e := ioutil.WriteFile(idFn, dat, 0644)
				if e != nil {
					PrintErrorAndQuit(e)
				}
				irdr := NewInteractiveReader()
				parId := &parentId
				verifCode := irdr.Prompt("You will receive a text message with a verification code. Please input it now:")
				err := zetabase.ConfirmUserIdentity(context.Background(), res.GetId(), parId, strings.Trim(verifCode, " \r\n"), client)
				for err != nil {
					verifCode = irdr.Prompt("Wrong code. Please re-input it now:")
					err = zetabase.ConfirmUserIdentity(context.Background(), res.GetId(), parId, strings.Trim(verifCode, " \r\n"), client)
				}
				Logf("Success!\n\nResult: saved subidentity to file: %s\n > ID %s", idFn, id)
			}

		} else if task == AdminTaskListSubUsers {
			identity := loadIdentityFromConfigs()
			rig := zetabase.NewZetabaseClient(identity.Id)
			provHost := viper.GetString(ConfigKeyZbHostPort)
			if len(provHost) > 0 {
				rig.SetServerAddr(provHost)
			} else {
				rig.SetServerAddr("api.zetabase.io:443")
			}
			if viper.GetBool(ConfigKeyConnectInsecure) {
				rig.SetInsecure()
			}
			rig.SetIdKey(identity.PrivKey, identity.PubKey)
			err := rig.Connect()
			if err != nil {
				PrintErrorAndQuit(err)
			}

			subus, err := rig.GetSubIdentities()
			if err != nil {
				PrintErrorAndQuit(err)
			} else {
				var rows [][]string
				for _, x := range subus {
					row := []string{x.Id, x.Name, x.GroupId}
					rows = append(rows, row)
				}
				PrintSubUsersList(rows)
			}


		} else if task == AdminTaskGetWebhook {
			identity := loadIdentityFromConfigs()
			if len(args) != 1 {
				PrintErrorStringAndQuit("Please specify a table name, e.g. `zb manage -t webhook tablename`.")
			}
			tbl := args[0]
			if len(tbl) == 0 {
				PrintErrorStringAndQuit("Please specify a table name, e.g. `zb manage -t webhook tablename`.")
			}
			tblOwner := chooseDefaultTableOwnerId(identity)
			if tblOwner != identity.Id {
				if isVerbose() {
					Logf("Using table owner: %s", tblOwner)
				}
			}

			log.Printf("Webhook address: https://zetabase.io/api/webhooks/%s/%s\n", tblOwner, tbl)

		} else if task == AdminTaskTestClient {
			// This is the tweets test corresponding to the JS example
			tblId := "tweetstest2"
			parentId := "a68d5254-206c-4782-bb10-eb33037e0d4e"
			client := zetabase.NewZetabaseClient("")
			client.SetInsecure()
			client.SetServerAddr("localhost:9991")
			client.SetParent(parentId)
			client.SetIdPassword("testweb5", "testweb5")
			err := client.Connect()
			if err != nil {
				panic(err)
			}
			keys, err := client.ListKeys("a68d5254-206c-4782-bb10-eb33037e0d4e", tblId).KeysAll()
			log.Printf("Got keys: %v\n", keys)

			if err != nil {
				panic(err)
			} else {
				log.Printf("Connected client with ID: %s\n", client.Id())
			}

			tweet := "hello go world!"
			tweetData := fmt.Sprintf("{\"text\":\"%s\",\"uid\":\"%s\"}", tweet, client.Id())
			key := fmt.Sprintf("tweet/%s/%d", client.Id(), time.Now().Unix())

			err = client.PutData(parentId, tblId, key, []byte(tweetData), false)
			if err != nil {
				panic(err)
			}
		} else if task == AdminTaskTestClient {
			identity := loadIdentityFromConfigs()
			rig := zetabase.NewZetabaseClient(identity.Id)
			rig.SetIdKey(identity.PrivKey, identity.PubKey)
			runLocal := false
			if !runLocal {
				rig.SetServerAddr("api.zetabase.io:443")
			} else {
				rig.SetInsecure()
				rig.SetServerAddr("localhost:9991")
			}
			err := rig.Connect()
			if err != nil {
				PrintErrorAndQuit(err)
			}

			tblId := "perf9"
			var keys []string
			var valus [][]byte


			log.Printf("Connected, generating data...\n")
			for i := 0; i < 100000; i++ {
				k := fmt.Sprintf("key/%d", i)
				s := randLetterStringRunes(512)
				num := 2*rand.Intn(100)
				v := fmt.Sprintf("{\"num\": %d, \"randomstring\":\"%s\"}", num, s)
				keys = append(keys, k)
				valus = append(valus, []byte(v))
			}
			log.Printf("Running `PutMulti`...\n")
			err = rig.PutMulti(identity.Id, tblId, keys, valus, true)
			if err != nil {
				PrintErrorAndQuit(err)
			}

		}  else if task == AdminTaskTestingGround {
			uid := "18259baf-b9e7-4cbd-9027-ca6a4dae1af1"
			handle := "test_user"
			password := "test_pass"

			testClient := zetabase.NewZetabaseClient(uid)
			//testClient.SetInsecure()
			//testClient.SetServerAddr("127.0.0.1:9991")
			testClient.SetServerAddr("api.zetabase.io:443")
			testClient.SetIdPassword(handle, password)
			testClient.SetMaxItemSize(int64(20000))
			err := testClient.Connect()
			if err != nil {
				panic(err)
			}

			append := zbprotocol.PermissionLevel_APPEND
			individual := zbprotocol.PermissionAudienceType_INDIVIDUAL

			perm := zetabase.NewPermissionEntry(append, individual, "1551e62d-41d6-4406-940b-25dccf8d5220")
			permConstraint := zetabase.NewPermConstraintOrder("age")
			perm.AddConstraint(permConstraint)

			testClient.AddPermission(uid, "PermTest1", perm)


			print("sucess")


			// startKeys := time.Now()

			// listKeys := testClient.ListKeys(uid, "testBigDf3")
			// keys, err := listKeys.KeysAll()

			// elapsedKeys := time.Since(startKeys)
			// log.Printf("Keys took %s", elapsedKeys)
			
			// if err != nil {
			// 	panic(err)
			// }

			// startData := time.Now()
		
			// res := testClient.Get(uid, "testBigDf3", keys)
			// data, _ := res.DataAll()
			// print(data[keys[0]])

			// elapsedData := time.Since(startData)
			// log.Printf("Data took %s", elapsedData)			

			// nothing
			// uid := "a68d5254-206c-4782-bb10-eb33037e0d4e"
			// rig := zetabase.NewZetabaseClient(uid)
			// rig.SetServerAddr("localhost:9991")
			// rig.SetInsecure()
			// rig.SetIdPassword("jasonpy1", "jasonpy1")
			// err := rig.Connect()
			// if err != nil {
			// 	panic(err)
			// }
			// ksPgs := rig.ListKeysWithPattern(uid, "simulation6", "data/%")
			// ks, err := ksPgs.KeysAll()
			// if err != nil {
			// 	panic(err)
			// }
			// Logf("Number of keys: %d", len(ks))
			// m := map[string]bool{}
			// for _, k := range ks {
			// 	if _, ok := m[k]; ok {
			// 		Logf("Duplicate! %s", k)
			// 	}
			// 	m[k] = true
			// }

		} else {
			Logf("No such task `%s`", task)
		}
	},
}

func isVerbose() bool {
	return viper.GetBool(ConfigKeyVerbose)
}


func makeNewClient(uid string, privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) *zetabase.ZetabaseClient {
	insec := viper.GetBool(ConfigKeyConnectInsecure)
	host := viper.GetString(ConfigKeyZbHostPort)
	loginParentId := viper.GetString(ConfigKeyLoginParentId)
	loginHandl := viper.GetString(ConfigKeyLoginId)
	loginPass := viper.GetString(ConfigKeyIdPassword)
	cli := zetabase.NewZetabaseClient(uid)
	if len(loginParentId) > 0 {
		cli.SetParent(loginParentId)
	}
	if insec {
		cli.SetInsecure()
	}
	cli.SetServerAddr(host)
	if privKey != nil && pubKey != nil {
		cli.SetIdKey(privKey, pubKey)
	} else if len(loginHandl) > 0 {
		cli.SetIdPassword(loginHandl, loginPass)
	}
	err := cli.Connect()
	if err != nil {
		return nil
	} else {
		return cli
	}
}

func dialRemote() *grpc.ClientConn {
	insec := viper.GetBool(ConfigKeyConnectInsecure)
	host := viper.GetString(ConfigKeyZbHostPort)
	noCertVerif := viper.GetBool(ConfigKeyConnectNoCertVerify)
	if insec {
		wi := grpc.WithInsecure()
		if len(host) == 0 {
			// TODO: for debug purposes only
			host = "localhost:9991"
		}
		conn, err := grpc.Dial(host, wi)
		if err != nil {
			PrintErrorAndQuit(err)
		}
		return conn
	} else if noCertVerif {
		//Logf("debug: %s", "Not doing insec skip verif...")
		//configT := &tls.Config{InsecureSkipVerify: true}
		configT := &tls.Config{InsecureSkipVerify: false}
		wtc := grpc.WithTransportCredentials(credentials.NewTLS(configT))
		conn, err := grpc.Dial(host, wtc)
		if err != nil {
			Logf("Connection error: %s", err.Error())
			return nil
		}
		return conn
	} else {
		//wi := grpc.WithInsecure()
		tlsCert, x509, err := getcert.FromTLSServer(host, true)

		if viper.GetBool(ConfigKeyVerbose) && zetabase.DoSaveCertificates {
			Logf("certexport: EXPORTING CERTIFICATES...")
			fnPre := "zbcert."
			for i, c := range x509 {
				fn := fmt.Sprintf("%s%d", fnPre, i)
				bs := convertCertificate(c, c.PublicKey)
				ioutil.WriteFile(fn, bs, 0644)
			}
		}

		//cr := credentials.NewClientTLSFromCert(tlsCert.Leaf, "")
		cr := credentials.NewServerTLSFromCert(&tlsCert)
		wtc := grpc.WithTransportCredentials(cr)
		if len(host) == 0 {
			// TODO: for debug purposes only
			host = "localhost:9991"
		}
		//conn, err := grpc.Dial(host, wi)
		conn, err := grpc.Dial(host, wtc)
		if err != nil {
			PrintErrorAndQuit(err)
		}
		return conn
	}
}

func collectNewIdDetailsFromConfig() (string, string, string, string, string, string) {
	email := viper.GetString(ConfigKeyUserEmail)
	name := viper.GetString(ConfigKeyUserName)
	mobile := viper.GetString(ConfigKeyUserMobile)
	admPass := viper.GetString(ConfigKeyAdminPassword)
	keyfn := viper.GetString(ConfigKeyPubKeyFile)
	pkfn := viper.GetString(ConfigKeyPrivKeyFile)
	reader := NewInteractiveReader()


	x := reader.Prompt("Please review the Zetabase terms of service available online at https://zetabase.io/tos\n" +
		" Do you accept and agree to the terms and conditions established in the Agreement?")
	if strings.ToLower(strings.TrimSpace(x))[0] != 'y' {
		os.Exit(2)
	}

	if len(email) == 0 {
		x := reader.Prompt("Email address:")
		email = strings.TrimSpace(x)
	}
	if len(name) == 0 {
		x := reader.Prompt("Enter a user handle (e.g. your name):")
		name = strings.TrimSpace(x)
	}
	if len(mobile) == 0 {
		x := reader.Prompt("Your mobile number with region code (e.g. +12125551212 for U.S.):")
		for !validatePhoneNumber(x) {
			x = reader.Prompt("Invalid number. A phone number is required to confirm your account and set up 2FA.\nPlease provide your mobile number in international format (e.g. +12125551212 for U.S.):")
		}
		mobile = strings.TrimSpace(x)
	}
	if len(admPass) == 0 {
		x := reader.Prompt("Your administrator website password:")
		for len(x) < 6 {
			x = reader.Prompt("Invalid password (too short). Your administrator website password:")
		}
		admPass = strings.TrimSpace(x)
	}
	if len(keyfn) == 0 || len(pkfn) == 0 {
		x := reader.Prompt("Would you like us to generate a key for this identity?")
		if strings.ToLower(x)[0] == 'y' {
			privkey, pubkey := zetabase.GenerateKeyPair()
			ts := time.Now().Unix()
			fnPub := fmt.Sprintf("zetabase.%d.pub", ts)
			fnPriv := fmt.Sprintf("zetabase.%d.priv", ts)
			keyfn = fnPub
			pkfn = fnPriv
			bs, err := zetabase.EncodeEcdsaPublicKey(pubkey)
			if err != nil {
				PrintErrorAndQuit(err)
			}
			err = ioutil.WriteFile(fnPub, bs, 0644)
			if err != nil {
				PrintErrorAndQuit(err)
			}
			bs, err = zetabase.EncodeEcdsaPrivateKey(privkey)
			if err != nil {
				PrintErrorAndQuit(err)
			}
			err = ioutil.WriteFile(fnPriv, bs, 0644)
			if err != nil {
				PrintErrorAndQuit(err)
			}
		} else {
			x := reader.Prompt("Path to public key file:")
			keyfn = strings.TrimSpace(x)
			x = reader.Prompt("Path to private key file:")
			pkfn = strings.TrimSpace(x)
		}
	}
	return email, name, mobile, admPass, keyfn, pkfn
}
