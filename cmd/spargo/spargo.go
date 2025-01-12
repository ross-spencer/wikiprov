package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ross-spencer/wikiprov/pkg/spargo"
)

// SHEBANG provides some way of recognizing a .sparql file compatible
// with spargo, aka. our .sparql magic number.
var SHEBANG string = "#!"

// ENDPOINT must be specified in a .sparql file so that a query can be
// sent to the appropriate SPARQL endpoint.
const ENDPOINT string = "ENDPOINT"

// WIKIBASE describes the Wikibase URL from where we will retrieve
// provenance information.
const WIKIBASEURL string = "WIKIBASEURL"

// HISTORY describes the length of history to return.
const HISTORY string = "HISTORY"

// SUBJECTPARAM describes the ?param to use as the subject of the query
// for which we want provenance for.
const SUBJECTPARAM string = "SUBJECTPARAM"

// wikiEndpoint allows us to check that only the Wikidata endpoint is
// supplied to the utility.
const wikiEndpoint string = "https://query.wikidata.org/sparql"

var (
	vers       bool
	query      string
	endpoint   string
	param      string
	lenHistory int
	threads    int
)

type wbQuery struct {
	url      string
	wikibase string
	query    string
	param    string
	subject  string
	history  int
}

func (wb wbQuery) String() string {
	return fmt.Sprintf("---\nurl: %s\nwikibaseURL: %s\nquery: %s\nhistory:%d\n---\n", wb.url, wb.wikibase, wb.query, wb.history)
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "endpoint to query")
	flag.StringVar(&query, "query", "", "sparql query to run")
	flag.StringVar(&param, "param", "", "for provenance a SPARQL ?param needs to be specified that contains a Wikidata IRI")
	flag.IntVar(&lenHistory, "history", 5, "length of history to return to the caller")
	flag.IntVar(&threads, "threads", 10, "number of go routines to use to fetch provenance")
	flag.BoolVar(&vers, "version", false, "application version and user-agent")
}

// extractKey will extract a key from the input script.
func extractKey(line string, key string) string {
	str := strings.SplitN(line, "=", 2)
	if len(str) < 2 {
		fmt.Fprintf(os.Stderr, "cannot extract key '%s' from file", key)
		return ""
	}
	return strings.TrimSpace(str[1])
}

// Extract the query from the .sparql input.
func extractQuery(sparqlFile string) (wbQuery, error) {
	var err error
	var wb wbQuery
	for _, line := range strings.Split(sparqlFile, "\n") {
		if strings.HasPrefix(line, SHEBANG) {
			continue
		} else if line == "" {
			continue
		} else if strings.Contains(strings.ToUpper(line), ENDPOINT) {
			url := extractKey(line, ENDPOINT)
			if !strings.Contains(url, wikiEndpoint) {
				errString := fmt.Sprintf("endpoint does not look like a valid Wikidata endpoint: %s", url)
				err = fmt.Errorf(errString)
			}
			wb.url = url
		} else if strings.Contains(strings.ToUpper(line), WIKIBASEURL) {
			wbURL := extractKey(line, WIKIBASEURL)
			wb.wikibase = wbURL
		} else if strings.Contains(strings.ToUpper(line), HISTORY) {
			wbURL := extractKey(line, HISTORY)
			wb.history, _ = strconv.Atoi(wbURL)
		} else if strings.Contains(strings.ToUpper(line), SUBJECTPARAM) {
			wb.param = strings.Replace(extractKey(line, SUBJECTPARAM), "?", "", 1)
		} else {
			wb.query = wb.query + strings.TrimSpace(line) + "\n"
		}
	}
	return wb, err
}

func runQuery(sparqlFile string) {
	wb, err := extractQuery(sparqlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "connecting to: %s", wb)
	fmt.Fprintf(os.Stderr, "threads: %d\n", threads)
	if wb.param == "" {
		fmt.Fprintf(os.Stderr, "?param not set, not returning provenance for query\n")
	}
	provResults, err := spargo.SPARQLWithProv(wb.url, wb.query, wb.param, wb.history, threads)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	fmt.Println(provResults)
}

func isPipeInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if (info.Mode() & os.ModeNamedPipe) != 0 {
		return true
	}
	return false
}

// interpreterInput tests for a file as the second argument in a call to
// spargo.
//
// TODO: there may be another pattern here using Open, but we are also
// anticipating other arguments to the program at different times, so...
func interpreterInput() (bool, string) {
	if len(os.Args) == 2 {
		sparql := os.Args[1]
		if _, err := os.Stat(sparql); err == nil {
			return true, sparql
		} else if os.IsNotExist(err) {
			// Does not exist.
		} else {
			// Another error.
		}
	}
	return false, ""
}

func handlePipedInput() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return string(output)
}

func handleInterpreterInput(sparql string) string {
	data, err := ioutil.ReadFile(sparql)
	if err != nil {
		return ""
	}
	return string(data)
}

func main() {
	// Parse our input and let spargo generate a response.
	flag.Parse()
	if isPipeInput() {
		queryString := handlePipedInput()
		runQuery(queryString)
		os.Exit(0)
	} else {
		_, sparql := interpreterInput()
		if sparql != "" {
			query := handleInterpreterInput(sparql)
			runQuery(query)
			os.Exit(0)
		}
	}
	flag.Parse()
	if vers {
		fmt.Fprintf(os.Stderr, "%s (%s)\n", getVersion(), spargo.DefaultAgent)
		os.Exit(0)
	} else if flag.NFlag() == 0 {
		fmt.Fprintln(os.Stderr, "spargo (with provenance): run sparql queries from the command-line.")
		fmt.Fprintln(os.Stderr, "usage:  spargo {options}              ")
		fmt.Fprintln(os.Stderr, "                 OPTIONAL: [-sparql] ...")
		fmt.Fprintln(os.Stderr, "                 OPTIONAL: [-query]  ...")
		fmt.Fprintln(os.Stderr, "                 OPTIONAL: [-variable]  ...")
		fmt.Fprintln(os.Stderr, "                 OPTIONAL: [-version]   ")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "output: [JSON]   {url}")
		fmt.Fprintf(os.Stderr, "output: [STRING] '%s (%s) ...'\n\n", getVersion(), spargo.DefaultAgent)
		flag.Usage()
		os.Exit(0)
	}
}
