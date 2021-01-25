// Package spargo is a Wrapper for the generic spargo package:
//
//    * github.com/ross-spencer/spargo/pkg/spargo
//
// The package exists to enable to inclusion of Wikibase provenance in
// those results. Where spargo is a generic package this version is
// specific to Wikidata implementations on-top of Wikibase.
package spargo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/ross-spencer/spargo/pkg/spargo"
	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

// DefaultAgent as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
const DefaultAgent = spargo.DefaultAgent

// Binding as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type Binding = spargo.Binding

// Item as it exists in the spargo package exported to enable dropping
// this package into host packages/executables.
type Item = spargo.Item

// SPARQLClient as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type SPARQLClient = spargo.SPARQLClient

// SPARQLResult as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type SPARQLResult = spargo.SPARQLResult

// WikiProv wraps spargo's standard results so that we can attach
// provenance without attempting to modify the generic capabilities of
// the wikiprov's sister package.
type WikiProv struct {
	Head       map[string]interface{} `json:"head"`
	Binding    `json:"results"`
	Provenance []wikiprov.Provenance `json:"provenance,omitempty"`
}

// maxChannels determines the number of channels to use in requests to
// Wikibase for its provenance data. Ostensibly it's a throttle.
// Wikidata will return an error if we ask for too much too quickly, 20
// caused an error previously for over 1000 records. 10 seems to work
// fairly well. Without requesting this information in threads,
// processing can be pretty slow.
var maxChannels = 10

// SPARQLWithProv is used to query the Wikidata query service and attach
// Wikibase provenance. History can be configured as well as the number
// of threads used to connect to Wikibase. The key provided this
// function must exist as a parameter in the SPARQL query, e.g. SELECT
// `?uri` where `?uri` is the key. This parameter must also be a
// Wikidata IRI from which the QID will be returned. The QID is then
// used to grab the provenance information for the record. If key is
// empty then provenance functions will not be called.
func SPARQLWithProv(endpoint string, queryString string, key string, lenHistory int, threads int) (WikiProv, error) {
	sparqlMe := SPARQLClient{}
	sparqlMe.ClientInit(endpoint, queryString)
	res := sparqlMe.SPARQLGo()
	provResults := WikiProv{}
	provResults.Head = res.Head
	provResults.Binding = res.Results
	if key == "" {
		return provResults, nil
	}
	if threads > maxChannels {
		threads = maxChannels
	}
	err := provResults.attachProvenance(key, lenHistory, threads)
	if err != nil {
		return WikiProv{}, err
	}
	return provResults, nil
}

// validateIRI will attempt to perform some basic validation on IRI's
// we're trying to retrieve provenance information for. We need to build
// up a set of rules.
func validateIRI(iri string) bool {
	const wikidata string = "wikidata.org"
	const statement string = "statement"
	if !strings.Contains(iri, wikidata) {
		return false
	}
	if strings.Contains(iri, statement) {
		return false
	}
	return true
}

// getQID will retrieve the QID from a Wikidata IRI. It can handle
// Properties which require a special suffice, and entities which are
// our standard QIDs e.g. Q12345.
func getQID(iri string) (string, error) {
	const prop string = "prop"
	const property string = "Property"
	parsedIRI, err := url.Parse(iri)
	if err != nil {
		return "", err
	}
	qid := path.Base(parsedIRI.Path)
	if strings.Contains(iri, prop) {
		return fmt.Sprintf("%s:%s", property, qid), nil
	}
	return qid, nil
}

// AttachProvenance will attach WikiBase provenance to SPARQL results
// from Wikidata.
func (sparql *WikiProv) attachProvenance(key string, lenHistory int, threads int) error {
	var qids map[string]bool
	qids = make(map[string]bool)
	for _, value := range sparql.Bindings {
		wikidataIRI := value[key].Value
		if !validateIRI(wikidataIRI) {
			continue
		}
		qid, err := getQID(wikidataIRI)
		if err != nil {
			return err
		}
		qids[qid] = false
	}
	if len(qids) < 1 {
		return fmt.Errorf("No results returned from given key")
	}
	var uniqueQIDs []string
	for key := range qids {
		uniqueQIDs = append(uniqueQIDs, key)
	}
	provCache := getProvThreaded(uniqueQIDs, lenHistory, threads)
	sparql.Provenance = provCache
	return nil
}

// getProvThreaded goes out to Wikibase and collects the provenance
// associated with a record. The function takes an argument that limits
// the number of channels to be used to do work to provide some level
// of throttling and to also increase performance of this. For ~5000
// records this can take 15 minutes without concurrency.
func getProvThreaded(qids []string, lenHistory int, maxChan int) []wikiprov.Provenance {
	ch := make(chan wikiprov.Provenance)
	var mutex sync.Mutex
	counter := 0
	for channels := 0; channels < maxChan; channels++ {
		go func(ch chan wikiprov.Provenance, mutex *sync.Mutex) {
			for {
				mutex.Lock()
				idx := counter
				counter++
				mutex.Unlock()
				if counter > len(qids) {
					// Finished processing, exit.
					return
				}
				qid := qids[idx]
				// Retrieve the provenance information from Wikibase.
				prov := getProvenance(qid, lenHistory)
				ch <- prov
			}
		}(ch, &mutex)
	}
	var provCache []wikiprov.Provenance
	provCache = make([]wikiprov.Provenance, len(qids))
	getData(ch, provCache)
	return provCache
}

// getData invokes the go routines and then adds the results to the
// provenance array.
func getData(ch <-chan wikiprov.Provenance, provCache []wikiprov.Provenance) {
	for idx := 0; idx < len(provCache); idx++ {
		provCache[idx] = <-ch
	}
}

// getProvenance is a helper which is used to call wikiprov's primary
// function collecting provenance for a record from the underlying
// Wikibase implementation.
func getProvenance(qid string, lenHistory int) wikiprov.Provenance {
	prov, err := wikiprov.GetWikidataProvenance(qid, lenHistory)
	if err != nil {
		panic(err)
	}
	return prov
}

// String will return a summary of a Wikiprov structure as JSON.
func (sparql WikiProv) String() string {
	str, err := json.MarshalIndent(sparql, "", "  ")
	if err != nil {
		return ""
	}

	// THe encoder now escapes these values, this is for browser
	// compatibility, and I don't think it matters to us too much.
	//
	//    * https://stackoverflow.com/a/24657016
	//
	str = bytes.Replace(str, []byte("\\u003c"), []byte("<"), -1)
	str = bytes.Replace(str, []byte("\\u003e"), []byte(">"), -1)
	str = bytes.Replace(str, []byte("\\u0026"), []byte("&"), -1)

	return fmt.Sprintf("%s", str)
}
