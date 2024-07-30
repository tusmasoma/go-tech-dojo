package model

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModel_NewRanking(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			userName string
			rank     int
			score    int
		}
		want struct {
			ranking *Ranking
			err     error
		}
	}{
		{
			name: "success",
			arg: struct {
				userName string
				rank     int
				score    int
			}{
				userName: "user",
				rank:     1,
				score:    100,
			},
			want: struct {
				ranking *Ranking
				err     error
			}{
				ranking: &Ranking{
					UserName: "user",
					Rank:     1,
					Score:    100,
				},
				err: nil,
			},
		},
		{
			name: "Fail: userName is required",
			arg: struct {
				userName string
				rank     int
				score    int
			}{
				userName: "",
				rank:     1,
				score:    100,
			},
			want: struct {
				ranking *Ranking
				err     error
			}{
				ranking: nil,
				err:     fmt.Errorf("userName is empty"),
			},
		},
		{
			name: "Fail: rank is less than 1",
			arg: struct {
				userName string
				rank     int
				score    int
			}{
				userName: "user",
				rank:     0,
				score:    100,
			},
			want: struct {
				ranking *Ranking
				err     error
			}{
				ranking: nil,
				err:     fmt.Errorf("rank is less than 1"),
			},
		},
		{
			name: "Fail: score is less than 0",
			arg: struct {
				userName string
				rank     int
				score    int
			}{
				userName: "user",
				rank:     1,
				score:    -1,
			},
			want: struct {
				ranking *Ranking
				err     error
			}{
				ranking: nil,
				err:     fmt.Errorf("score is less than 0"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ranking, err := NewRanking(tt.arg.userName, tt.arg.rank, tt.arg.score)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewRanking() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewRanking() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(ranking, tt.want.ranking); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
