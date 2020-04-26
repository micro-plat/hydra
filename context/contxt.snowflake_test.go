package context

import (
	"fmt"
	"testing"
)

func TestShowflake(t *testing.T) {
	wk := NewWorker(1)
	fmt.Println(wk.GetID())
}
