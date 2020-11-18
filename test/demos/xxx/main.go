package main

import (
	"fmt"
	"reflect"
)

func main() {

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
