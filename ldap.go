package main

import (
	"fmt"
	"log"
	"github.com/go-ldap/ldap"
	"os"
	"strings"
)

func ldap_search(filter string, attributes []string) [][]string {
	l, err := ldap.Dial("tcp", Config["config"]["server"] + ":" + Config["config"]["port"])
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		Config["config"]["basedn"], // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, // The filter to apply
		attributes,                    // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var result [][]string
	for _, entry := range sr.Entries {
		var line []string
		for _, attr := range attributes {
			line = append(line, strings.Join(entry.GetAttributeValues(attr), " & "))
		}
		result = append(result, line)
	}
	return result
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tldaporg <filter> [<attr>...]\n")
}

func main() {
	var filter, real_filter string
	var attrs, real_attrs []string

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	filter = os.Args[1]
	real_filter = Config["filters"][filter]
	if real_filter == "" {
		real_filter = filter
	}
	if len(os.Args) > 2 {
		attrs = os.Args[2:]
	} else {
		tmp := Config["default_attributes"][filter]
		if tmp != "" {
			attrs = strings.Split(tmp, " ")
		} else {
			fmt.Fprintf(os.Stderr, "Error: no default attributes for filter \"%s\"\n", filter)
			os.Exit(1)
		}
	}
	for _, name := range attrs {
		tmp := Config["attributes"][name]
		if tmp == "" {
			real_attrs = append(real_attrs, name)
		} else {
			real_attrs = append(real_attrs, tmp)
		}
	}

	result := ldap_search(real_filter, real_attrs)

	write_orgtable(os.Stdout, attrs, result)
}
