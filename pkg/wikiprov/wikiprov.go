// Package wikiprov provides functions to enable simple reification,
// i.e. provenance/fixity of Wikidata entities. The module can be
// extended to other Wikibase sites in the future.
//
// The package constructs an API call for the Wikibase Query endpoint
// and then uses that data to create a normalized "provenance" struct
// that should be easier to work with for the caller.
//
// An example API query we need to construct:
//
//  https://www.wikidata.org/w/api.php?action=query&format=json&prop=revisions&titles=Q5381415&rvlimit=200&rvprop=ids|user|comment|timestamp|sha1
//
// We'll also use some of these values to build a permalink for that
// provenance struct which looks as follows:
//
//  https://www.wikidata.org/w/index.php?title=Q178051&oldid=1301912874&format=json
//
//  https://www.wikidata.org/w/index.php?title=<QID>&oldid=<REVISION>&format=json
//
package wikiprov

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getRevisionProperties() string {
	return strings.Join(revisionPropertiesDefault[:], "|")
}

// buildRequest will build the request we want to send to Wikibase.
// An error is returned if the request is malformed.
func buildRequest(id string, history int) (*http.Request, error) {
	const paramFormat = "format"
	const paramAction = "action"
	const paramTitles = "titles"
	const paramProps = "prop"
	const paramLimit = "rvlimit"
	const paramRevisionProps = "rvprops"

	req, err := http.NewRequest("GET", wikibaseAPI, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set(paramFormat, format)
	query.Set(paramAction, action)
	query.Set(paramTitles, id)
	query.Set(paramProps, prop)
	query.Set(paramLimit, fmt.Sprintf("%d", history))
	query.Set(paramRevisionProps, getRevisionProperties())

	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", agent)

	return req, nil
}

// GetWikidataProvenance requests the entity data we need from the
// Wikibase API and returns a structure containing the information that
// we're interested in, augmented with a permalink to the record.
func GetWikidataProvenance(id string, history int) (Provenance, error) {

	request, err := buildRequest(id, history)
	if err != nil {
		return Provenance{}, err
	}

	var client http.Client

	resp, err := client.Do(request)
	if err != nil {
		return Provenance{}, err
	}

	const expectedCode int = 200
	if resp.StatusCode != expectedCode {
		responseError := ResponseError{}
		return Provenance{}, responseError.makeError(200, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return Provenance{}, err
	}

	var wdRevisions wdRevisions

	err = json.Unmarshal(data, &wdRevisions)
	if err != nil {
		return Provenance{}, err
	}

	return wdRevisions.normalize(), nil
}

// Version returns the agent string for this package.
func Version() string {
	return agent
}

// SetWikibaseAPIURL lets the caller configure its own Wikibase API
// service to connect to.
func SetWikibaseAPIURL(newURL string) {
	wikibaseAPI = newURL
}

// SetWikibasePermalinkBaseURL lets the caller configure the Wikibase
// base URL for the permalink that needs to be built.
func SetWikibasePermalinkBaseURL(newURL string) {
	wikibasePermalinkBase = newURL
}
