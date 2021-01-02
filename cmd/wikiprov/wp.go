package main

import (
	"fmt"

	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

func main() {
	var defaultIDs = "Q27229608"
	res, _ := wikiprov.GetWikidataProvenance(defaultIDs)
	fmt.Printf("%+v\n", res)
}
