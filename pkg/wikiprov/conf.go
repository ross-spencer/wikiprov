package wikiprov

// Consts and variables used internally to request or create the data
// that we want.

const agent string = "wikiprov/0.0.2 (https://github.com/ross-spencer/wikiprov/; all.along.the.watchtower+github@gmail.com)"

var wikibaseAPI = "https://www.wikidata.org/w/api.php"
var wikidataBase = "https://www.wikidata.org/w/index.php"

var format = "json"
var action = "query"
var prop = "revisions"

var revisionPropertiesDefault = [...]string{"ids", "user", "comment", "timestamp", "sha1"}
