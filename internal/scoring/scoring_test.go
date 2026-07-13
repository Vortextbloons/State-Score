package scoring

import (
	"math"
	"testing"
)

func ptr(v float64) *float64 { return &v }

func TestNormalizeMinMaxDirectionAndClamping(t *testing.T) {
	observations := []Observation{{1, ptr(10)}, {2, ptr(20)}, {3, ptr(30)}, {4, nil}}
	got, err := Normalize(observations, MinMax, true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got[1], 0)
	assertClose(t, got[2], 50)
	assertClose(t, got[3], 100)
	if _, ok := got[4]; ok {
		t.Fatal("missing value must not be assigned a score")
	}
	got, err = Normalize(observations, MinMax, false, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got[1], 100)
	assertClose(t, got[3], 0)
}

func TestNormalizePercentileUsesAverageRankForTies(t *testing.T) {
	got, err := Normalize([]Observation{{1, ptr(10)}, {2, ptr(20)}, {3, ptr(20)}, {4, ptr(40)}}, Percentile, true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got[1], 0)
	assertClose(t, got[2], 50)
	assertClose(t, got[3], 50)
	assertClose(t, got[4], 100)
}

func TestNormalizeZScoreAndFixedThreshold(t *testing.T) {
	z, err := Normalize([]Observation{{1, ptr(0)}, {2, ptr(10)}, {3, ptr(20)}}, ZScore, true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, z[2], 50)
	if !(z[1] < z[2] && z[2] < z[3]) {
		t.Fatalf("z-score ordering wrong: %#v", z)
	}
	fixed, err := Normalize([]Observation{{1, ptr(-5)}, {2, ptr(50)}, {3, ptr(120)}}, FixedThreshold, true, ptr(0), ptr(100))
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, fixed[1], 0)
	assertClose(t, fixed[2], 50)
	assertClose(t, fixed[3], 100)
}

func TestWeightedAverageRedistributesMissingWeight(t *testing.T) {
	got, err := WeightedAverage(map[int64]*float64{1: ptr(80), 2: nil}, map[int64]float64{1: 1, 2: 3})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got.Score, 80)
	assertClose(t, got.Completeness, .25)
	if !got.Incomplete {
		t.Fatal("missing weighted metric must mark score incomplete")
	}
}

func TestCalculateCategoryOverallAndCompleteness(t *testing.T) {
	got, err := Calculate([]int64{10}, []CategoryInput{{1, .6}, {2, .4}}, []MetricInput{
		{ID: 1, CategoryID: 1, Weight: 1, Scores: map[int64]float64{10: 100}},
		{ID: 2, CategoryID: 1, Weight: 1, Scores: map[int64]float64{}},
		{ID: 3, CategoryID: 2, Weight: 1, Scores: map[int64]float64{10: 50}},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, got[0].Categories[1].Score, 100)
	assertClose(t, got[0].Overall.Score, 80)
	assertClose(t, got[0].Overall.Completeness, .7) // .6*50% + .4*100%
}

func TestValidationAndDegenerateSeries(t *testing.T) {
	constant, err := Normalize([]Observation{{1, ptr(7)}, {2, ptr(7)}}, MinMax, true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, constant[1], 50)
	assertClose(t, constant[2], 50)
	if _, err = Normalize([]Observation{{1, ptr(1)}}, "mystery", true, nil, nil); err == nil {
		t.Fatal("expected unsupported method error")
	}
	if _, err = WeightedAverage(map[int64]*float64{}, map[int64]float64{1: -1}); err == nil {
		t.Fatal("expected invalid weight error")
	}
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("got %v, want %v", got, want)
	}
}
