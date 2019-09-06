package main

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap"
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
