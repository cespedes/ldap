package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	ini "github.com/glacjay/goini"
)

func init() {
	ldapConf, err := os.Open("/etc/ldap/ldap.conf")
	if err == nil {
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
					LdapServer = value
				case "BASE":
					LdapDN = value
				}
			}
		}
	}
	filename := os.Getenv("HOME") + "/.ldap.ini"
	userConf, err := ini.Load(filename)
	if err != nil {
		return
	}
	if userConf["general"]["server"] != "" {
		LdapServer = userConf["general"]["server"]
	}
	if userConf["general"]["basedn"] != "" {
		LdapDN = userConf["general"]["basedn"]
	}
	if userConf["general"]["attributes_order"] != "" {
		AttributesOrder = strings.Split(userConf["general"]["attributes_order"], " ")
	}
	UserFilters = userConf["filters"]
	if userConf["alias"] != nil {
		Alias = userConf["alias"]
	}
	for k, v := range Alias {
		ReverseAlias[v] = k
	}
	fmt.Println("Alias =", Alias)
	fmt.Println("ReverseAlias =", ReverseAlias)
}
