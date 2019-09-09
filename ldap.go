package main

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap"
)

func ldap_search(filter string, reqAttributes []string) (attributes []string, table [][]string) {
	l, err := ldap.Dial("tcp", Config["config"]["server"] + ":" + Config["config"]["port"])
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		Config["config"]["basedn"], // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, // The filter to apply
		reqAttributes,                    // A list attributes to retrieve
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
	var seenAttrs []string
	for _, a := range strings.Split(Config["attributes_order"]["order"], " ") {
		if mapAttrs[a] {
			seenAttrs = append(seenAttrs, a)
			delete(mapAttrs, a)
		}
	}
	for a, _ := range mapAttrs {
		seenAttrs = append(seenAttrs, a)
	}

	var result [][]string

	for _, entry := range sr.Entries {
		var line []string
		if *flagDebug {
			log.Println()
			for _, a := range entry.Attributes {
				log.Printf("%v: %v", a.Name, a.Values)
			}
		}
		for _, attr := range seenAttrs {
			line = append(line, strings.Join(entry.GetAttributeValues(attr), " & "))
		}
		result = append(result, line)
	}
	return seenAttrs, result
}
