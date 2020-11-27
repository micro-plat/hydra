package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
)

func main() {

	fmt.Println()
	fmt.Println(os.Getegid())
	fmt.Println(os.Getppid())

	fmt.Println("----")
	chanxx := make(chan os.Signal)
	signal.Notify(chanxx, os.Interrupt)

	<-chanxx

	r := `{"a":"b"}`
	result := map[string]string{}
	json.Unmarshal([]byte(r), &result)
	fmt.Println(result)

	o := XX{}

	v := tttt(o)

	fmt.Println("0", o["a"])
	fmt.Println("v", v["a"])

	t := reflect.TypeOf(v)

	fmt.Println("1", reflect.ValueOf(t).Interface())
	fmt.Println("2", t.Name())
	fmt.Println("3", t.Kind().String())

	fmt.Println("4", v)

	fv := fs()

	nfv := fv.(XX)

	fmt.Println("fv", fv)
	fmt.Println("nfv", nfv)

}

func tttt(x XX) XX {
	x["a"] = "1"
	return x
}

func ff() IF {

	return nil
}

func fs() interface{} {
	return ifs()
}

func ifs() XX {
	return nil
}

type XX map[string]string

type IF interface {
	F() string
}

type S struct {
	Name string
}
