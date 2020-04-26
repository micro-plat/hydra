package context

import (
	"testing"
)

type input struct {
	dct    int
	ect    int
	ivalue interface{}
	evalue interface{}
}

func TestJSONResponse(t *testing.T) {
	// res := []input{
	// 	input{ect: CT_JSON, ivalue: "success", evalue: map[string]interface{}{"data": "success"}},
	// 	input{dct: CT_JSON, ect: CT_JSON, ivalue: "success", evalue: map[string]interface{}{"data": "success"}},
	// 	input{dct: CT_PLAIN, ect: CT_PLAIN, ivalue: "success", evalue: "success"},
	// 	input{dct: CT_PLAIN, ect: CT_PLAIN, ivalue: `{"data":"success"}`, evalue: `{"data":"success"}`},
	// 	input{dct: CT_XML, ect: CT_XML, ivalue: "success", evalue: "success"},
	// 	input{ect: CT_XML, ivalue: "<?xml><root></root>", evalue: []byte(`<?xml><root></root>`)},
	// 	input{ect: CT_HTML, ivalue: "<!DOCTYPE html><html></html>", evalue: `<!DOCTYPE html><html></html>`},
	// 	input{ect: CT_JSON, ivalue: map[string]interface{}{"data": "success"}, evalue: map[string]interface{}{"data": "success"}},
	// 	input{ect: CT_JSON, ivalue: 123, evalue: map[string]interface{}{"data": 123}},
	// }
	// for _, r := range res {
	// 	checkJSON(t, &Response{
	// 		Content: r.ivalue,
	// 		Params: map[string]interface{}{
	// 			"Content-Type": ContentTypes[r.dct],
	// 		},
	// 	}, r.ect, r.evalue)
	// }
}
