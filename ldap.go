package main

import (
	"log"
	"sort"
	"strings"

	"github.com/go-ldap/ldap"
)

func ldapSearch(baseDN string, filter string, reqAttributes []string, orderBy string) (dnList []string, attributes []string, table [][]string) {
	l, err := ldap.DialURL(LdapServer)
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
	for _, a := range AttributesOrder {
		if mapAttrs[a] {
			attributes = append(attributes, a)
			delete(mapAttrs, a)
		}
	}
	for a := range mapAttrs {
		attributes = append(attributes, a)
	}

	if orderBy != "" {
		for _, name := range attributes {
			if orderBy == name || orderBy == ReverseAlias[name] {
				sort.Slice(sr.Entries, func(a, b int) bool {
					return sr.Entries[a].GetAttributeValue(name) < sr.Entries[b].GetAttributeValue(name)
				})
				break
			}
		}
	}

	for _, entry := range sr.Entries {
		dnList = append(dnList, entry.DN)
		var line []string
		for _, attr := range attributes {
			line = append(line, strings.Join(entry.GetAttributeValues(attr), " & "))
		}
		table = append(table, line)
	}
	for i, a := range attributes {
		if ReverseAlias[a] != "" {
			attributes[i] = ReverseAlias[a]
		}
	}

	return dnList, attributes, table
}
