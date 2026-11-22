package analyzer

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	MinCount      = 5
	MaxCount      = 100
	MaxCityLength = 100

	MinWindSpeed = 0.0
	MaxWindSpeed = 150.0

	CategoryCalm     = "Штиль"
	CategoryLight    = "Слабкий"
	CategoryModerate = "Помірний"
	CategoryStrong   = "Сильний"
)

type CategoryStat struct {
	Name       string
	RangeLabel string
	Count      int
}
type RepeatStat struct {
	Value float64
	Count int
}
type WindSeries struct {
	city   string
	speeds []float64
}
type WindReport struct {
	City       string
	Count      int
	Min        float64
	Max        float64
	Sum        float64
	Average    float64
	Categories []CategoryStat
	Repeats    []RepeatStat
}

func NewWindSeries(city string, speeds []float64) (*WindSeries, error) {
	trimmedCity := strings.TrimSpace(city)
	if trimmedCity == "" {
		return nil, errors.New("назва міста не може бути порожньою")
	}
	if len([]rune(trimmedCity)) > MaxCityLength {
		return nil, fmt.Errorf("назва міста перевищує допустиму довжину (%d символів)", MaxCityLength)
	}
	if len(speeds) == 0 {
		return nil, errors.New("набір швидкостей вітру не може бути порожнім")
	}
	if len(speeds) < MinCount || len(speeds) > MaxCount {
		return nil, fmt.Errorf("кількість вимірювань повинна бути від %d до %d (отримано %d)", MinCount, MaxCount, len(speeds))
	}

	copied := make([]float64, len(speeds))
	for i, v := range speeds {
		if math.IsNaN(v) || math.IsInf(v, 0) || v < MinWindSpeed || v > MaxWindSpeed {
			return nil, fmt.Errorf("значення швидкості [%d] = %.2f виходить за межі від %.0f до %.0f км/год", i, v, MinWindSpeed, MaxWindSpeed)
		}
		copied[i] = v
	}

	return &WindSeries{
		city:   trimmedCity,
		speeds: copied,
	}, nil
}
func (s *WindSeries) City() string {
	return s.city
}
func (s *WindSeries) Speeds() []float64 {
	out := make([]float64, len(s.speeds))
	copy(out, s.speeds)
	return out
}
func (s *WindSeries) Analyze() WindReport {
	count := len(s.speeds)
	minVal := s.speeds[0]
	maxVal := s.speeds[0]
	sum := 0.0

	var calm, light, moderate, strong int
	freqMap := make(map[float64]int, count)
	order := make([]float64, 0, count)

	for _, v := range s.speeds {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += v

		switch {
		case v <= 5.0:
			calm++
		case v <= 30.0:
			light++
		case v <= 60.0:
			moderate++
		default:
			strong++
		}

		if freqMap[v] == 0 {
			order = append(order, v)
		}
		freqMap[v]++
	}

	average := sum / float64(count)

	categories := []CategoryStat{
		{Name: CategoryCalm, RangeLabel: "<= 5", Count: calm},
		{Name: CategoryLight, RangeLabel: "5.1 - 30", Count: light},
		{Name: CategoryModerate, RangeLabel: "30.1 - 60", Count: moderate},
		{Name: CategoryStrong, RangeLabel: "> 60", Count: strong},
	}

	var repeats []RepeatStat
	for _, v := range order {
		if freqMap[v] > 1 {
			repeats = append(repeats, RepeatStat{
				Value: v,
				Count: freqMap[v],
			})
		}
	}

	return WindReport{
		City:       s.city,
		Count:      count,
		Min:        minVal,
		Max:        maxVal,
		Sum:        sum,
		Average:    average,
		Categories: categories,
		Repeats:    repeats,
	}
}
