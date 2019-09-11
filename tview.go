package main

import (
	// "github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func myTview(columns []string, data [][]string) {
	app := tview.NewApplication()
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle(" LDAP ")
	// table.SetBorders(true)
	table.SetSeparator(tview.Borders.Vertical)
	table.SetFixed(1, 0)
	table.SetSelectable(true, false)
	for i := 0; i < len(columns); i++ {
		cell := tview.NewTableCell("[yellow]" + columns[i])
		cell.SetSelectable(false)
		table.SetCell(0, i, cell)
		for j := 1; j <= len(data); j++ {
			content := data[j-1][i]
			if runes := []rune(content); len(runes) > 20 {
				content = string(runes[:20]) + "[green]…"
			}
			cell := tview.NewTableCell(content)
			table.SetCell(j, i, cell)
		}
	}
	app.SetRoot(table, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
