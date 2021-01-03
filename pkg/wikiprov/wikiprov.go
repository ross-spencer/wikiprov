// Package wikiprov provides functions to enable simple reification,
// i.e. provenance/fixity of Wikidata entities. The module can be
// extended to other Wikibase sites in the future.
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
func buildRequest(ids string) (*http.Request, error) {
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
	query.Set(paramTitles, ids)
	query.Set(paramProps, prop)
	query.Set(paramLimit, fmt.Sprintf("%d", revisionLimitDefault))
	query.Set(paramRevisionProps, getRevisionProperties())

	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", agent)

	return req, nil
}

// GetWikidataProvenance requests the entity data we need from the
// Wikibase API and returns a structure containing the information that
// we're interested in, augmented with a permalink to the record.
func GetWikidataProvenance(ids string) (Provenance, error) {

	request, err := buildRequest(ids)
	if err != nil {
		return Provenance{}, err
	}

	var client http.Client

	resp, err := client.Do(request)
	if err != nil {
		return Provenance{}, err
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
