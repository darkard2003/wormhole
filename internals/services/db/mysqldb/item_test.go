package mysqldb_test

import (
	"testing"

	"github.com/darkard2003/wormhole/internals/services/db/mysqldb"
)

func TestMySqlRepo_PopLatestItem(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		channelId int
		want      any
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var r mysqldb.MySqlRepo
			got, gotErr := r.PopLatestItem(tt.channelId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("PopLatestItem() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("PopLatestItem() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("PopLatestItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
