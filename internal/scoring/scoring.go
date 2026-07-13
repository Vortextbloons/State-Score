// Package scoring implements StateScore's deterministic scoring formulas.
package scoring

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

const CalculationVersion = "1"

type NormalizationMethod string

const (
	MinMax         NormalizationMethod = "minmax"
	Percentile     NormalizationMethod = "percentile"
	ZScore         NormalizationMethod = "zscore"
	FixedThreshold NormalizationMethod = "fixed"
)

// Observation is a value belonging to one state. A nil Value is missing data.
type Observation struct {
	StateID int64
	Value   *float64
}

// Normalize converts present observations to scores in [0, 100]. For min-max
// and fixed normalization, minimum and maximum are the observed and policy
// bounds respectively. Missing observations are omitted from the result.
func Normalize(observations []Observation, method NormalizationMethod, higherIsBetter bool, minimum, maximum *float64) (map[int64]float64, error) {
	values := make([]Observation, 0, len(observations))
	seen := make(map[int64]struct{}, len(observations))
	for _, o := range observations {
		if _, ok := seen[o.StateID]; ok {
			return nil, fmt.Errorf("duplicate observation for state %d", o.StateID)
		}
		seen[o.StateID] = struct{}{}
		if o.Value == nil {
			continue
		}
		if math.IsNaN(*o.Value) || math.IsInf(*o.Value, 0) {
			return nil, fmt.Errorf("state %d has a non-finite value", o.StateID)
		}
		values = append(values, o)
	}
	result := make(map[int64]float64, len(values))
	if len(values) == 0 {
		return result, nil
	}

	method = NormalizationMethod(strings.ToLower(string(method)))
	switch method {
	case MinMax, FixedThreshold:
		lo, hi := observedBounds(values)
		if method == FixedThreshold {
			if minimum == nil || maximum == nil {
				return nil, errors.New("fixed normalization requires minimum and maximum thresholds")
			}
			lo, hi = *minimum, *maximum
		} else if minimum != nil && maximum != nil {
			lo, hi = *minimum, *maximum
		}
		if hi <= lo {
			if hi == lo {
				for _, o := range values {
					result[o.StateID] = 50
				}
				return result, nil
			}
			return nil, errors.New("normalization maximum must be greater than minimum")
		}
		for _, o := range values {
			score := (*o.Value - lo) / (hi - lo) * 100
			if !higherIsBetter {
				score = 100 - score
			}
			result[o.StateID] = clamp(score)
		}
	case Percentile:
		sorted := append([]Observation(nil), values...)
		sort.Slice(sorted, func(i, j int) bool { return *sorted[i].Value < *sorted[j].Value })
		if len(sorted) == 1 {
			result[sorted[0].StateID] = 50
			return result, nil
		}
		// Ties receive their average rank, preventing input order from affecting scores.
		for start := 0; start < len(sorted); {
			end := start + 1
			for end < len(sorted) && *sorted[end].Value == *sorted[start].Value {
				end++
			}
			rank := float64(start+end-1) / 2
			score := rank / float64(len(sorted)-1) * 100
			if !higherIsBetter {
				score = 100 - score
			}
			for i := start; i < end; i++ {
				result[sorted[i].StateID] = score
			}
			start = end
		}
	case ZScore:
		var mean float64
		for _, o := range values {
			mean += *o.Value
		}
		mean /= float64(len(values))
		var variance float64
		for _, o := range values {
			d := *o.Value - mean
			variance += d * d
		}
		variance /= float64(len(values))
		if variance == 0 {
			for _, o := range values {
				result[o.StateID] = 50
			}
			return result, nil
		}
		stddev := math.Sqrt(variance)
		for _, o := range values {
			// The standard normal CDF provides a stable, bounded 0-100 conversion.
			score := 50 * (1 + math.Erf(((*o.Value-mean)/stddev)/math.Sqrt2))
			if !higherIsBetter {
				score = 100 - score
			}
			result[o.StateID] = clamp(score)
		}
	default:
		return nil, fmt.Errorf("unsupported normalization method %q", method)
	}
	return result, nil
}

type WeightedScore struct {
	Score, Completeness float64
	Incomplete          bool
}

// WeightedAverage excludes missing scores and redistributes their weights.
// Completeness is the included positive weight divided by all positive weight.
func WeightedAverage(scores map[int64]*float64, weights map[int64]float64) (WeightedScore, error) {
	var total, included, sum float64
	for id, weight := range weights {
		if math.IsNaN(weight) || math.IsInf(weight, 0) || weight < 0 {
			return WeightedScore{}, fmt.Errorf("weight %d must be finite and non-negative", id)
		}
		if weight == 0 {
			continue
		}
		total += weight
		if score := scores[id]; score != nil {
			if math.IsNaN(*score) || math.IsInf(*score, 0) || *score < 0 || *score > 100 {
				return WeightedScore{}, fmt.Errorf("score %d must be between 0 and 100", id)
			}
			included += weight
			sum += *score * weight
		}
	}
	if total == 0 {
		return WeightedScore{}, errors.New("at least one positive weight is required")
	}
	if included == 0 {
		return WeightedScore{Completeness: 0, Incomplete: true}, nil
	}
	completeness := included / total
	return WeightedScore{Score: sum / included, Completeness: completeness, Incomplete: completeness < 1}, nil
}

type MetricInput struct {
	ID, CategoryID int64
	Weight         float64
	Scores         map[int64]float64
}
type CategoryInput struct {
	ID     int64
	Weight float64
}
type StateScore struct {
	StateID    int64
	Overall    WeightedScore
	Categories map[int64]WeightedScore
}

// Calculate composes normalized metric scores into category and overall scores.
func Calculate(stateIDs []int64, categories []CategoryInput, metrics []MetricInput) ([]StateScore, error) {
	categoryWeights := make(map[int64]float64, len(categories))
	for _, c := range categories {
		categoryWeights[c.ID] = c.Weight
	}
	results := make([]StateScore, 0, len(stateIDs))
	for _, stateID := range stateIDs {
		categoryScores := make(map[int64]WeightedScore, len(categories))
		overallValues := make(map[int64]*float64, len(categories))
		for _, category := range categories {
			values, weights := map[int64]*float64{}, map[int64]float64{}
			for _, metric := range metrics {
				if metric.CategoryID != category.ID {
					continue
				}
				weights[metric.ID] = metric.Weight
				if score, ok := metric.Scores[stateID]; ok {
					s := score
					values[metric.ID] = &s
				}
			}
			categoryScore, err := WeightedAverage(values, weights)
			if err != nil {
				return nil, fmt.Errorf("category %d: %w", category.ID, err)
			}
			categoryScores[category.ID] = categoryScore
			if categoryScore.Completeness > 0 {
				s := categoryScore.Score
				overallValues[category.ID] = &s
			}
		}
		overall, err := WeightedAverage(overallValues, categoryWeights)
		if err != nil {
			return nil, fmt.Errorf("overall score for state %d: %w", stateID, err)
		}
		// Overall completeness accounts for partially complete categories too.
		var expected, present float64
		for _, c := range categories {
			if c.Weight > 0 {
				expected += c.Weight
				present += c.Weight * categoryScores[c.ID].Completeness
			}
		}
		if expected > 0 {
			overall.Completeness = present / expected
			overall.Incomplete = overall.Completeness < 1
		}
		results = append(results, StateScore{StateID: stateID, Overall: overall, Categories: categoryScores})
	}
	return results, nil
}

func observedBounds(values []Observation) (float64, float64) {
	lo, hi := *values[0].Value, *values[0].Value
	for _, o := range values[1:] {
		if *o.Value < lo {
			lo = *o.Value
		}
		if *o.Value > hi {
			hi = *o.Value
		}
	}
	return lo, hi
}
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
