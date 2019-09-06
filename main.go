package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "\tldaporg [<options>] <filter> [<attr>...]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "The options are:")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

var (
	flagEdit = flag.Bool("e", false, "Edit (not view) LDAP information (not implemented)")
	flagSort = flag.String("s", "", "Sort by that attribute")
	// flagDN   = flag.String("b", "", "Use this Base DN (not implemented)")
)

func main() {
	var filter, real_filter string
	var attrs, real_attrs []string

	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}
	filter = flag.Args()[0]
	real_filter = Config["filters"][filter]
	if real_filter == "" {
		real_filter = filter
	}
	if len(flag.Args()) > 1 {
		attrs = flag.Args()[1:]
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

	if (*flagSort != "") {
		for i, name := range attrs {
			if name == *flagSort {
				sort.Slice(result, func(a, b int) bool { return result[a][i] < result[b][i] })
				goto sortDone
			}
		}
		log.Fatal("Cannot sort by " + *flagSort + " (unknown attribute)")
sortDone:
	}

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
