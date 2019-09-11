package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
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
	flagOrg   = flag.Bool("o", false, "Use org-tables instead of tview")
	flagSort  = flag.String("s", "", "Sort by that attribute")
	flagDebug = flag.Bool("d", false, "Show debugging info")
	// flagDN   = flag.String("b", "", "Use this Base DN (not implemented)")
)

func main() {
	var filter, realFilter string
	var attrs, realAttrs []string

	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}
	filter = flag.Args()[0]
	realFilter = Config["filters"][filter]
	if realFilter == "" {
		realFilter = filter
	}
	if len(flag.Args()) > 1 {
		attrs = flag.Args()[1:]
	} else {
		tmp := Config["default_attributes"][filter]
		if tmp != "" {
			attrs = strings.Split(tmp, " ")
		} else {
			attrs = strings.Split(Config["default_attributes"]["default"], " ")
			if attrs == nil {
				fmt.Fprintf(os.Stderr, "Error: no default attributes for filter \"%s\"\n", filter)
				os.Exit(1)
			}
		}
	}
	for _, name := range attrs {
		tmp := Config["attributes"][name]
		if tmp == "" {
			realAttrs = append(realAttrs, name)
		} else {
			realAttrs = append(realAttrs, tmp)
		}
	}

	attrs, result := ldapSearch(realFilter, realAttrs)

	if *flagDebug {
		return
	}

	if *flagSort != "" {
		for i, name := range attrs {
			if name == *flagSort {
				sort.Slice(result, func(a, b int) bool { return result[a][i] < result[b][i] })
				goto sortDone
			}
		}
		log.Fatal("Cannot sort by " + *flagSort + " (unknown attribute)")
	sortDone:
	}

	if *flagOrg {
		writeOrgtable(os.Stdout, attrs, result)
		return
	}
	myTview(attrs, result)
}
