# wikiprov

Rudimentary reification in Wikidata. This package takes the spargo generic
SPARQL handling package I created and wraps Wikibase provenance around it.

Where the generic SPARQL results from any service look as follows:

```json
{
	"Head": {},
	"Bindings": {}
}
```
With Wikidata sitting on-top of a Wikibase instance, it allows us to try and
retrieve some amount of provenance for the IRIS returned. While not pure linked
data as we might like we can make the best of what we've got and return it with
our query results anyway. Those results take the form:

```json
{
	"Head": {},
	"Bindings": {},
	"Provenance": {}
}
```
The `Provenance` block comes from the wikiprov package and (per unique QID)
tentatively looks as follows:

```json
{
  "Title": "Q5381415",
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
This enables users to look up a QID and see what last happened to that record
from the same SPARQL results source.
