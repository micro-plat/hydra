package pub

import (
	"io/ioutil"
	"log"
)

func getScript() (string, string) {
	bytes, err := ioutil.ReadFile("pub.sh")
	if err != nil {
		log.Fatal(err)
	}
	return "pub.sh", string(bytes)
}
