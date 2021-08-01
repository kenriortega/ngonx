package genkey

import (
	"testing"
)

func Test_genkey(t *testing.T) {
	word := "1q2w3e4r5t"
	apiKey := ApiKeyGenerator(word)
	if apiKey != "28f0116ef42bf718324946f13d787a1d41274a08335d52ee833d5b577f02a32a" {
		t.Error("Expected that ApiKeyGenerator return this hash [28f0116ef42bf718324946f13d787a1d41274a08335d52ee833d5b577f02a32a]")
	}
}
