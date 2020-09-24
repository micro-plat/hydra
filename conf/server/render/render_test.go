package render

import (
	"os"
	"testing"
	"text/template"
	"time"
)

type TFuncs map[string]interface{}

type rspns struct {
	Status      int
	ContentType string
	Content     interface{}
}

//
func TestGet(t *testing.T) {
	// renderObj := NewRender()
	// renderObj.Get("/taosy/test",)

	tmplFuncs := TFuncs{}
	tmplFuncs["getString"] = func(name string) string {
		return name
	}
	tmplFuncs["getRequestID"] = func() string {
		return time.Now().Format("20060102150405")
	}
	tm, err := translate("/taosy/test", `{{.Status}}
	{{.Content}}
	{{getString .Content}}`, tmplFuncs, rspns{Status: 200, ContentType: "text/xml", Content: "success"})
	if err != nil {
		t.Errorf("getTmplt err:%+v", err)
		return
	}
	t.Errorf("result:%s", tm)

	ttype, err := translate("/taosy/test", "{{getString}}", tmplFuncs, rspns{Status: 200, ContentType: "text/xml", Content: "success"})
	if err != nil {
		t.Errorf("getTmplt err:%+v", err)
		return
	}
	t.Errorf("ttype:%s", ttype)

	tc, err := translate("/taosy/test", "fail", tmplFuncs, rspns{Status: 200, ContentType: "text/xml", Content: "success"})
	if err != nil {
		t.Errorf("getTmplt err:%+v", err)
		return
	}
	t.Errorf("tc:%s", tc)

}

func TestFuncMap(T *testing.T) {
	//os.Chdir("controller")
	t := template.New("render_test.html")
	tmplFuncs := template.FuncMap{}
	tmplFuncs["getString"] = func(name string) string {
		return name
	}
	t.Funcs(tmplFuncs)
	t, err := t.Parse(`{{.name}}
	{{.homeworkid}}
	{{getString .homeworkid}}`)
	if err != nil {
		panic(err)
	}
	render := make(map[string]interface{})
	render["homeworkid"] = time.Now().Format("20060102150405")
	render["name"] = "title man"
	// buff := bytes.NewBufferString("")
	if err = t.Execute(os.Stdout, render); err != nil {
		panic(err)
	}
}
