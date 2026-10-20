// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregateutil // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/aggregateutil"

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func Test_CopyMetricDetails(t *testing.T) {
	gaugeFunc := func() pmetric.Metric {
		m := pmetric.NewMetric()
		m.SetDescription("desc")
		m.SetName("name")
		m.SetUnit("unit")
		m.SetEmptyGauge()
		return m
	}

	sumFunc := func() pmetric.Metric {
		m := pmetric.NewMetric()
		m.SetDescription("desc")
		m.SetName("name")
		m.SetUnit("unit")
		s := m.SetEmptySum()
		s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
		s.SetIsMonotonic(true)
		return m
	}

	summaryFunc := func() pmetric.Metric {
		m := pmetric.NewMetric()
		m.SetDescription("desc")
		m.SetName("name")
		m.SetUnit("unit")
		m.SetEmptySummary()
		return m
	}

	histogramFunc := func() pmetric.Metric {
		m := pmetric.NewMetric()
		m.SetDescription("desc")
		m.SetName("name")
		m.SetUnit("unit")
		m.SetEmptyHistogram().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
		return m
	}

	expHistogramFunc := func() pmetric.Metric {
		m := pmetric.NewMetric()
		m.SetDescription("desc")
		m.SetName("name")
		m.SetUnit("unit")
		m.SetEmptyExponentialHistogram().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
		return m
	}
	tests := []struct {
		name string
		from func() pmetric.Metric
		to   func() pmetric.Metric
	}{
		{
			name: "gauge",
			from: gaugeFunc,
			to:   gaugeFunc,
		},
		{
			name: "summary",
			from: summaryFunc,
			to:   summaryFunc,
		},
		{
			name: "sum",
			from: sumFunc,
			to:   sumFunc,
		},
		{
			name: "histogram",
			from: histogramFunc,
			to:   histogramFunc,
		},
		{
			name: " exp histogram",
			from: expHistogramFunc,
			to:   expHistogramFunc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pmetric.NewMetric()
			from := tt.from()
			to := tt.to()
			CopyMetricDetails(from, result)
			require.Equal(t, to, result)
		})
	}
}

func Test_FilterAttributes(t *testing.T) {
	tests := []struct {
		name string
		attr []string
		want func() pmetric.Metric
	}{
		{
			name: "nil",
			attr: nil,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
		},
		{
			name: "empty",
			attr: []string{},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				s.DataPoints().AppendEmpty()
				return m
			},
		},
		{
			name: "valid",
			attr: []string{"attr1"},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := pmetric.NewMetric()
			s := m.SetEmptySum()
			d := s.DataPoints().AppendEmpty()
			d.Attributes().PutStr("attr1", "val1")
			d.Attributes().PutStr("attr2", "val2")

			FilterAttrs(m, tt.attr)
			require.Equal(t, tt.want(), m)
		})
	}
}

func Test_RangeDataPointAttributes(t *testing.T) {
	fun := func(attrs pcommon.Map) bool {
		attrs.RemoveIf(func(k string, _ pcommon.Value) bool {
			return !slices.Contains([]string{"attr1"}, k)
		})
		return true
	}

	tests := []struct {
		name string
		in   func() pmetric.Metric
		want func() pmetric.Metric
	}{
		{
			name: "sum",
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
		{
			name: "gauge",
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyGauge()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyGauge()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
		{
			name: "summary",
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySummary()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySummary()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
		{
			name: "histogram",
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
		{
			name: "exp histogram",
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyExponentialHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.Attributes().PutStr("attr2", "val2")
				return m
			},
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyExponentialHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.in()
			RangeDataPointAttributes(m, fun)
			require.Equal(t, tt.want(), m)
		})
	}
}

func Test_GroupDataPoints(t *testing.T) {
	mapAttr := pcommon.NewMap()
	mapAttr.PutStr("attr1", "val1")
	hash := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}))

	hashHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	hashExpHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	tests := []struct {
		name     string
		in       func() pmetric.Metric
		aggGroup AggGroups
		want     AggGroups
	}{
		{
			name: "sum",
			aggGroup: AggGroups{
				sum: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumber(),
				},
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				d := s.DataPoints().AppendEmpty()
				d.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
				d.Attributes().PutStr("attr1", "val1")
				d.SetIntValue(5)
				return m
			},
			want: AggGroups{
				sum: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumberDouble(),
				},
			},
		},
		{
			name: "gauge",
			aggGroup: AggGroups{
				gauge: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumber(),
				},
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyGauge()
				d := s.DataPoints().AppendEmpty()
				d.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
				d.Attributes().PutStr("attr1", "val1")
				d.SetIntValue(5)
				return m
			},
			want: AggGroups{
				gauge: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumberDouble(),
				},
			},
		},
		{
			name: "histogram",
			aggGroup: AggGroups{
				histogram: map[string]pmetric.HistogramDataPointSlice{
					hashHistogram: testDataHistogram(),
				},
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyHistogram()
				s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				d := s.DataPoints().AppendEmpty()
				d.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
				d.Attributes().PutStr("attr1", "val1")
				d.SetCount(1)
				return m
			},
			want: AggGroups{
				histogram: map[string]pmetric.HistogramDataPointSlice{
					hashHistogram: testDataHistogramDouble(),
				},
			},
		},
		{
			name: "exp histogram",
			aggGroup: AggGroups{
				expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{
					hashExpHistogram: testDataExpHistogram(),
				},
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyExponentialHistogram()
				s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				d := s.DataPoints().AppendEmpty()
				d.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
				d.Attributes().PutStr("attr1", "val1")
				d.SetCount(9)
				d.SetZeroCount(2)
				d.Positive().BucketCounts().Append(0, 1, 2, 3)
				d.Negative().BucketCounts().Append(0, 1)
				return m
			},
			want: AggGroups{
				expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{
					hashExpHistogram: testDataExpHistogramDouble(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.aggGroup
			GroupDataPoints(tt.in(), &a)
			require.Equal(t, tt.want, a)
		})
	}
}

func Test_MergeDataPoints(t *testing.T) {
	mapAttr := pcommon.NewMap()
	mapAttr.PutStr("attr1", "val1")

	hash := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}))

	hashHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	hashExpHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	tests := []struct {
		name     string
		typ      AggregationType
		aggGroup AggGroups
		want     func() pmetric.Metric
		in       func() pmetric.Metric
	}{
		{
			name: "sum",
			aggGroup: AggGroups{
				sum: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumberDouble(),
				},
			},
			typ: Sum,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.SetIntValue(6)
				return m
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptySum()
				s.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				return m
			},
		},
		{
			name: "gauge",
			aggGroup: AggGroups{
				gauge: map[string]pmetric.NumberDataPointSlice{
					hash: testDataNumberDouble(),
				},
			},
			typ: Sum,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyGauge()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.SetIntValue(6)
				return m
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				m.SetEmptyGauge()
				return m
			},
		},
		{
			name: "histogram",
			aggGroup: AggGroups{
				histogram: map[string]pmetric.HistogramDataPointSlice{
					hashHistogram: testDataHistogramDouble(),
				},
			},
			typ: Sum,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.SetCount(3)
				return m
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				m.SetEmptyHistogram()
				return m
			},
		},
		{
			name: "exp histogram",
			aggGroup: AggGroups{
				expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{
					hashExpHistogram: testDataExpHistogramDouble(),
				},
			},
			typ: Sum,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyExponentialHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.SetCount(16)
				d.SetZeroCount(3)
				d.Positive().BucketCounts().Append(0, 2, 4, 3)
				d.Negative().BucketCounts().Append(0, 2, 2)
				return m
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				m.SetEmptyExponentialHistogram()
				return m
			},
		},
		{
			name: "exp histogram with different offsets",
			aggGroup: AggGroups{
				expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{
					hashExpHistogram: testDataExpHistogramWithDifferentOffsets(),
				},
			},
			typ: Sum,
			want: func() pmetric.Metric {
				m := pmetric.NewMetric()
				s := m.SetEmptyExponentialHistogram()
				d := s.DataPoints().AppendEmpty()
				d.Attributes().PutStr("attr1", "val1")
				d.SetCount(12)
				d.SetZeroCount(3)
				// First datapoint: positive offset 0, buckets [1, 2]
				// Second datapoint: positive offset 1, buckets [3, 4] (represents [0, 3, 4])
				// Third datapoint: positive offset 5, buckets [1] (represents [0, 0, 0, 0, 0, 1])
				// Result: [1+0+0, 2+3+0, 4+0+0, 0+0+0, 0+0+0, 0+0+1] = [1, 5, 4, 0, 0, 1]
				d.Positive().BucketCounts().Append(1, 5, 4, 0, 0, 1)
				d.Positive().SetOffset(0)

				// First data point: negative offset 3, buckets [1, 2] (represents [0, 0, 0, 1, 2])
				// Second data point: negative offset 1, buckets [5, 6] (represents [0, 5, 6])
				// Third data point: negative side empty
				// Result: [0+0, 0+5, 0+6, 1+0, 2+0] = [0, 5, 6, 1, 2] (with offset 1: [5, 6, 1, 2])
				d.Negative().BucketCounts().Append(5, 6, 1, 2)
				d.Negative().SetOffset(1)

				return m
			},
			in: func() pmetric.Metric {
				m := pmetric.NewMetric()
				m.SetEmptyExponentialHistogram()
				return m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.in()
			MergeDataPoints(m, tt.typ, tt.aggGroup)
			require.Equal(t, tt.want(), m)
		})
	}
}

// Test_MergeDataPoints_SumPresence is a regression test for #49379: merging
// histogram datapoints must add their sums when all of them carry one, and must
// leave the merged datapoint without a sum as soon as any contributor lacks one
// (rather than reporting a fabricated sum of 0).
func Test_MergeDataPoints_SumPresence(t *testing.T) {
	mapAttr := pcommon.NewMap()
	mapAttr.PutStr("attr1", "val1")
	hashHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)
	hashExpHistogram := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	histogramPair := func(setSum0, setSum1 bool) pmetric.HistogramDataPointSlice {
		s := pmetric.NewHistogramDataPointSlice()
		d0 := s.AppendEmpty()
		d0.Attributes().PutStr("attr1", "val1")
		d0.SetCount(2)
		if setSum0 {
			d0.SetSum(1.5)
		}
		d1 := s.AppendEmpty()
		d1.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
		d1.Attributes().PutStr("attr1", "val1")
		d1.SetCount(1)
		if setSum1 {
			d1.SetSum(2.5)
		}
		return s
	}
	expHistogramPair := func(setSum0, setSum1 bool) pmetric.ExponentialHistogramDataPointSlice {
		s := pmetric.NewExponentialHistogramDataPointSlice()
		d0 := s.AppendEmpty()
		d0.Attributes().PutStr("attr1", "val1")
		d0.SetCount(2)
		if setSum0 {
			d0.SetSum(1.5)
		}
		d1 := s.AppendEmpty()
		d1.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
		d1.Attributes().PutStr("attr1", "val1")
		d1.SetCount(1)
		if setSum1 {
			d1.SetSum(2.5)
		}
		return s
	}

	t.Run("histogram both have sum -> sum added", func(t *testing.T) {
		m := pmetric.NewMetric()
		m.SetEmptyHistogram()
		MergeDataPoints(m, Sum, AggGroups{histogram: map[string]pmetric.HistogramDataPointSlice{hashHistogram: histogramPair(true, true)}})
		d := m.Histogram().DataPoints().At(0)
		require.True(t, d.HasSum())
		require.Equal(t, 4.0, d.Sum())
	})
	t.Run("histogram one missing sum -> sum absent", func(t *testing.T) {
		m := pmetric.NewMetric()
		m.SetEmptyHistogram()
		MergeDataPoints(m, Sum, AggGroups{histogram: map[string]pmetric.HistogramDataPointSlice{hashHistogram: histogramPair(true, false)}})
		d := m.Histogram().DataPoints().At(0)
		require.False(t, d.HasSum())
	})
	t.Run("exp histogram both have sum -> sum added", func(t *testing.T) {
		m := pmetric.NewMetric()
		m.SetEmptyExponentialHistogram()
		MergeDataPoints(m, Sum, AggGroups{expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{hashExpHistogram: expHistogramPair(true, true)}})
		d := m.ExponentialHistogram().DataPoints().At(0)
		require.True(t, d.HasSum())
		require.Equal(t, 4.0, d.Sum())
	})
	t.Run("exp histogram one missing sum -> sum absent", func(t *testing.T) {
		m := pmetric.NewMetric()
		m.SetEmptyExponentialHistogram()
		MergeDataPoints(m, Sum, AggGroups{expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{hashExpHistogram: expHistogramPair(false, true)}})
		d := m.ExponentialHistogram().DataPoints().At(0)
		require.False(t, d.HasSum())
	})
}

func testDataNumber() pmetric.NumberDataPointSlice {
	data := pmetric.NewNumberDataPointSlice()
	d := data.AppendEmpty()
	d.Attributes().PutStr("attr1", "val1")
	d.SetIntValue(1)
	return data
}

func testDataNumberDouble() pmetric.NumberDataPointSlice {
	dataWant := pmetric.NewNumberDataPointSlice()
	dWant := dataWant.AppendEmpty()
	dWant.Attributes().PutStr("attr1", "val1")
	dWant.SetIntValue(1)
	dWant2 := dataWant.AppendEmpty()
	dWant2.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
	dWant2.Attributes().PutStr("attr1", "val1")
	dWant2.SetIntValue(5)
	return dataWant
}

func testDataHistogram() pmetric.HistogramDataPointSlice {
	data := pmetric.NewHistogramDataPointSlice()
	d := data.AppendEmpty()
	d.Attributes().PutStr("attr1", "val1")
	d.SetCount(2)
	return data
}

func testDataHistogramDouble() pmetric.HistogramDataPointSlice {
	dataWant := pmetric.NewHistogramDataPointSlice()
	dWant := dataWant.AppendEmpty()
	dWant.Attributes().PutStr("attr1", "val1")
	dWant.SetCount(2)
	dWant2 := dataWant.AppendEmpty()
	dWant2.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
	dWant2.Attributes().PutStr("attr1", "val1")
	dWant2.SetCount(1)
	return dataWant
}

func testDataExpHistogram() pmetric.ExponentialHistogramDataPointSlice {
	data := pmetric.NewExponentialHistogramDataPointSlice()
	d := data.AppendEmpty()
	d.Attributes().PutStr("attr1", "val1")
	d.SetCount(7)
	d.SetZeroCount(1)
	d.Positive().BucketCounts().Append(0, 1, 2)
	d.Negative().BucketCounts().Append(0, 1, 2)
	return data
}

func testDataExpHistogramDouble() pmetric.ExponentialHistogramDataPointSlice {
	dataWant := pmetric.NewExponentialHistogramDataPointSlice()

	dWant := dataWant.AppendEmpty()
	dWant.Attributes().PutStr("attr1", "val1")
	dWant.SetCount(7)
	dWant.SetZeroCount(1)
	dWant.Positive().BucketCounts().Append(0, 1, 2)
	dWant.Negative().BucketCounts().Append(0, 1, 2)

	dWant2 := dataWant.AppendEmpty()
	dWant2.SetTimestamp(pcommon.NewTimestampFromTime(time.Time{}))
	dWant2.Attributes().PutStr("attr1", "val1")
	dWant2.SetCount(9)
	dWant2.SetZeroCount(2)
	// Use a larger number of buckets than above to check that we expand the
	// destination array as needed while merging.
	dWant2.Positive().BucketCounts().Append(0, 1, 2, 3)
	// Use a smaller number of buckets than above to check that we merge values
	// into the correct, existing buckets.
	dWant2.Negative().BucketCounts().Append(0, 1)

	return dataWant
}

func testDataExpHistogramWithDifferentOffsets() pmetric.ExponentialHistogramDataPointSlice {
	dataWant := pmetric.NewExponentialHistogramDataPointSlice()

	dWant := dataWant.AppendEmpty()
	dWant.Attributes().PutStr("attr1", "val1")
	dWant.SetCount(6)
	dWant.SetZeroCount(2)
	dWant.Positive().SetOffset(0)
	dWant.Positive().BucketCounts().Append(1, 2) // [1, 2]
	dWant.Negative().SetOffset(3)
	dWant.Negative().BucketCounts().Append(1, 2) // [0, 0, 0, 1, 2]

	// Different offsets
	dWant2 := dataWant.AppendEmpty()
	dWant2.Attributes().PutStr("attr1", "val1")
	dWant2.SetCount(5)
	dWant2.SetZeroCount(1)
	dWant2.Positive().SetOffset(1)
	dWant2.Positive().BucketCounts().Append(3, 4) // [0, 3, 4]
	dWant2.Negative().SetOffset(1)
	dWant2.Negative().BucketCounts().Append(5, 6) // [0, 5, 6]

	// Set large offset for positive, no offset for negative
	dWant3 := dataWant.AppendEmpty()
	dWant3.Attributes().PutStr("attr1", "val1")
	dWant3.SetCount(1)
	dWant3.SetZeroCount(0)
	dWant3.Positive().SetOffset(5)
	dWant3.Positive().BucketCounts().Append(1) // [0, 0, 0, 0, 0, 1]

	return dataWant
}

func Test_mergeExponentialHistogramBuckets(t *testing.T) {
	type side struct {
		scale  int32
		offset int32
		counts []uint64
	}
	gap := func(n int, first, last uint64) []uint64 {
		c := make([]uint64, n)
		c[0], c[n-1] = first, last
		return c
	}
	tests := []struct {
		name string
		tgt  side
		src  side
		want side
	}{
		{
			name: "same offset",
			tgt:  side{0, 3, []uint64{1, 2}}, src: side{0, 3, []uint64{3, 4, 5}},
			want: side{0, 3, []uint64{4, 6, 5}},
		},
		{
			name: "source offset higher",
			tgt:  side{0, 0, []uint64{1, 2}}, src: side{0, 1, []uint64{3, 4}},
			want: side{0, 0, []uint64{1, 5, 4}},
		},
		{
			name: "source offset lower with a gap",
			tgt:  side{0, 5, []uint64{10}}, src: side{0, 3, []uint64{1}},
			want: side{0, 3, []uint64{1, 0, 10}},
		},
		{
			name: "source offset lower and overlapping",
			tgt:  side{0, 5, []uint64{10, 20}}, src: side{0, 4, []uint64{1, 2}},
			want: side{0, 4, []uint64{1, 12, 20}},
		},
		{
			name: "source offset lower with zero buckets",
			tgt:  side{0, 5, []uint64{10}}, src: side{0, 3, []uint64{1, 0, 2}},
			want: side{0, 3, []uint64{1, 0, 12}},
		},
		{
			name: "leading zero bucket is kept",
			tgt:  side{0, 5, []uint64{10}}, src: side{0, 2, []uint64{0, 1}},
			want: side{0, 2, []uint64{0, 1, 0, 10}},
		},
		{
			name: "empty source keeps target",
			tgt:  side{0, 5, []uint64{1}}, src: side{0, 0, []uint64{}},
			want: side{0, 5, []uint64{1}},
		},
		{
			name: "empty target takes source",
			tgt:  side{0, 0, []uint64{}}, src: side{0, 7, []uint64{1, 2}},
			want: side{0, 7, []uint64{1, 2}},
		},
		{
			name: "different scales merge at the lower one",
			tgt:  side{2, 4, []uint64{1, 1}}, src: side{1, 2, []uint64{3}},
			want: side{1, 2, []uint64{5}},
		},
		{
			name: "large gap downscales to the max size",
			tgt:  side{20, 0, []uint64{1}}, src: side{20, 1_000_000, []uint64{1}},
			want: side{7, 0, gap(123, 1, 1)},
		},
		{
			name: "offsets more than 2^31 apart",
			tgt:  side{20, 1 << 30, []uint64{1}}, src: side{20, -(1 << 30) - 1, []uint64{1}},
			want: side{-4, -65, gap(130, 1, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := func(s side) pmetric.ExponentialHistogramDataPoint {
				dp := pmetric.NewExponentialHistogramDataPoint()
				dp.SetScale(s.scale)
				dp.Positive().SetOffset(s.offset)
				dp.Positive().BucketCounts().FromRaw(s.counts)
				return dp
			}
			dp := point(tt.tgt)

			mergeExponentialHistogramBuckets(dp, point(tt.src))

			require.Equal(t, tt.want.scale, dp.Scale())
			require.Equal(t, tt.want.offset, dp.Positive().Offset())
			require.Equal(t, tt.want.counts, dp.Positive().BucketCounts().AsRaw())
			require.Zero(t, dp.Negative().BucketCounts().Len())
		})
	}
}

func Test_MergeDataPoints_ExpHistogramOrderIndependent(t *testing.T) {
	mapAttr := pcommon.NewMap()
	mapAttr.PutStr("attr1", "val1")
	hash := dataPointHashKey(mapAttr, pcommon.NewTimestampFromTime(time.Time{}), false, false, 0)

	points := func(offsets []int32, buckets [][]uint64) pmetric.ExponentialHistogramDataPointSlice {
		s := pmetric.NewExponentialHistogramDataPointSlice()
		for i := range offsets {
			d := s.AppendEmpty()
			d.Attributes().PutStr("attr1", "val1")
			var count uint64
			for _, c := range buckets[i] {
				count += c
			}
			d.SetCount(count)
			d.Positive().SetOffset(offsets[i])
			d.Positive().BucketCounts().FromRaw(buckets[i])
		}
		return s
	}

	merge := func(dps pmetric.ExponentialHistogramDataPointSlice) pmetric.ExponentialHistogramDataPoint {
		m := pmetric.NewMetric()
		m.SetEmptyExponentialHistogram()
		MergeDataPoints(m, Sum, AggGroups{expHistogram: map[string]pmetric.ExponentialHistogramDataPointSlice{hash: dps}})
		return m.ExponentialHistogram().DataPoints().At(0)
	}

	forward := merge(points([]int32{5, 3, 4}, [][]uint64{{10, 20}, {1, 0, 2}, {3}}))
	reverse := merge(points([]int32{4, 3, 5}, [][]uint64{{3}, {1, 0, 2}, {10, 20}}))

	for _, dp := range []pmetric.ExponentialHistogramDataPoint{forward, reverse} {
		require.Equal(t, int32(3), dp.Positive().Offset())
		require.Equal(t, []uint64{1, 3, 12, 20}, dp.Positive().BucketCounts().AsRaw())
		require.Equal(t, uint64(36), dp.Count())
	}
}

func Test_MergeDataPoints_ExpHistogramDifferentScales(t *testing.T) {
	in := pmetric.NewMetric()
	dps := in.SetEmptyExponentialHistogram().DataPoints()
	for _, p := range []struct {
		scale  int32
		offset int32
		counts []uint64
	}{
		{scale: 3, offset: 8, counts: []uint64{1, 1}},
		{scale: 2, offset: 4, counts: []uint64{3}},
	} {
		dp := dps.AppendEmpty()
		dp.Attributes().PutStr("attr1", "val1")
		dp.SetScale(p.scale)
		dp.SetCount(p.counts[0] + uint64(len(p.counts)-1))
		dp.Positive().SetOffset(p.offset)
		dp.Positive().BucketCounts().FromRaw(p.counts)
	}

	ag := AggGroups{}
	GroupDataPoints(in, &ag)
	out := pmetric.NewMetric()
	out.SetEmptyExponentialHistogram()
	MergeDataPoints(out, Sum, ag)

	// one point, so a backend keeping one point per series loses nothing
	require.Equal(t, 1, out.ExponentialHistogram().DataPoints().Len())
	dp := out.ExponentialHistogram().DataPoints().At(0)
	require.Equal(t, int32(2), dp.Scale())
	require.Equal(t, int32(4), dp.Positive().Offset())
	require.Equal(t, []uint64{5}, dp.Positive().BucketCounts().AsRaw())
	require.Equal(t, uint64(5), dp.Count())
}
