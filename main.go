package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

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

	var pager string
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		pager = ""
	} else if _, err := exec.LookPath("pager"); err == nil {
		pager = "pager"
	} else if _, err := exec.LookPath("less"); err == nil {
		pager = "less"
	} else if _, err := exec.LookPath("more"); err == nil {
		pager = "more"
	} else {
		pager = ""
	}
	if pager != "" {
		cmd := exec.Command(pager)
		output, err := cmd.StdinPipe()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err == nil {
			cmd.Start()
			write_orgtable(output, attrs, result)
			output.Close()
			cmd.Wait()
			return
		}
	}
	write_orgtable(os.Stdout, attrs, result)
}
