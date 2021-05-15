package wikiprov

// Consts and variables used internally to request or create the data
// that we want.

const agent string = "wikiprov/0.1.0 (https://github.com/ross-spencer/wikiprov/; all.along.the.watchtower+github@gmail.com)"

const defaultWikibaseAPI = "https://www.wikidata.org/w/api.php"
const wikibasePermaURL = "https://www.wikidata.org/w/index.php"

var wikibaseAPI = defaultWikibaseAPI
var wikibasePermalinkBase = wikibasePermaURL

var format = "json"
var action = "query"
var prop = "revisions"

var revisionPropertiesDefault = [...]string{"ids", "user", "comment", "timestamp", "sha1"}
