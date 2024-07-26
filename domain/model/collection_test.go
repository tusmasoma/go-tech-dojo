package model

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestModel_NewCollection(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			name   string
			rarity int
			weight int
		}
		want struct {
			collection *Collection
			err        error
		}
	}{
		{
			name: "success",
			arg: struct {
				name   string
				rarity int
				weight int
			}{
				name:   "collection",
				rarity: 3,
				weight: 10,
			},
			want: struct {
				collection *Collection
				err        error
			}{
				collection: &Collection{
					Name:   "collection",
					Rarity: 3,
					Weight: 10,
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is required",
			arg: struct {
				name   string
				rarity int
				weight int
			}{
				name:   "",
				rarity: 3,
				weight: 10,
			},
			want: struct {
				collection *Collection
				err        error
			}{
				collection: nil,
				err:        fmt.Errorf("name is empty"),
			},
		},
		{
			name: "Fail: rarity is invalid",
			arg: struct {
				name   string
				rarity int
				weight int
			}{
				name:   "collection",
				rarity: 6,
				weight: 10,
			},
			want: struct {
				collection *Collection
				err        error
			}{
				collection: nil,
				err:        fmt.Errorf("rarity is invalid"),
			},
		},
		{
			name: "Fail: weight is invalid",
			arg: struct {
				name   string
				rarity int
				weight int
			}{
				name:   "collection",
				rarity: 3,
				weight: -1,
			},
			want: struct {
				collection *Collection
				err        error
			}{
				collection: nil,
				err:        fmt.Errorf("weight is invalid"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getCollection, err := NewCollection(tt.arg.name, tt.arg.rarity, tt.arg.weight)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewCollection() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewCollection() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(getCollection, tt.want.collection, cmpopts.IgnoreFields(Collection{}, "ID")); len(d) != 0 {
				t.Errorf("NewCollection() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
