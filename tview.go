package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func tviewFillTable(table *tview.Table, columns []string, data [][]string) {
	for i := 0; i < len(columns); i++ {
		cell := tview.NewTableCell("[yellow]" + columns[i]).SetBackgroundColor(tcell.ColorBlue)
		cell.SetSelectable(false)
		table.SetCell(0, i, cell)
		for j := 0; j < len(data); j++ {
			content := data[j][i]
			cell := tview.NewTableCell(content)
			cell.SetMaxWidth(32)
			table.SetCell(j+1, i, cell)
		}
	}
}

func myTview(rows []string, columns []string, data [][]string) {
	app := tview.NewApplication()
	table := tview.NewTable()
	text := tview.NewTextView()
	flex := tview.NewFlex()
	var lastLine tview.Primitive
	var lastSearch string

	tviewSearch := func(row int, text string) bool {
		text = strings.ToLower(text)
		for i := 0; i < len(data); i++ {
			for j := 0; j < len(columns); j++ {
				cellContent := strings.ToLower(data[(row+i)%len(data)][j])
				if strings.Contains(cellContent, text) {
					table.Select(((row+i)%len(data))+1, 0)
					return true
				}
			}
		}
		return false
	}

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
					dn := rows[row-1]
					cmd := exec.Command("ldapvi", "-s", "base", "-b", dn)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						log.Printf("ldapvi: " + err.Error())
						time.Sleep(5 * time.Second)
					}
					rows, columns, data = ldapSearch(LdapDN, LdapFilter, LdapAttrs)
					tviewFillTable(table, columns, data)
				})
			case '/':
				row, _ := table.GetSelection()
				row--
				search := tview.NewInputField()
				search.SetLabel("Search: ")
				search.SetFieldBackgroundColor((tcell.ColorBlack))
				search.SetChangedFunc(func(text string) {
					if tviewSearch(row, text) {
						search.SetFieldTextColor(tcell.ColorWhite)
					} else {
						search.SetFieldTextColor((tcell.ColorRed))
					}
				})
				search.SetDoneFunc(func(key tcell.Key) {
					lastSearch = search.GetText()
					flex.RemoveItem(lastLine)
					lastLine = tview.NewTextView().SetText(fmt.Sprintf("Last search: %q from line %d", lastSearch, row))
					flex.AddItem(lastLine, 1, 0, false)
					app.SetFocus(table)
				})
				flex.RemoveItem(lastLine)
				lastLine = search
				flex.AddItem(lastLine, 1, 0, false)
				app.SetFocus(search)

			case 'n':
				row, _ := table.GetSelection()
				tviewSearch(row, lastSearch)
				flex.RemoveItem(lastLine)
				lastLine = tview.NewTextView().SetText(fmt.Sprintf("Searching again: %q from line %d", lastSearch, row))
				flex.AddItem(lastLine, 1, 0, false)
			}
		}
		return event
	})
	text.SetBackgroundColor(tcell.ColorBlue)
	text.SetDynamicColors(true)
	text.SetText(" [yellow]q:quit   e:edit   f:filter   s:sort   /:search   n:next")
	flex.SetBackgroundColor(tcell.ColorRed)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(text, 1, 0, false)
	lastLine = tview.NewBox()
	flex.AddItem(lastLine, 1, 0, false)
	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
