package wikiprov

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// testInit allows us to reset any values we need to reset before our
// next test... The use of testInit() might point to a different pattern
// that we can use in another release.
func testInit() {
	wikibaseAPI = constructWikibaseAPIURL(defaultBaseURI)
	wikibasePermalinkBase = constructWikibaseIndexURL(defaultBaseURI)
}

// TestConstructAPIURL ensures that we correctly create the URL needed
// to talk to the Wikimedia API.
func TestConstructAPIURL(t *testing.T) {
	urlA := constructWikibaseAPIURL("http://example.com")
	urlB := constructWikibaseAPIURL("http://example.com/")
	// res is expected to be the same in both cases.
	res := "http://example.com/w/api.php"
	if urlA != res {
		t.Errorf("Incorrect URL created, expected: '%s', received: '%s'", res, urlA)
	}
	if urlB != res {
		t.Errorf("Incorrect URL created, expected: '%s', received: '%s'", res, urlB)
	}
}

// TestConstructIndexURL ensures that we correctly create the URL needed
// for the Wikimedia index page, e.g. for resolution of permalinks.
func TestConstructIndexURL(t *testing.T) {
	urlA := constructWikibaseIndexURL("http://example.com")
	urlB := constructWikibaseIndexURL("http://example.com/")
	// res is expected to be the same in both cases.
	res := "http://example.com/w/index.php"
	if urlA != res {
		t.Errorf("Incorrect URL created, expected: '%s', received: '%s'", res, urlA)
	}
	if urlB != res {
		t.Errorf("Incorrect URL created, expected: '%s', received: '%s'", res, urlB)
	}
}

// TestGetWikidataProvenance provides a simple test case to ensure that
// data isn't mutated between request and response and that Wikiprov's
// primary function returns something predictable from a sensible
// response.
func TestGetWikidataProvenance(t *testing.T) {

	testInit()

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(testJSON))
	}))
	defer func() { testServer.Close() }()

	// Replace Wikibase API URL with that of the test server.
	wikibaseAPI = testServer.URL

	// The integer here 1,000,000 is related to number of lines to request
	// from Wikidata. Wikidata handles that. We don't get to control the
	// return so precisely from this test.
	prov, err := GetWikidataProvenance("Q12345", 1000000)

	if err != nil {
		t.Errorf("We were expecting 'nil' error from this test, received: %s", err)
	}

	// The data we're testing is not directly from Wikibase but a
	// normalized format. I.e. it has been parsed as it came out of
	// Wikidata and converted int something we can use. For example,
	// the permalink URL is created in this library not Wikibase.
	title := "Q12345"
	revisions := 1419131078
	modified := "2021-05-11T20:17:31Z"

	permalink := "https://www.wikidata.org/w/index.php?oldid=1419131078&title=Q12345"

	// We're not doing anything with these values right now. An integer
	// work just as well if we were really desperate to change this.
	history := []string{"val1", "val2", "val3", "val4", "val5"}

	if prov.Title != title {
		t.Errorf("Provenance title '%s' is incorrect, expected: '%s'", prov.Title, title)
	}

	if prov.Revision != revisions {
		t.Errorf("Provenance title '%d' is incorrect, expected: '%d'", prov.Revision, revisions)
	}

	if prov.Modified != modified {
		t.Errorf("Provenance title '%s' is incorrect, expected: '%s'", prov.Modified, modified)
	}

	if prov.Permalink != permalink {
		t.Errorf("Provenance title '%s' is incorrect, expected: '%s'", prov.Permalink, permalink)
	}

	if len(prov.History) != len(history) {
		t.Errorf("Provenance title '%d' is incorrect, expected: '%d'", len(prov.History), len(history))
	}

	// Array is a predictable structure and we've created some user values to
	// test. Fairly arbitrary. We can find better tests in time.
	for userX, val := range prov.History {
		testUser := fmt.Sprintf("user%d", userX+1)
		if !strings.Contains(val, testUser) {
			t.Errorf("User string for '%s' has been lost somewhere in '%s'", testUser, val)
		}
	}
}

// TestGetWikidataProvenanceError provides a simple test case to ensure
// that something sensible is returned when the data that we receives
// is not what was expected. In this case, a 400 error.
func TestGetWikidataProvenanceError(t *testing.T) {

	testInit()

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		res.Write([]byte(testJSON))
	}))
	defer func() { testServer.Close() }()

	// Replace Wikibase API URL with that of the test server.
	wikibaseAPI = testServer.URL

	// The integer here 1,000,000 is related to number of lines to request
	// from Wikidata. Wikidata handles that. We don't get to control the
	// return so precisely from this test.
	prov, err := GetWikidataProvenance("Q12345", 1000000)

	const expectedErrorString string = "Returned unexpected status code"

	// Ensure that the error is not nil.
	if err == nil {
		t.Errorf("We were expecting an error from function, received: 'nil'")
	}

	// Ensure that it is the error we anticipated.
	responseError := ResponseError{}
	if err != nil && !errors.As(err, &responseError) {
		t.Errorf("We were expecting a specific error '%s' from this test, received: '%s'",
			responseError,
			err,
		)
	}

	// Ensure that provenance returned is empty.
	nilProv := Provenance{}
	if !reflect.DeepEqual(prov, nilProv) {
		t.Errorf("Function should return empty Provenance{} returned: %s", prov)
	}
}

// TestBuildRequest performs some rudimentary testing of the request
// builder. We can also improve on this in time.
func TestBuildRequest(t *testing.T) {

	testInit()

	req, err := buildRequest("Q12345", 1)
	if err != nil {
		t.Errorf("Expected 'nil' err from buildREquest, received: '%s'", err)
	}

	// The request to get provenance for the given QID, in this case,
	// QID needs to be:
	const expectedURL string = "https://www.wikidata.org/w/api.php?action=query&format=json&prop=revisions&rvlimit=1&rvprop=ids%7Cuser%7Ccomment%7Ctimestamp%7Csha1&titles=item%3AQ12345"
	// We test that here...
	if req.URL.String() != expectedURL {
		t.Errorf("Requested string we built isn't correct, \nreceived: '%s', \nexpected: '%s'",
			req.URL.String(),
			expectedURL,
		)
	}

	// Test that the user-agent is added as we'd like.
	builtAgent := req.Header["User-Agent"]
	if len(builtAgent) > 0 {
		if builtAgent[0] != agent {
			t.Errorf("Unexpected user agent: '%s', expected: '%s'",
				req.Header["User-Agent"][0],
				agent,
			)
		}
	} else {
		t.Errorf("User-Agent field not added to Wikibase response as required")
	}

}
