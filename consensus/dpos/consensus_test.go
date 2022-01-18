

package dpos

import (
	"fmt"
	"testing"

	"github.com/gcchains/chain/types"
	"github.com/stretchr/testify/assert"
)

func TestDpos_VerifyHeader(t *testing.T) {

	tests := []struct {
		name          string
		verifySuccess bool
		wantErr       bool
	}{
		{"verifyHeader success", true, false},
		{"verifyHeader failed", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Dpos{
				dh: &fakeDposHelper{verifySuccess: tt.verifySuccess},
			}

			err := c.VerifyHeader(&FakeReader{}, newHeader(), true, newHeader())
			fmt.Println("err:", err)
			if err := c.VerifyHeader(&FakeReader{}, newHeader(), true, newHeader()); (err != nil) != tt.wantErr {
				t.Errorf("Dpos.VerifyHeaders() got = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDpos_VerifyHeaders(t *testing.T) {
	tests := []struct {
		name          string
		verifySuccess bool
		wantErr       bool
	}{
		{"verifyHeader success", true, true},
		{"verifyHeader failed", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Dpos{
				dh: &fakeDposHelper{verifySuccess: tt.verifySuccess},
			}
			_, results := c.VerifyHeaders(
				&FakeReader{},
				[]*types.Header{newHeader()},
				[]bool{true},
				[]*types.Header{newHeader()})

			got := <-results
			fmt.Println("got:", got)
			if tt.wantErr != (got == nil) {
				t.Errorf("Dpos.VerifyHeaders() got = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestDpos_APIs(t *testing.T) {
	c := &Dpos{
		dh: &fakeDposHelper{},
	}
	got := c.APIs(nil)
	assert.Equal(t, 1, len(got), "only 1 api should be created")
}

// ========================================================
