package spargo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

// TestGetProvThreadedError will test the response from the function
// where a non-expected response is returned from the server.
func TestGetProvThreadedError(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		res.Write([]byte("no value"))
	}))
	defer func() { testServer.Close() }()

	wikiprov.SetWikibaseAPIURL(testServer.URL)

	for _, val := range errorTests {

		provs := getProvThreaded(val.qids, 5, val.threads)

		if len(provs) != len(val.qids) {
			t.Errorf("Despite testing an error condition results returned are not correct length. Got '%d', expected '%d'",
				len(provs),
				len(val.qids),
			)
		}

		for _, prov := range provs {
			if prov.Error == nil {
				t.Errorf("Expecting a non-'nil' error from getProvThreaded, received: '%s'", prov.Error)
			}
			responseError := wikiprov.ResponseError{}
			if !errors.As(prov.Error, &responseError) {
				t.Errorf("Unexpected error condition returned, expecting: '%s' received %s",
					responseError,
					prov.Error,
				)
			}
		}
	}
}

// TestGetProvThreaded will test the return of provenance results and
// parsing them given a number of different thread numbers, e.g. do we
// get the right number of results for the number of QIDs being queried
// over the as many threads as we specify.
func TestGetProvThreaded(t *testing.T) {

	for _, test := range threadTests {

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(test.result))
		}))
		defer func() { testServer.Close() }()

		wikiprov.SetWikibaseAPIURL(testServer.URL)

		// All the results can be the same. They just need to all be
		// accounted for given various different configurations of
		// threads etc. If there is an opportunity then these tests can
		// be expanded to be more varied.

		provs := getProvThreaded(test.qids, 5, 10)

		if len(provs) != len(test.qids) {
			t.Errorf("Results length from getProvThreaded: '%d' not what was expected: '%d'",
				len(provs),
				len(test.qids),
			)
		}

		// Expected test output. Results from getProvThreaded should match.
		testProvOutput := wikiprov.Provenance{}
		testProvOutput.Title = "Q12345"
		testProvOutput.Revision = 2600
		testProvOutput.Modified = "2020-08-31T23:13:00Z"
		testProvOutput.Permalink = "https://www.wikidata.org/w/index.php?oldid=2600&title=Q12345"
		testProvOutput.History = append(testProvOutput.History, "2020-08-31T23:13:00Z (oldid: 2600): 'Emmanuel Goldstein' edited: 'edit comment #1'")
		testProvOutput.History = append(testProvOutput.History, "2020-08-01T23:13:00Z (oldid: 1000): 'Robert Smith' edited: 'edit comment #2'")
		testProvOutput.Error = nil

		if !reflect.DeepEqual(provs[0], testProvOutput) {
			t.Errorf("Provenance result structure: '%v' does not match expected output structure: '%v'",
				provs[0],
				testProvOutput,
			)
		}

		if len(provs) > 1 {
			if !reflect.DeepEqual(provs[len(provs)-1], testProvOutput) {
				t.Errorf("Provenance result structure: '%v' does not match expected output structure: '%v'",
					provs[0],
					testProvOutput,
				)
			}
		}
	}
}

// TestSPARQLWithProvError will test what happens when there is a server
// error and provenance cannot be attached to the SPARQL results that
// we already have.
func TestSPARQLWithProvError(t *testing.T) {

	// sparqlTestServer returns valid Wikidata SPARQL JSON.
	sparqlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(wikidataResultsJSON))
	}))
	defer func() { sparqlTestServer.Close() }()

	// apiTestServer returns an unexpected response and a value that
	// doesn't need to be handled.
	apiTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		res.Write([]byte("no value"))
	}))
	defer func() { apiTestServer.Close() }()

	// Replace the API test server with our custom URL.
	wikiprov.SetWikibaseAPIURL(apiTestServer.URL)

	lenResults := 5
	threads := 10

	prov, err := SPARQLWithProv(sparqlTestServer.URL, "testQuery", "uri", lenResults, threads)

	if !reflect.DeepEqual(prov, WikiProv{}) {
		t.Errorf("Expected an empty WikiProv{} struct to be returned, returned: '%s'", prov)
	}

	if err == nil {
		t.Errorf("Anticipating an error from SPARQLWithProv, received 'nil': %s", err)
	}

	if !errors.Is(err, ErrProvAttach) {
		t.Errorf("Expecting error: '%s' but received: '%s'", ErrProvAttach, err)
	}
}

// TestSPARQLWithProv is used to look at the provenance attached to a
// SPARQL result from a Wikidata like service and ensures that the
// data is constructed as we need.
func TestSPARQLWithProv(t *testing.T) {

	// The number of records being queried that we expect provenance
	// entries for.
	const expectedResultsLength int = 6

	for _, val := range []int{1, 5, 7, 10, 100} {

		// sparqlTestServer returns valid Wikidata SPARQL JSON. The test
		// string used used an example.com URL throughout to help us
		// test against Wikidata artifacts where other Wikibase sites
		// can be used.
		sparqlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(wikidataResultsJSONExampleDotCom))
		}))
		defer func() { sparqlTestServer.Close() }()

		// apiTestServer returns an unexpected response and a value that
		// doesn't need to be handled.
		apiTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(attachedProvenance))
		}))
		defer func() { apiTestServer.Close() }()

		// Replace the API test server with our custom URL.
		wikiprov.SetWikibaseAPIURL(apiTestServer.URL)

		lenResults := 2 // Unimportant as the results from the tests are deterministic.
		threads := val

		// Using the example.com URL base we want to make sure the
		// results can reflect other services and not just Wikidata.
		wikiprov.SetWikibasePermalinkBaseURL("http://example.com")

		provs, err := SPARQLWithProv(sparqlTestServer.URL, "testQuery", "uri", lenResults, threads)

		if err != nil {
			t.Errorf("Unexpected error '%s' from SPARQLWithProv", err)
		}

		// All the results will be the same as we're only returning one
		// value from the test server, but they must all be accounted
		// for and correct for the different number of threads we're
		// using.

		if len(provs.Provenance) != expectedResultsLength {
			t.Errorf("Expected results length '%d', but got '%d'", expectedResultsLength, len(provs.Provenance))
		}

		// Expected test output. Results from getProvThreaded should match.
		testProvOutput := wikiprov.Provenance{}
		testProvOutput.Title = "Q12345"
		testProvOutput.Revision = 2600
		testProvOutput.Modified = "2020-08-31T23:13:00Z"
		testProvOutput.Permalink = "http://example.com?oldid=2600&title=Q12345"
		testProvOutput.History = append(testProvOutput.History, "2020-08-31T23:13:00Z (oldid: 2600): 'Emmanuel Goldstein' edited: 'edit comment #1'")
		testProvOutput.History = append(testProvOutput.History, "2020-08-01T23:13:00Z (oldid: 1000): 'Robert Smith' edited: 'edit comment #2'")
		testProvOutput.Error = nil

		// Test some characteristics from the provenance struct and ensure that
		// we are getting back the latest values.
		for _, prov := range provs.Provenance {

			if !reflect.DeepEqual(prov, testProvOutput) {
				t.Errorf("Provenance result structure: '%v' does not match expected output structure: '%v'",
					prov,
					testProvOutput,
				)
			}
		}
	}
}
