package main

import (
	"log"
	"sort"
	"strings"

	"github.com/go-ldap/ldap"
)

func ldapSearch(c Config) (dnList []string, attributes []string, table [][]string) {
	l, err := ldap.DialURL(c.LdapServer)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		c.LdapDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		c.LdapFilter, // The filter to apply
		c.LdapAttrs,  // A list attributes to retrieve
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
	for _, a := range c.AttributesOrder {
		if mapAttrs[a] {
			attributes = append(attributes, a)
			delete(mapAttrs, a)
		}
	}
	for a := range mapAttrs {
		attributes = append(attributes, a)
	}

	if c.OrderBy != "" {
		for _, name := range attributes {
			if c.OrderBy == name || c.OrderBy == c.ReverseAlias[name] {
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
		if c.ReverseAlias[a] != "" {
			attributes[i] = c.ReverseAlias[a]
		}
	}

	return dnList, attributes, table
}
