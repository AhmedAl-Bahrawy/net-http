package restapi27oct_test

import (
	restapi27oct "net-http/restAPI27-OCT"
	"testing"
)

func TestCreateRestAPIServer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restapi27oct.CreateRestAPIServer()
		})
	}
}
