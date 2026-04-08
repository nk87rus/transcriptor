package sbertoken

import (
	"fmt"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	at := Token{authKey: "???"}
	resultData, resultErr := at.Fetch(t.Context())
	fmt.Printf("RESP: %v\nERR: %v\n\n", resultData, resultErr)
	println(time.UnixMilli(resultData.ExpiresAt).String())
}
