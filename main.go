package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "\tldaporg [<options>] <filter> [<attr>...]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "The options are:")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

var (
	flagOrg   = flag.Bool("o", false, "Use org-tables instead of tview")
	flagSort  = flag.String("s", "", "Sort by that attribute")
	flagDebug = flag.Bool("d", false, "Show debugging info")
	flagDN    = flag.String("b", "", "Use this Base DN")
)

var (
	// LdapFilter is used to restrict the LDAP query
	LdapFilter string

	// LdapAttrs is a list of attributes to ask in the LDAP query
	LdapAttrs []string

	// LdapDN is the Base DN in the LDAP query
	LdapDN string
)

func main() {
	var filter string
	var attrs []string

	flag.Parse()

	if len(flag.Args()) >= 1 {
		filter = flag.Args()[0]
		attrs = flag.Args()[1:]
	} else {
		filter = "(objectClass=*)"
		attrs = []string{"*", "createTimestamp", "modifyTimestamp"}
	}
	LdapFilter = Config["filters"][filter]
	if LdapFilter == "" {
		LdapFilter = filter
	}
	if *flagDN != "" {
		LdapDN = *flagDN
	} else {
		LdapDN = Config["config"]["basedn"]
	}
	tmp := Config["default_attributes"][filter]
	if tmp != "" {
		attrs = strings.Split(tmp, " ")
	} else {
		attrs = strings.Split(Config["config"]["attributes"], " ")
		if attrs == nil {
			fmt.Fprintf(os.Stderr, "Error: no default attributes for filter \"%s\"\n", filter)
			os.Exit(1)
		}
	}
	for _, name := range attrs {
		tmp := Config["attributes"][name]
		if tmp == "" {
			LdapAttrs = append(LdapAttrs, name)
		} else {
			LdapAttrs = append(LdapAttrs, tmp)
		}
	}

	rows, attrs, result := ldapSearch(LdapDN, LdapFilter, LdapAttrs)

	if *flagSort != "" {
		for i, name := range attrs {
			if name == *flagSort {
				sort.Slice(result, func(a, b int) bool { return result[a][i] < result[b][i] })
				goto sortDone
			}
		}
		log.Fatal("Cannot sort by " + *flagSort + " (unknown attribute)")
	sortDone:
	}

	if *flagOrg {
		writeOrgtable(os.Stdout, attrs, result)
		return
	}
	fmt.Println("len(rows) = ", len(rows))
	fmt.Println("len(attrs) = ", len(attrs))
	fmt.Println("len(result) = ", len(result))
	myTview(rows, attrs, result)
}
