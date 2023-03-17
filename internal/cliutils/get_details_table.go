package cliutils

import (
	"bytes"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func GetDetailsTable(keyValueMap map[string]string, disableSort ...bool) string {
	keys := []string{}
	for key, _ := range keyValueMap {
		keys = append(keys, key)
	}
	if len(disableSort) == 0 || !disableSort[0] {
		sort.Strings(keys)
	}
	var tableData bytes.Buffer
	table := tablewriter.NewWriter(&tableData)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"property", "value"})
	for _, key := range keys {
		table.Append([]string{key, keyValueMap[key]})
	}
	table.Render()
	return tableData.String()
}
