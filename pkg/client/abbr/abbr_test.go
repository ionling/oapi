package abbr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAbbrs(t *testing.T) {
	var (
		uid = "your-uid"
		tid = "your-tid"
		ctx = context.Background()
	)
	t.Run("no-parent-category", func(t *testing.T) {
		res, err := GetAbbrs(ctx, &GetAbbrsReq{
			UID:        uid,
			TokenID:    tid,
			Term:       "example",
			SearchType: STReverse,
		})

		require.NoError(t, err)
		require.Len(t, res.Result, 3)
		assert.Equal(t, &Term{
			commonTerm: commonTerm{
				ID:           "1440380",
				Term:         "E.G",
				Definition:   "Example",
				Category:     "MISCELLANEOUS",
				CategoryName: "Miscellaneous",
				Score:        "3.17",
			},
		}, res.Result[2])
	})

	t.Run("exact-type", func(t *testing.T) {
		res, err := GetAbbrs(ctx, &GetAbbrsReq{
			UID:        uid,
			TokenID:    tid,
			Term:       "accept",
			SearchType: STExact,
		})
		require.NoError(t, err)
		require.Len(t, res.Result, 12)
		assert.Equal(t, &Term{
			commonTerm: commonTerm{
				ID:           "2153349",
				Term:         "ACCEPT",
				Definition:   "Admissions Community Cultivating Equity Peace Today",
				Category:     "COMMUNITY",
				CategoryName: "Community",
				Score:        "1.00",
			},
		}, res.Result[1])
	})

	t.Run("single-result", func(t *testing.T) {
		res, err := GetAbbrs(ctx, &GetAbbrsReq{
			UID:        uid,
			TokenID:    tid,
			Term:       "accept",
			SearchType: STReverse,
		})
		require.NoError(t, err)
		require.Len(t, res.Result, 1)
		assert.Equal(t, &Term{
			commonTerm: commonTerm{
				ID:           "1419160",
				Term:         "A",
				Definition:   "Accept",
				Category:     "USGOV",
				CategoryName: "US Government",
				Score:        "1.83",
			},
			ParentCategory:     "GOVERNMENTAL",
			ParentCategoryName: "Governmental",
		}, res.Result[0])
	})
}
