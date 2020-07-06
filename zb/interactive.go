package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"github.com/zetabase/zetabase-client"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	CliPrompt = "> "
)

var (
	ErrorExplanations = map[string]string{
		"NoSuchSymbol":            "The table or object you are searching for does not exist by that symbol.",
		"NoSuchUser":              "The user cannot be properly identified. Please check your identity (and table owner ID if applicable.)",
		"InvalidSignature":        "The signature provided was invalid. Please verify your identity settings and make sure your key files are in the proper locations.",
		"InsufficientCredentials": "The identity you are using is not allowed to access the resource. Check that your identity is properly loaded and all needed permissions have been added.",
		"BadSignupCode":           "You omitted or entered an invalid signup code for establishing a new user identity. Please obtain an up-to-date signup code.",
		"BetaRestriction":         "This action is restricted based on the rulse of the Beta Program. Please go to zetabase.io to switch account types.",
	}
)

type InteractiveReader struct {
	reader *bufio.Reader
}

func NewInteractiveReader() *InteractiveReader {
	reader := bufio.NewReader(os.Stdin)
	return &InteractiveReader{reader}
}

func (r *InteractiveReader) Prompt(question string) string {
	fmt.Print(question + "\n" + CliPrompt)
	text, _ := r.reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return strings.TrimSpace(text)
	return text
}

func (r *InteractiveReader) PromptInline(prompt string) string {
	fmt.Print(prompt)
	text, _ := r.reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return strings.TrimSpace(text)
	return text
}

func PrintErrorStringAndQuit(err string) {
	exp := "Other error: " + err
	Logf("Error:\n\t%s\nTerminating...", exp)
	os.Exit(1)
}

func symbolFromGrpcError(err error) string {
	arr := strings.Split(err.Error(), " ")
	return arr[len(arr)-1]
}

func PrintErrorAndQuit(err error) {
	exp := "Other error: " + err.Error()
	if x, ok := ErrorExplanations[symbolFromGrpcError(err)]; ok {
		exp = x
	}
	Logf("Error:\n\t%s\nTerminating...", exp)
	os.Exit(1)
}

func shellParse(usrCmd string) []string {
	var res []string
	var curTok string
	inQuot := false
	for i := 0; i < len(usrCmd); i++ {
		curChar := usrCmd[i]
		if curChar == '"' && (i == 0 || usrCmd[i-1] != '\\') { // Check for escape
			inQuot = !inQuot
		} else if (!inQuot) && (curChar == ' ') {
			if len(curTok) > 0 {
				res = append(res, curTok)
				curTok = ""
			}
		} else {
			if curChar == '\\' && (i == 0 || usrCmd[i-1] != '\\') {
				// nothing
			} else if curChar == '\\' && (i == 0 || usrCmd[i-1] == '\\') {
				curTok = curTok + string(curChar)
			} else {
				curTok = curTok + string(curChar)
			}
		}
	}
	if len(curTok) > 0 {
		res = append(res, curTok)
	}
	return res
}

func printShellHelp() {
	cmds := []string{"ls", "get", "put", "set", "create", "rm", "query"}
	usages := []string{"ls (tables|keys) [table name] [keypattern/%]", "get <table name> <pattern>",
		"put <table name> <key> <value>", "set <table name> <key> <value>",
		"create <table name> <data type> <permissions> [field1 order1 field2 order2...]", "rm (table|key) [table-name] [key]",
		"query <table-name> (<query>)"}
	descs := []string{"List tables or keys, where % denotes the suffix wildcard (e.g. `key%`)",
		"View keys/values for a given pattern where `%` is the suffix wildcard (e.g. `view mydata user/%`)",
		"Add data to a given table. Value will be UTF-8 encoded.",
		"Set a key in a given table (and overwrite if exists).",
		"Create a new table with the given permissions and field indexing.",
		"Delete a single object by key or drop a full table.",
		"Query <table-name> based on indexed fields. Syntax: e.g. `query mytable (fld1 = v1 and (fld2 = v2 or fld3 = v3))`"}
	PrintShellCommandsTable(cmds, usages, descs)
}

func printShellError(err error) {
	fmt.Printf("\tError:  %s\n", err.Error())
}

var (
	tblOwnerId = ""
)

func shellHandle(promptSuff, quitToken string, identity *UserIdentity, zbclient *zetabase.ZetabaseClient, s string, nonceMaker *zetabase.NonceMaker, tblOwnerId0 string, shellState *ShellState) {
	tblOwnerId = tblOwnerId0
	client := zbclient.GrpcClient()
	args := shellParse(s)
	if len(args) < 1 {
		return
	}
	nonce := nonceMaker.Get()
	loginUid  := identity.Id

	getCredential := func(nonce int64, xBytes []byte) (*zbprotocol.ProofOfCredential, error) {
		uid, poc, err := getUserCredential(identity, nonce, xBytes, client)
		if err != nil {
			//Logf("debug - error returned: %s...", err.Error())
			return nil, err
		}
		if len(uid) > 0 && (len(identity.Id) == 0 || len(tblOwnerId) == 0  || len(loginUid) == 0) {
			loginUid = uid
			identity.Id = loginUid
			tblOwnerId = identity.Id
			//Logf("debug - tbl owner ID set to %s...", tblOwnerId)
		} else {
			//Logf("debug - no UID returned...")
		}
		return poc, nil
	}

	switch args[0] {
	case "query":
		if len(args) < 3 {
			printShellHelp()
			return
		}
		tblId := args[1]
		qryStart := strings.Index(s, "(")
		qryTxt := s[qryStart:]
		parser := zetabase.NewBQParser(qryTxt)
		res, err := parser.Parse()
		if err != nil {
			printShellError(err)
			return
		}
		qry := res.ToQuery()
		pages := zbclient.Query(identity.Id, tblId, qry)
		var pairs []*zbprotocol.DataPair
		var keyLst []string
		dat, _ := pages.Data()
		for k, v := range dat {
			pairs = append(pairs, &zbprotocol.DataPair{
				Key:   k,
				Value: v,
			})
			keyLst = append(keyLst, k)
		}
		dTyp := viper.GetString(ConfigKeyOutputDataType)
		shellState.SetKeyBuffer(tblId, keyLst)
		shellState.AddTableIdToHistory(tblId)
		PrintKeyValuePairs(pairs, dTyp)
	case "rm":
		if len(args) < 3 {
			printShellHelp()
			return
		}
		switch args[1] {
		case "key":
			if len(args) >= 4 {
				tblId, key := args[2], args[3]
				//_, poc, err := getUserCredential(identity, nonce, []byte(key), client)
				poc, err := getCredential(nonce, []byte(key))
				if err != nil {
					printShellError(err)
				} else {
					err := deleteKey(identity, identity.Id, tblId, key, nonce, poc, client)
					if err != nil {
						printShellError(err)
					} else {
						Logf("Success.")
					}
				}
			} else {
				printShellHelp()
			}
		case "table":
			tblId := args[2]
			//_, poc, err := getUserCredential(identity, nonce, []byte(tblId), client)
			poc, err := getCredential(nonce, []byte(tblId))
			if err != nil {
				printShellError(err)
			} else {
				err = deleteTable(identity, identity.Id, tblId, nonce, poc, client)
				if err != nil {
					printShellError(err)
				} else {
					Logf("Success.")
				}
			}
		default:
			printShellHelp()
		}
	case "ls", "list":
		if len(args) >= 2 && args[1] == "tables" {
			poc, err := getCredential(nonce, nil)
			//uid, poc, err := getUserCredential(identity, nonce, nil, client)
			//if len(uid) > 0 && len(identity.Id) == 0 {
			//	loginUid = uid
			//	identity.Id = loginUid
			//	tblOwnerId = identity.Id
			//}
			if err != nil {
				printShellError(err)
			} else {
				ts, err := listTables(identity, tblOwnerId, nonce, poc, client)
				if err != nil {
					printShellError(err)
				} else {
					shellState.IngestTablesData(ts)
					PrintTableDefinitions(ts)
				}
			}
		} else if len(args) >= 3 && args[1] == "keys" {
			keyPat := "%"
			if len(args) >= 4 {
				Logf("With pattern: %s", args[3])
				keyPat = args[3]
			}
			//_, poc, err := getUserCredential(identity, nonce, nil, client)
			poc, err := getCredential(nonce, nil)
			if err != nil {
				printShellError(err)
			} else {
				//Logf("debug - tbl owner = %s", tblOwnerId)
				keyLst, err := listKeys(identity, args[2], tblOwnerId, keyPat, nonce, poc, client)
				//keyHistBuf = &keyLst
				//tblHistBuf = &args[2]
				if err != nil {
					printShellError(err)
				} else {
					shellState.SetKeyBuffer(args[2], keyLst)
					shellState.AddTableIdToHistory(args[2])
					PrintSingleColumn("Key", keyLst)
				}
			}
		} else {
			// print error
			Logf("Bad syntax.")
			printShellHelp()
		}
	case "help":
		printShellHelp()
	case "put", "set":
		doOvr := (args[0] == "set")
		if len(args) != 4 {
			Logf("Invalid usage.")
			printShellHelp()
		} else {
			tbl, key, valu := args[1], args[2], args[3]
			nonce := nonceMaker.Get()
			xBytes := zetabase.TablePutExtraSigningBytes(key, []byte(valu))
			//_, poc, err := getUserCredential(identity, nonce, xBytes, client)
			poc, err := getCredential(nonce, xBytes)
			if isVerbose() {
				//Logf("Signing bytes for put: %x", xBytes)
			}
			if err != nil {
				printShellError(err)
			} else {
				resp, err := client.PutData(context.Background(), &zbprotocol.TablePut{
					Id:           identity.Id,
					TableOwnerId: tblOwnerId,
					TableId:      tbl,
					Key:          key,
					Value:        []byte(valu),
					Overwrite:    doOvr,
					Nonce:        nonce,
					Credential:   poc,
				})
				if err != nil {
					printShellError(err)
				} else {
					if resp == nil || resp.GetCode() == 0 {
						shellState.AddTableIdToHistory(tbl)
						Logf("Success.")
					} else {
						printShellError(errors.New(resp.GetMessage()))
					}
				}
			}
		}
	case "create":
		if len(args) < 3 {
			Logf("Invalid usage.")
			printShellHelp()
		} else {
			var idxArgs []string
			tblid, typ, perms := args[1], args[2], ""
			if len(args) > 4 {
				perms = args[3]
				idxArgs = args[4:]
			} else if len(args) == 4 {
				perms = args[3]
			}
			//_, poc, err := getUserCredential(identity, nonce, []byte(tblid), client)
			poc, err := getCredential(nonce, []byte(tblid))
			if err != nil {
				printShellError(err)
			} else {
				//rig, err := createTable(perms, identity, tblid, typ, idxArgs, nonce, poc, client)

				ctr := getTableCreate(perms, identity, tblid, typ, idxArgs, nonce, poc)
				sigBytes := zetabase.TableCreateSigningBytes(ctr.GetTableId(), ctr.GetPermissions())
				//_, poc, err = getUserCredential(identity, nonce, sigBytes, client)
				poc, err = getCredential(nonce, sigBytes)
				if err != nil {
					printShellError(err)
					return
				}
				ctr.Credential = poc
				rig, err := createTable(ctr, client)

				if err != nil {
					printShellError(err)
				} else {
					if rig == nil || rig.GetCode() == 0 {
						shellState.AddTableIdToHistory(tblid)
						Logf("Success.")
					} else {
						printShellError(errors.New(rig.GetMessage()))
					}
				}
			}

		}
	case "get":
		var searchKeys []string
		var tblId string
		if len(args) == 1 {
			// Use previous keys
			keyHistBuf := shellState.GetKeyBuffer()
			if keyHistBuf != nil {
				Logf("Using keys from search buffer: %s", strings.Join(keyHistBuf, ", "))
				searchKeys = keyHistBuf
				tblId = shellState.GetKeyBufferTable()
			} else {
				Logf("Invalid usage. The keyword `get` without arguments will fetch the results from the last `ls keys...` command.")
				printShellHelp()
				return
			}
		} else if len(args) == 2 {
			Logf("Invalid usage.")
			printShellHelp()
			return
		}
		if len(args) == 3 || len(searchKeys) > 0 {
			nonce := nonceMaker.Get()
			//_, poc, err := getUserCredential(identity, nonce, nil, client)
			poc, err := getCredential(nonce, nil)
			if len(searchKeys) == 0 {
				keyPat := args[2]

				if !strings.HasSuffix(keyPat, "%") {
					searchKeys = []string{keyPat}
				} else {
					ks, err := listKeys(identity, args[1], tblOwnerId, keyPat, nonce, poc, client)
					if err != nil {
						Logf("Key lookup error: check syntax and symbols.")
						printShellError(err)
						return
					} else {
						searchKeys = ks
						//Logf("Got relevant keys: %v...", searchKeys)
						nonce = nonceMaker.Get()
						// Refresh user POC
						//_, poc, err = getUserCredential(identity, nonce, nil, client)
						poc, err = getCredential(nonce, nil)
						if err != nil {
							printShellError(err)
							return
						} else {
							//Logf("Refreshed credential: %v", poc)
						}
					}
				}
				if len(args) >= 2 {
					tblId = args[1]
					shellState.AddTableIdToHistory(args[1])
				}
			}

			if err != nil {
				printShellError(err)
			} else {
				//Logf("Getting data from %s for %s (table %s)", tblOwnerId, identity.Id, tblId)
				data, err := client.GetData(context.Background(), &zbprotocol.TableGet{
					Id:           identity.Id,
					TableOwnerId: tblOwnerId,
					TableId:      tblId,
					Nonce:        nonce,
					Credential:   poc,
					Keys:         searchKeys,
				})
				if err != nil {
					printShellError(err)
				} else {
					if len(args) >= 2 {
						shellState.AddTableIdToHistory(args[1])
					}
					dTyp := viper.GetString(ConfigKeyOutputDataType)
					PrintKeyValuePairs(data.GetData(), dTyp)
				}
			}
		} else {
			Logf("Invalid usage.")
			printShellHelp()
		}
	case quitToken:
		os.Exit(0)
	default:
		fmt.Printf("Unknown command `%s`. Enter `help` for a list of available commands (`quit` to exit).\n", args[0])
	}
}

func makeSuggestionsFrom(tableIds []string) []prompt.Suggest {
	var suggs []prompt.Suggest
	for _, x := range tableIds {
		suggs = append(suggs, prompt.Suggest{Text: x})
	}
	return suggs
}

func makeCompleter(state *ShellState) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		dText := strings.ToLower(d.Text)
		if strings.HasPrefix(dText, "ls") {
			if strings.HasPrefix(dText, "ls keys") {
				suggs := makeSuggestionsFrom(state.GetUsedTableIds())
				//Logf("Suggs = %v", suggs)
				//return suggs
				return prompt.FilterHasPrefix(suggs, d.GetWordBeforeCursor(), true)
			}
			s := []prompt.Suggest{
				{Text: "tables", Description: "List tables"},
				{Text: "keys", Description: "List keys by pattern"},
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		} else if strings.HasPrefix(dText, "rm") {
			if strings.HasPrefix(dText, "rm key") || strings.HasPrefix(dText, "rm table") {
				suggs := makeSuggestionsFrom(state.GetUsedTableIds())
				return prompt.FilterHasPrefix(suggs, d.GetWordBeforeCursor(), true)
			}
			s := []prompt.Suggest{
				{Text: "table", Description: "Delete a table"},
				{Text: "key", Description: "Delete a table object by key"},
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		} else if (strings.HasPrefix(dText, "query ")) && len(strings.Split(dText, " ")) >= 2 {
			suggs := makeSuggestionsFrom(state.GetUsedTableIds())
			//return suggs
			return prompt.FilterHasPrefix(suggs, d.GetWordBeforeCursor(), true)
		} else if (strings.HasPrefix(dText, "get ") || strings.HasPrefix(dText, "set ") || strings.HasPrefix(dText, "put ")) && len(strings.Split(dText, " ")) >= 2 {
			suggs := makeSuggestionsFrom(state.GetUsedTableIds())
			//return suggs
			if len(strings.Split(dText, " ")) >= 3 || strings.Contains(dText, "(") {
				return nil
			}
			return prompt.FilterHasPrefix(suggs, d.GetWordBeforeCursor(), true)
		} else if len(strings.Split(dText, " ")) < 2 {
			s := []prompt.Suggest{
				{Text: "ls", Description: "List keys and tables"},
				{Text: "get", Description: "View data (by key or key pattern)"},
				{Text: "put", Description: "Add data to a table (append-only; will not overwrite)"},
				{Text: "create", Description: "Create a new table with custom structure, permissions, and field constraints"},
				{Text: "query", Description: "Query table based on indexed fields"},
				{Text: "rm", Description: "Delete a single key or entire table"},
				{Text: "set", Description: "Add data to a table (with overwrite)"},
				{Text: "help", Description: "Show available commands"},
				{Text: "quit", Description: "Exit shell"},
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}
		return nil
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	if strings.HasPrefix(d.Text, "ls") {
		if strings.HasPrefix(d.Text, "ls keys") {
			return nil
		}
		s := []prompt.Suggest{
			{Text: "tables", Description: "List tables"},
			{Text: "keys", Description: "List keys by pattern"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	} else if len(strings.Split(d.Text, " ")) < 2 {
		s := []prompt.Suggest{
			{Text: "ls", Description: "List keys and tables"},
			{Text: "get", Description: "View data (by key or key pattern)"},
			{Text: "put", Description: "Add data to a table"},
			{Text: "set", Description: "Add data to a table (with overwrite)"},
			{Text: "create", Description: "Create a new table with custom structure, permissions, and field constraints"},
			{Text: "help", Description: "Show available commands"},
			{Text: "quit", Description: "Exit shell"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return nil
}

func WindowsShellLoop(promptStuff, quitToken string, identity *UserIdentity, zbclient *zetabase.ZetabaseClient) error {
	nonceMaker := zetabase.NewNonceMaker()
	idPref := identity.Id
	if len(idPref) > 3 {
		idPref = idPref[len(idPref)-3:]
	}
	promptStr := idPref + "." + promptStuff
	tblOwnerId := identity.Id
	stateHolder := NewShellState()
	ir := NewInteractiveReader()
	x := ir.PromptInline(promptStr)
	for {
		shellHandle(promptStuff, quitToken, identity, zbclient, x, nonceMaker, tblOwnerId, stateHolder)
		x = ir.PromptInline(promptStr)
	}
}

func ShellLoop(promptStuff, quitToken string, identity *UserIdentity, zbclient *zetabase.ZetabaseClient) error {
	nonceMaker := zetabase.NewNonceMaker()
	idPref := identity.Id
	if len(idPref) > 3 {
		idPref = idPref[len(idPref)-3:]
	}
	promptStr := idPref + "." + promptStuff
	tblOwnerId := identity.Id
	stateHolder := NewShellState()
	//p := prompt.New(func(s string) {
	//	shellHandle(promptStuff, quitToken, identity, client, s, nonceMaker, tblOwnerId, stateHolder)
	//}, completer, prompt.OptionPrefix(promptStr), )
	p := prompt.New(func(s string) {
		shellHandle(promptStuff, quitToken, identity, zbclient, s, nonceMaker, tblOwnerId, stateHolder)
	}, makeCompleter(stateHolder), prompt.OptionPrefix(promptStr), )
	p.Run()
	return nil
}
