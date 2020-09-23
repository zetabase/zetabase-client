package main

import (
	"encoding/hex"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"os"
	"strings"
	"unicode/utf8"
)

func PrintSingleColumn(title string, valu []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", title})
	var rows []table.Row
	for i, x := range valu {
		rows = append(rows, table.Row{i+1, x})
	}
	t.AppendRows(rows)
	t.SetStyle(getStyleForOs())
	t.Render()
}

func stringifyPermEntry(p *zbprotocol.PermissionsEntry) string {
	at := p.AudienceType.String()
	lvl := p.Level.String()
	aid := p.AudienceId
	var cstrs []string
	for _, c := range p.Constraints {
		cs := ""
		if c.GetFieldConstraint() != nil {
			reqVal := c.GetFieldConstraint().GetRequiredValue()
			switch c.GetFieldConstraint().GetValueType() {
			case zbprotocol.FieldConstraintValueType_UID:
				reqVal = "@uid"
			case zbprotocol.FieldConstraintValueType_TIMESTAMP:
				reqVal = "@time"
			case zbprotocol.FieldConstraintValueType_NATURAL_ORDER:
				reqVal = "@order"
			case zbprotocol.FieldConstraintValueType_RANDOM:
				reqVal = "@random"
			case zbprotocol.FieldConstraintValueType_CONSTANT:
				reqVal = reqVal // no change
			}
			opDesc := " = "
			cs = c.GetFieldConstraint().GetFieldKey() + opDesc + reqVal
		} else {
			reqVal := ""
			switch c.GetKeyConstraint().GetValueType() {
			case zbprotocol.FieldConstraintValueType_UID:
				reqVal = "@uid"
			case zbprotocol.FieldConstraintValueType_TIMESTAMP:
				reqVal = "@time"
			case zbprotocol.FieldConstraintValueType_NATURAL_ORDER:
				reqVal = "@order"
			case zbprotocol.FieldConstraintValueType_RANDOM:
				reqVal = "@random"
			case zbprotocol.FieldConstraintValueType_CONSTANT:
				reqVal = reqVal // no change
			}
			cs = "@key = " + c.GetKeyConstraint().GetRequiredPrefix() + reqVal + c.GetKeyConstraint().GetRequiredSuffix()
		}
		cstrs = append(cstrs, cs)
	}
	constStr := ""
	if len(cstrs) > 0 {
		constStr = "(" + strings.Join(cstrs, ", ") + ")"
	}
	if len(constStr) > 0 {
		return fmt.Sprintf("%s %s %s %s", at, lvl, aid, constStr)
	} else {
		return fmt.Sprintf("%s %s %s", at, lvl, aid)
	}
}

func stringifyPermissions(perms []*zbprotocol.PermissionsEntry) string {
	var ps []string
	for _, x := range perms {
		ps = append(ps, stringifyPermEntry(x))
	}
	return strings.Join(ps, ", ")
}

func PrintTableDefinitions(defns []*zbprotocol.TableCreate) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Table Id", "Indexed Fields", "Data Format", "Permissions"})
	var rows []table.Row
	for i, x := range defns {
		var df string
		switch x.DataFormat {
		case zbprotocol.TableDataFormat_PLAIN_TEXT:
			df = "text"
		case zbprotocol.TableDataFormat_JSON:
			df = "json"
		case zbprotocol.TableDataFormat_BINARY:
			df = "binary"
		default:
			df = "other"
		}
		var idxs []string
		if x.GetIndices() != nil {
			for _, idf := range x.GetIndices().GetFields() {
				idxs = append(idxs, idf.Field)
			}
		}
		idxStr := strings.Join(idxs, ", ")
		rows = append(rows, table.Row{i+1, x.TableId, idxStr, df, stringifyPermissions(x.GetPermissions())})
	}
	t.AppendRows(rows)
	t.SetStyle(getStyleForOs())
	t.Render()
}

func areDataPairsAllStringValues(data []*zbprotocol.DataPair) bool {
	success := true
	for _, d := range data {
		if !utf8.Valid(d.GetValue()) {
			success = false
			break
		}
	}
	return success
}

func getStyleForOs() table.Style {
	if isWindows() {
		return table.StyleDefault
	} else {
		return table.StyleRounded
	}
}

func PrintSubUsersList(data [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"User ID", "Handle", "User Group"})
	var rows []table.Row
	for i := 0; i < len(data); i++ {
		rows = append(rows, table.Row{data[i][0], data[i][1], data[i][2]})
	}
	t.AppendRows(rows)
	t.SetStyle(getStyleForOs())
	t.Render()
}

func PrintShellCommandsTable(cmds, usages, descs []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Command", "Usage", "Description"})
	var rows []table.Row
	for i := 0; i < len(cmds); i++ {
		rows = append(rows, table.Row{cmds[i], usages[i], descs[i]})
	}
	t.AppendRows(rows)
	t.SetStyle(getStyleForOs())
	t.Render()
}

func PrintKeyValuePairs(data []*zbprotocol.DataPair, dTyp string) {
	if dTyp == "" {
		if areDataPairsAllStringValues(data) {
			dTyp = "text"
			if isVerbose() {
				Logf("\t NOTE:  Printing data in text mode. Use `-X binary` to print hex.")
			}
		} else {
			if dTyp != "json" && dTyp != "text" && isVerbose() {
				Logf("\t NOTE:  Printing data in binary mode (hex). Use parameters `-X json` or `-X text` to render as text.")
			}
		}
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Key", "Value"})
	var rows []table.Row
	for d, x := range data {
		var suffix string
		dataToRender := x.GetValue()
		if len(dataToRender) > 64 {
			dataToRender = dataToRender[:64]
			suffix = "..."
		}
		var strEnc string
		switch dTyp {
		case "json":
			strEnc = string(dataToRender)
		case "text":
			strEnc = string(dataToRender)
		default:
			strEnc = "0x" + hex.EncodeToString(dataToRender)
		}
		k, s := x.Key, strEnc
		rows = append(rows, table.Row{d+1, k, s+suffix})
		//Logf(" [%d]  %s -> %s", d, k, s)
	}
	t.AppendRows(rows)
	t.SetStyle(getStyleForOs())
	//t.Style().Color.Header = text.Colors{text.BgBlack, text.FgWhite}
	//t.Style().Format.Footer = text.FormatLower
	//t.Style().Options.DrawBorder = false
	t.Render()
}
