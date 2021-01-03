package wikiprov

// Hard-codes configuration of the wikiprov package.
//
// An example API query we need to construct:
//
//  https://www.wikidata.org/w/api.php?action=query&format=json&prop=revisions&titles=Q5381415&rvlimit=200&rvprop=ids|user|comment|timestamp|sha1
//
// We'll also use some of these values to build a permalink which looks
// as follows:
//
//	https://www.wikidata.org/w/index.php?title=Q178051&oldid=1301912874&format=json
//

const agent string = "wikiprov/0.0.2 (https://github.com/ross-spencer/wikiprov/; all.along.the.watchtower+github@gmail.com)"

var wikibaseAPI = "https://www.wikidata.org/w/api.php"
var wikidataBase = "https://www.wikidata.org/w/index.php"

var format = "json"
var action = "query"
var prop = "revisions"

var revisionPropertiesDefault = [...]string{"ids", "user", "comment", "timestamp", "sha1"}
