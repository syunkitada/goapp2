package str_utils

import (
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetTable() (tableString *strings.Builder, table *tablewriter.Table) {
	tableString = &strings.Builder{}
	table = tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	return
}
