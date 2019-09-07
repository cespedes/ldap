package main

import (
	// "github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func my_tview(columns []string, data [][]string) {
	app := tview.NewApplication()
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	box = box
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle("title")
	// table.SetBorders(true)
	table.SetSeparator(tview.Borders.Vertical)
	table.SetFixed(1, 0)
	table.SetSelectable(true, false)
	for i:=0; i < len(columns); i++ {
		cell := tview.NewTableCell("[yellow]" + columns[i])
		cell.SetSelectable(false)
		table.SetCell(0, i, cell)
		for j:=1; j <= len(data); j++ {
			cell := tview.NewTableCell(data[j-1][i])
			table.SetCell(j, i, cell)
		}
	}
	app.SetRoot(table, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

