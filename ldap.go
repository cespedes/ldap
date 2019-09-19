package main

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap"
)

func ldapSearch(baseDN string, filter string, reqAttributes []string) (rows []string, attributes []string, table [][]string) {
	l, err := ldap.DialURL(Config["config"]["server"])
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		baseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,        // The filter to apply
		reqAttributes, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	mapAttrs := make(map[string]bool)
	for _, entry := range sr.Entries {
		for _, attr := range entry.Attributes {
			mapAttrs[attr.Name] = true
		}
	}
	for _, a := range strings.Split(Config["attributes_order"]["order"], " ") {
		if mapAttrs[a] {
			attributes = append(attributes, a)
			delete(mapAttrs, a)
		}
	}
	for a := range mapAttrs {
		attributes = append(attributes, a)
	}

	var result [][]string

	for _, entry := range sr.Entries {
		if *flagDebug {
			log.Println()
			for _, a := range entry.Attributes {
				log.Printf("%v: %v", a.Name, a.Values)
			}
		}
		rows = append(rows, entry.DN)
		var line []string
		for _, attr := range attributes {
			line = append(line, strings.Join(entry.GetAttributeValues(attr), " & "))
		}
		result = append(result, line)
	}
	return rows, attributes, result
}
