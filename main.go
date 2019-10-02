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
)

var (
	// LdapServer is the URL of the LDAP server to connect to
	LdapServer string

	// LdapDN is the Base DN in the LDAP query
	LdapDN string

	// LdapFilter is used to restrict the LDAP query
	LdapFilter string = "(cn=*)"

	// LdapAttrs is a list of attributes to ask in the LDAP query
	LdapAttrs []string = []string{"*", "createTimestamp", "modifyTimestamp", "structuralObjectClass"}

	// AttributesOrder is a list of possible attributes to display
	AttributesOrder []string

	// UserFilters is a definition of some filters using a short name
	UserFilters = make(map[string]string)

	// Alias is a map of friendly names to LDAP attributes
	Alias = make(map[string]string)

	// ReverseAlias is a map of LDAP attributes to friendly names
	ReverseAlias = make(map[string]string)
)

func main() {
	var attrs []string

	flag.Parse()

	if len(flag.Args()) >= 1 {
		LdapFilter = flag.Args()[0]
	}
	if len(flag.Args()) > 1 {
		LdapAttrs = flag.Args()[1:]
	}
	if UserFilters[LdapFilter] != "" {
		LdapFilter = UserFilters[LdapFilter]
	}
	if *flagDN != "" {
		LdapDN = *flagDN
	}

	dnList, attrs, result := ldapSearch(LdapDN, LdapFilter, LdapAttrs, *flagSort)

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
		dnList, columns, data = ldapSearch(LdapDN, LdapFilter, LdapAttrs, *flagSort)
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
		dnList, columns, data = ldapSearch(LdapDN, LdapFilter, LdapAttrs, *flagSort)
		t.FillTable(columns, data)
	})
	t.Run()
}
