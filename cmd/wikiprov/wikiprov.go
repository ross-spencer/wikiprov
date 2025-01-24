package main

// Reference application for accessing provenance information via
// Wikidata. Generic instances, e.g. your own Wikibase are not handled
// by this app.

import (
	"flag"
	"fmt"
	"os"

	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

var (
	demo    bool
	history int
	qid     string
	vers    bool
)

func init() {
	flag.BoolVar(&demo, "demo", false, "Run the tool with a demo value and all provenance")
	flag.IntVar(&history, "history", 10, "length of history to return")
	flag.StringVar(&qid, "qid", "", "QID to look up provenance for")
	flag.BoolVar(&vers, "version", false, "Return version")
}

func main() {

	flag.Parse()
	if vers {
		fmt.Fprintf(os.Stderr, "%s \n", wikiprov.Version())
		os.Exit(0)
	} else if flag.NFlag() == 0 {
		fmt.Fprintln(os.Stderr, "wikiprov: return info about a QID from Wikidata")
		fmt.Fprintln(os.Stderr, "usage: wikiprov <QID e.g. Q27229608> {options}              ")
		fmt.Fprintln(os.Stderr, "                                     OPTIONAL: [-history] ...")
		fmt.Fprintln(os.Stderr, "                                     OPTIONAL: [-version]   ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "output: [JSON]   {wikidataProvenace}")
		fmt.Fprintf(os.Stderr, "output: [STRING] '%s ...'\n\n", wikiprov.Version())
		flag.Usage()
		os.Exit(0)
	}

	if demo {
		var demoQID = "Q49300657"
		res, _ := wikiprov.GetWikidataProvenance(demoQID, 10)
		fmt.Println(res)
		return
	}

	if qid == "" {
		fmt.Println("please provide a QID to lookup...")
		return
	}

	res, err := wikiprov.GetWikidataProvenance(qid, history)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
