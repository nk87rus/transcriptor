package gigachat

import (
	"fmt"
	"testing"
)

const ak = "SECRET"

func TestChat(t *testing.T) {
	g, ge := Init(t.Context(), ak)
	if ge != nil {
		t.Fatal(ge)
	}

	resultData, resultError := g.Chat(t.Context(), "что такое свет?")
	if resultError != nil {
		t.Fatal(resultError)
	}

	fmt.Printf("---\n%s\n---", resultData)

}
