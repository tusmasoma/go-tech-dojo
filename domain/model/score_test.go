package model

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestModel_NewScore(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()

	patterns := []struct {
		name string
		arg  struct {
			userID string
			value  int
		}
		want struct {
			score *Score
			err   error
		}
	}{
		{
			name: "success",
			arg: struct {
				userID string
				value  int
			}{
				userID: userID,
				value:  100,
			},
			want: struct {
				score *Score
				err   error
			}{
				score: &Score{
					UserID: userID,
					Value:  100,
				},
				err: nil,
			},
		},
		{
			name: "Fail: userID is required",
			arg: struct {
				userID string
				value  int
			}{
				userID: "",
				value:  100,
			},
			want: struct {
				score *Score
				err   error
			}{
				score: nil,
				err:   fmt.Errorf("userID is required"),
			},
		},
		{
			name: "Fail: value is less than 0",
			arg: struct {
				userID string
				value  int
			}{
				userID: userID,
				value:  -1,
			},
			want: struct {
				score *Score
				err   error
			}{
				score: nil,
				err:   fmt.Errorf("value is less than 0"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score, err := NewScore(tt.arg.userID, tt.arg.value)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewScore() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewScore() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(score, tt.want.score, cmpopts.IgnoreFields(Score{}, "ID")); len(d) != 0 {
				t.Errorf("NewScore() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
