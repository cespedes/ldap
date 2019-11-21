package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	ini "github.com/glacjay/goini"
)

var (
	defaultLdapServer string
	defaultLdapDN     string
	defaultLdapFilter = "(cn=*)"
	defaultLdapAttrs  = "* createTimestamp modifyTimestamp structuralObjectClass"
)

var userConf ini.Dict

func init() {
	if ldapConf, err := os.Open("/etc/ldap/ldap.conf"); err == nil {
		rd := bufio.NewReader(ldapConf)
		for {
			line, err := rd.ReadString('\n')
			if err != nil {
				break
			}
			var key, value string
			n, err := fmt.Sscan(line, &key, &value)
			if err == nil && n == 2 {
				switch key {
				case "URI":
					defaultLdapServer = value
				case "BASE":
					defaultLdapDN = value
				}
			}
		}
	}
	userConf, _ = ini.Load(os.Getenv("HOME") + "/.ldap.ini")
}

type Config struct {
	// LdapServer is the URL of the LDAP server to connect to
	LdapServer string

	// LdapDN is the Base DN in the LDAP query
	LdapDN string

	// LdapFilter is used to restrict the LDAP query
	LdapFilter string

	// LdapAttrs is a list of attributes to ask in the LDAP query
	LdapAttrs []string

	// AttributesOrder is a list of possible attributes to display
	AttributesOrder []string

	// Alias is a map of friendly names to LDAP attributes
	Alias map[string]string

	// ReverseAlias is a map of LDAP attributes to friendly names
	ReverseAlias map[string]string

	// OrderBy specifies the sorting order
	OrderBy string
}

func configGet(tag, name, def string) string {
	if s := userConf[tag][name]; s != "" {
		return s
	}
	if s := userConf["_default"][name]; s != "" {
		return s
	}
	return def
}

func config(name string) (c Config) {
	c.LdapServer = configGet(name, "server", defaultLdapServer)
	c.LdapDN = configGet(name, "basedn", defaultLdapDN)
	if name == "" {
		c.LdapFilter = configGet(name, "filter", defaultLdapFilter)
	} else {
		if s := userConf[name]["filter"]; s != "" {
			c.LdapFilter = s
		} else {
			c.LdapFilter = name
		}
	}
	c.LdapAttrs = strings.Split(configGet(name, "attributes", defaultLdapAttrs), " ")
	c.AttributesOrder = strings.Split(configGet(name, "order", ""), " ")
	c.OrderBy = configGet(name, "sort", "")
	c.Alias = userConf["_alias"]
	c.ReverseAlias = make(map[string]string)
	for k, v := range c.Alias {
		c.ReverseAlias[v] = k
	}
	return c
}
