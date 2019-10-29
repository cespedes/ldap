package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cespedes/tableview"
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "\tldap [<options>] <filter> [<attr>...]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "The options are:")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

var (
	flagOrg  = flag.Bool("o", false, "Use org-tables instead of tview")
	flagSort = flag.String("s", "", "Sort by that attribute")
	flagDN   = flag.String("b", "", "Use this Base DN")
	flagSSH  = flag.String("ssh", "", "ssh server to connect as a proxy")
)

func main() {
	var attrs []string

	var c Config

	flag.Parse()

	if len(flag.Args()) >= 1 {
		c = config(flag.Args()[0])
	} else {
		c = config("")
	}
	if len(flag.Args()) > 1 {
		c.LdapAttrs = flag.Args()[1:]
	}
	if *flagDN != "" {
		c.LdapDN = *flagDN
	}

	if *flagSort != "" {
		c.OrderBy = *flagSort
	}
	dnList, attrs, result := ldapSearch(c)

	if *flagOrg {
		writeOrgtable(os.Stdout, attrs, result)
		return
	}

	t := tableview.NewTableView()
	t.FillTable(attrs, result)
	t.NewCommand('e', "edit", func(row int) {
		dn := dnList[row]
		cmd := exec.Command("ldapvi", "-s", "base", "-b", dn)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("ldapvi: " + err.Error())
			time.Sleep(5 * time.Second)
		}
		var columns []string
		var data [][]string
		dnList, columns, data = ldapSearch(c)
		t.FillTable(columns, data)
	})
	t.NewCommand('E', "Edit", func(row int) {
		dn := dnList[row]
		cmd := exec.Command("ldapvi", "-m", "-s", "base", "-b", dn)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("ldapvi: " + err.Error())
			time.Sleep(5 * time.Second)
		}
		var columns []string
		var data [][]string
		dnList, columns, data = ldapSearch(c)
		t.FillTable(columns, data)
	})
	t.Run()
}
