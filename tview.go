package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func tviewFillTable(table *tview.Table, columns []string, data [][]string) {
	for i := 1; i < len(columns); i++ {
		cell := tview.NewTableCell("[yellow]" + columns[i]).SetBackgroundColor(tcell.ColorBlue)
		cell.SetSelectable(false)
		table.SetCell(0, i-1, cell)
		for j := 1; j <= len(data); j++ {
			content := data[j-1][i]
			cell := tview.NewTableCell(content)
			cell.SetMaxWidth(32)
			table.SetCell(j, i-1, cell)
		}
	}
}

func myTview(columns []string, data [][]string) {
	app := tview.NewApplication()
	table := tview.NewTable()
	// table.SetBorder(true)
	table.SetTitle(" LDAP ")
	table.SetFixedColumnsWidth(true)
	// table.SetBorders(true)
	table.SetSeparator(tview.Borders.Vertical)
	table.SetFixed(1, 0)
	table.SetSelectable(true, false)
	tviewFillTable(table, columns, data)
	table.SetDoneFunc(func(key tcell.Key) {
		app.Stop()
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.Stop()
				return nil
			case 'e':
				app.Suspend(func() {
					row, _ := table.GetSelection()
					dn := data[row-1][0]
					cmd := exec.Command("ldapvi", "-s", "base", "-b", dn)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						log.Printf("ldapvi: " + err.Error())
						time.Sleep(5 * time.Second)
					}
					columns, data := ldapSearch(LdapDN, LdapFilter, LdapAttrs)
					tviewFillTable(table, columns, data)
				})
			case '/':
				row, _ := table.GetSelection()
				app.Suspend(func() {
					fmt.Printf("search: current row=%d\n", row)
					time.Sleep(time.Second)
				})
			}
		}
		return event
	})
	text := tview.NewTextView()
	text.SetBackgroundColor(tcell.ColorBlue)
	text.SetDynamicColors(true)
	text.SetText(" [yellow]q:quit   e:edit   f:filter   s:sort   /:search   n:next")
	flex := tview.NewFlex()
	flex.SetBackgroundColor(tcell.ColorRed)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(text, 1, 0, false)
	flex.AddItem(tview.NewBox(), 1, 0, false)
	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
