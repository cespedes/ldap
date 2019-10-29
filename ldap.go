package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
)

func freePort() string {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic("freePort: binding to new port: " + err.Error())
	}
	defer ln.Close()

	addr := ln.Addr()
	tcp := addr.(*net.TCPAddr)

	return strconv.Itoa(tcp.Port)
}

func ldapDial(c Config) *ldap.Conn {
	if *flagSSH == "" {
		l, err := ldap.DialURL(c.LdapServer)
		if err != nil {
			log.Fatal("ldap.Dial: " + err.Error())
		}
		return l
	}
	url, err := url.Parse(c.LdapServer)
	if err != nil {
		log.Fatal("Unable to parse " + c.LdapServer + ": " + err.Error())
	}
	host := url.Hostname()
	remotePort := url.Port()
	if remotePort == "" {
		switch url.Scheme {
		case "ldap":
			remotePort = "389"
		case "ldaps":
			remotePort = "636"
		default:
			log.Fatal("Unknown scheme \"" + url.Scheme + "\"")
		}
	}
	localPort := freePort()

	cmd := exec.Command(
		"ssh", "-f",
		"-L", fmt.Sprintf("%s:%s:%s", localPort, host, remotePort),
		*flagSSH, "sleep", "60")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Connecting to ssh proxy: " + err.Error())
	}
	if url.Scheme == "ldap" {
		l, err := ldap.Dial("tcp", net.JoinHostPort("localhost", localPort))
		if err != nil {
			log.Fatal("ldap.Dial: " + err.Error())
		}
		return l
	}
	if url.Scheme == "ldaps" {
		l, err := ldap.DialTLS("tcp",
			net.JoinHostPort("localhost", localPort),
			&tls.Config{
				ServerName: host,
			})
		if err != nil {
			log.Fatal("ldap.Dial: " + err.Error())
		}
		return l
	}
	log.Fatal("Unknown scheme \"" + url.Scheme + "\"")
	return nil
}

func ldapSearch(c Config) (dnList []string, attributes []string, table [][]string) {
	l := ldapDial(c)
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
		log.Fatal("ldapSearch: Search: " + err.Error())
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
