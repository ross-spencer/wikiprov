# wikiprov

Rudimentary reification in Wikidata. This package takes the spargo generic
SPARQL handling package I created and wraps Wikibase provenance around it.

## Status and documentation

<!--markdownlint-disable-->

| report card | wikiprov | spargo |
|-------------|----------|--------|
| [![Go Report Card][report-badge]][report-card] | [![GoDoc][wikiprov-doc-badge]][godoc-1] | [![GoDoc][spargo-doc-badge]][godoc-2] |

<!--markdownlint-enable-->

[report-badge]:
    https://goreportcard.com/badge/github.com/ross-spencer/wikiprov
[report-card]:
    https://goreportcard.com/report/github.com/ross-spencer/wikiprov

[wikiprov-doc-badge]:
    https://godoc.org/github.com/ross-spencer/wikiprov?status.svg

[spargo-doc-badge]:
    https://godoc.org/github.com/ross-spencer/spargo?status.svg

[godoc-1]:
    https://pkg.go.dev/github.com/ross-spencer/wikiprov/pkg/wikiprov

[godoc-2]:
    https://godoc.org/github.com/ross-spencer/wikiprov/pkg/spargo

-----

## Introduction

Where the generic SPARQL results from any service look as follows:

```json
{
  "head": {},
  "results": {
    "bindings": [{}]
  }
}
```

With Wikidata sitting on-top of a Wikibase instance, it allows us to try and
retrieve some amount of provenance for the IRIs returned. While not 'pure'
linked data as we might like we can make the best of what we've got and return
it with our query results anyway. Those results take the form:

```json
{
  "head": {},
  "results": {
    "bindings": [{}]
  },
  "provenance": {}
}
```

The `Provenance` block comes from the wikiprov package and (per unique QID)
tentatively looks as follows:

<!--markdownlint-disable-->

```json
{
  "Title": "Q5381415",
  "Entity": "http://wikidata.org/entity/Q5381415",
  "Revision": 1343296571,
  "Modified": "2021-01-18T05:36:32Z",
  "Permalink": "https://www.wikidata.org/w/index.php?format=json\u0026oldid=1343296571\u0026title=Q5381415",
  "History": [
    "2021-01-18T05:36:32Z (oldid: 1343296571): 'Lockal' edited: '/* wbcreateclaim-create:1| */ [[Property:P646]]: /m/0fc557'",
    "2020-08-04T23:41:27Z (oldid: 1247209137): 'Beet keeper' edited: '/* wbsetclaim-update:2||1 */ [[Property:P4152]]: B297E169'",
    "2020-08-04T23:40:10Z (oldid: 1247208427): 'Beet keeper' edited: '/* wbsetclaim-update:2||1 */ [[Property:P4152]]: 325E1010'",
    "2020-02-21T14:40:33Z (oldid: 1120067133): 'YULdigitalpreservation' edited: '/* wbsetaliases-add:3|en */ Envoy Document File, Envoy Document, Envoy 1'",
    "2020-02-21T14:38:57Z (oldid: 1120066909): 'YULdigitalpreservation' edited: '/* wbsetclaim-create:2||1 */ [[Property:P348]]: 1'"
  ]
}
```

<!--markdownlint-enable-->

This enables users to look up a QID and see what last happened to that record
from the same SPARQL results source.

## Spargo package

It is anticipated wikiprov will be used primarily as a golang package.

A basic example might look as follows:

<!--markdownlint-disable-->

```go
package main

import (
	"fmt"

	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

func main() {
	var qid = "Q5381415"
	res, err := wikiprov.GetWikidataProvenance(qid, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

<!--markdownlint-enable-->

The results will look similar to the following:

<!--markdownlint-disable-->

```json
{
  "Title": "Q5381415",
  "Revision": 2036898689,
  "Modified": "2023-12-25T21:34:28Z",
  "Permalink": "https://www.wikidata.org/w/index.php?oldid=2036898689&title=Q5381415",
  "History": [
    "2023-12-25T21:34:28Z (oldid: 2036898689): 'Dragomouse' edited: '/* wbsetclaim-create:2||1 */ [[Property:P4839]]: Entity[\"FileFormat\", \"EVY-1\"]'",
    "2023-04-13T12:30:22Z (oldid: 1874135847): 'Maqivi' edited: '/* wbsetlabel-add:1|ru */ Envoy'",
    "2022-09-20T07:09:09Z (oldid: 1732599165): 'LogainmBot' edited: '/* wbeditentity-update-languages-short:0||ga */ Irish label added'",
    "2021-05-18T09:40:10Z (oldid: 1423352886): 'Edoderoobot' edited: '/* wbeditentity-update-languages-short:0||nl */ nl-description, [[User:Edoderoobot/Set-nl-description|python code]] - fileformat'",
    "2021-01-18T05:36:32Z (oldid: 1343296571): 'Lockal' edited: '/* wbcreateclaim-create:1| */ [[Property:P646]]: /m/0fc557'",
    "2020-08-04T23:41:27Z (oldid: 1247209137): 'Beet keeper' edited: '/* wbsetclaim-update:2||1 */ [[Property:P4152]]: B297E169'",
    "2020-08-04T23:40:10Z (oldid: 1247208427): 'Beet keeper' edited: '/* wbsetclaim-update:2||1 */ [[Property:P4152]]: 325E1010'",
    "2020-02-21T14:40:33Z (oldid: 1120067133): 'YULdigitalpreservation' edited: '/* wbsetaliases-add:3|en */ Envoy Document File, Envoy Document, Envoy 1'",
    "2020-02-21T14:38:57Z (oldid: 1120066909): 'YULdigitalpreservation' edited: '/* wbsetclaim-create:2||1 */ [[Property:P348]]: 1'",
    "2020-02-21T14:38:44Z (oldid: 1120066880): 'YULdigitalpreservation' edited: '/* wbsetclaim-create:2||1 */ [[Property:P2748]]: fmt/1286'"
  ]
}
```

<!--markdownlint-enable-->

Check out the godoc linked to at the top of this README for more info.

## Command line

Wikiprov can also be used on the command line through its partner app `spargo`.

### Spargo

example:

```text
cat examples/001-default-wikidata.sparql | ./spargo
```

### Spargo as interpreter

Given a correctly formatted [shebang][shebang-1], e.g.

```text
#!/usr/bin/spargo
```

alternatively:

```text
#!/your/own/path/to/spargo
```

You can write SPARQL into a text file and it can be run from the command line
as if it is an executable like any other on Linux. The output will be piped to
stdout.

E.g. a script called `001-default-wikidata.sparql` might look as follows:

<!--markdownlint-disable-->

```text
#!/usr/bin/spargo

ENDPOINT=https://query.wikidata.org/sparql
WIKIBASEURL=https://www.wikidata.org/
HISTORY=5
SUBJECTPARAM=?item

# Default query example on Wikidata:
SELECT ?item ?itemLabel
WHERE
{
  ?item wdt:P31 wd:Q146.
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en". }
}
limit 1
```

And its output, `./001-default-wikidata.sparql`:

```json
{
  "head": {
    "vars": [
      "item",
      "itemLabel"
    ]
  },
  "results": {
    "bindings": [
      {
        "item": {
          "type": "uri",
          "value": "http://www.wikidata.org/entity/Q378619"
        },
        "itemLabel": {
          "xml:lang": "en",
          "type": "literal",
          "value": "CC"
        }
      }
    ]
  },
  "provenance": [
    {
      "Title": "Q378619",
      "Entity": "http://wikidata.org/entity/Q378619",
      "Revision": 2252118087,
      "Modified": "2024-09-23T16:58:53Z",
      "Permalink": "https://www.wikidata.org/w/index.php?oldid=2252118087&title=Q378619",
      "History": [
        "2024-09-23T16:58:53Z (oldid: 2252118087): 'Skouratov' edited: '/* undo:0||2235229659|189.214.7.137 */'",
        "2024-08-24T18:50:49Z (oldid: 2235229659): '189.214.7.137' edited: '/* wbsetclaim-create:2||1 */ [[Property:P31]]: [[Q5]]'",
        "2024-07-09T16:58:39Z (oldid: 2199528978): 'Marek Mazurkiewicz' edited: '/* wbsetdescription-add:1|eo */ kato'",
        "2024-07-09T16:58:38Z (oldid: 2199528968): 'Marek Mazurkiewicz' edited: '/* wbsetlabel-add:1|eo */ CC'",
        "2024-04-22T15:54:53Z (oldid: 2134899800): 'MatSuBot' edited: '/* wbeditentity-update-languages-short:0||tw */ add missing labels'"
      ]
    }
  ]
}
```

If the script is called `query.sparql` and has been given executable
permissions `chmod +X query.sparql`, it can be run using:

```sh
./query.sparql
```

The output format is `json` and so can be interpreted using tools such as
[`jq`][jq-1].

[jq-1]: https://jqlang.github.io/jq/

<!--markdownlint-enable-->

[shebang-1]: https://en.wikipedia.org/wiki/Shebang_(Unix)

### File format

Two structures are compatible with `spargo` depending on whether you'd like to
return provenance.

### with provenance

A basic structure needs to list the following:

* wikidata query service endpoint,
* wikibase url,
* history length,
* parameter to return history for, e.g. which ?subject, ?predicate,
or '?object`.

```text
#!/usr/bin/spargo

ENDPOINT=...
WIKIBASEURL=...
HISTORY=...
SUBJECTPARAM=...

# Comment
{sparql query}
```

> NB. Examples can be found in [cmd][prov-examples].

[prov-examples]: cmd/spargo/prov-examples/

### without provenance

Without provenance you merely need the wikidata query service endpoint.

```text
#!/usr/bin/spargo

ENDPOINT=...

# Comment
{sparql query}
```

> NB. Examples can be found in [cmd][examples].

[examples]: cmd/spargo/examples/

## Architectural decisions

### Errors returned from SPARQLWithProv

I have flip-flopped with the error handling and currently try to be as
flexible as possible. If there are some errors retrieving provenance then the
provenance structure is still returned and gaps need to be identified upstream.
Callers are also able to request no-provenance which will of course render an
empty structure. As an experimental library it makes sense to be flexible while
the concept is proven.

## Feedback

Please leave an issue you have questions or want to develop this library
further.

## License

Apache License 2.0. More info [here](LICENSE).
