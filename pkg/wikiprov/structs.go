package wikiprov

// structs for wikiprov
//
// The primary JSON data we're interested in is the info struct from
// the API endpoint.
//
//	{
//		"entities": {
//			"Q27229608": {
//				"pageid": 29052990,
//				"ns": 0,
//				"title": "Q27229608",
//				"lastrevid": 784082439,
//				"modified": "2018-11-07T16:26:11Z",
//				"type": "item",
//				"id": "Q27229608"
//			}
//		},
//		"success": 1
//	}
//

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type infos struct {
	PageID     int    `json:"pageid"`
	NS         int    `json:"ns"`
	Title      string `json:"title"`
	LastRevID  int    `json:"lastrevid"`
	Modified   string `json:"modified"`
	EntityType string `json:"type"`
	ID         string `json:"id"`
}

type entities map[string]infos

type wdInfo struct {
	Entities  entities `json:"entities"`
	Success   int      `json:"success"`
	Permalink string
	ID        string
}

// Provenance provides simplified provenance information about a
// Wikidata record.
type Provenance struct {
	Title     string
	Revision  int
	Modified  string
	Permalink string
}

// buildPermalink creates a permalink for the entity requested at the
// time this function is called.
func (info *wdInfo) buildPermalink() string {
	const paramTitle = "title"
	const paramOldID = "oldid"
	const paramFormat = "format"
	req, _ := http.NewRequest("GET", wikidataBase, nil)
	query := req.URL.Query()
	title := info.Entities[info.ID].Title
	oldid := info.Entities[info.ID].LastRevID
	query.Set(paramTitle, title)
	query.Set(paramOldID, fmt.Sprintf("%d", oldid))
	query.Set(paramFormat, format)
	req.URL.RawQuery = query.Encode()
	return fmt.Sprintf("%s", req.URL)
}

// normalize simplifies the wdInfo structure so it can be easily used by
// the caller.
func (info *wdInfo) normalize() Provenance {
	prov := Provenance{}
	prov.Title = info.Entities[info.ID].Title
	prov.Revision = info.Entities[info.ID].LastRevID
	fmt.Println(prov.Revision)
	prov.Modified = info.Entities[info.ID].Modified
	prov.Permalink = info.buildPermalink()
	return prov
}

// String creates a human readable representation of the provenance
// struct.
func (prov Provenance) String() string {
	str, err := json.MarshalIndent(prov, "", "  ")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", str)
}
