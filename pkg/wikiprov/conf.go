package wikiprov

// Hard-codes configuration of the wikiprov package.
//
// An example API query we need to construct:
//
// 	https://www.wikidata.org/w/api.php?format=json&action=wbgetentities&ids=Q27229608
//
// We'll also use some of these values to build a permalink which looks
// as follows:
//
//	https://www.wikidata.org/w/index.php?title=Q178051&oldid=1301912874&format=json
//

const agent string = "wikiprov/0.0.1 (https://github.com/ross-spencer/wikiprov/; all.along.the.watchtower+github@gmail.com)"

var wikibaseAPI = "https://www.wikidata.org/w/api.php"
var wikidataBase = "https://www.wikidata.org/w/index.php"

var format = "json"
var action = "wbgetentities"
var props = "info"
