package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetric(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	var funcMetric *Metric
	var stdMetric *Metric

	t.Run("create func metric", func(t *testing.T) {
		funcMetric = CreateFuncMetric("name", "category", "unit", func() float64 {
			return 42
		})
	})

	t.Run("get func value", func(t *testing.T) {
		val := funcMetric.Get()
		require.EqualValues(42, val)
	})

	t.Run("create std metric", func(t *testing.T) {
		stdMetric = CreateMetric("name2", "category2", "unit2")
		assert.NotNil(stdMetric)
	})

	t.Run("set value", func(t *testing.T) {
		stdMetric.Set(20000)
	})

	t.Run("get value", func(t *testing.T) {
		val := stdMetric.Get()
		require.EqualValues(20000, val)
	})
}
