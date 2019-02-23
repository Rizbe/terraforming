package gen

import (
	"log"

	h "github.com/rodaine/hclencoder"
)

//GenHCL test
func GenHCL(resource interface{}) (string, error) {
	hcl, err := h.Encode(resource)
	if err != nil {
		log.Fatal("unable to encode: ", err)
	}

	// fmt.Print(string(hcl))
	return string(hcl), nil

}
