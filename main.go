package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
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
	/*
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
	*/

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
	myTview(rows, attrs, result)
}
