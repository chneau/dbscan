package dbscan

import (
	"math/rand"
	"slices"
	"testing"
)

type Numbers int64

func (e Numbers) Coordinates() []float64 {
	return []float64{float64(e)}
}

type Coordinates [2]int

func (c Coordinates) Coordinates() []float64 {
	return []float64{float64(c[0]), float64(c[1])}
}

func TestClusterNumbers(t *testing.T) {
	dataset := []Numbers{
		1, 2, 3,
		11, 12, 13,
		21, 22, 23,
		100,
	}
	minDensity := 1
	epsilon := 2.0
	rng := rand.New(rand.NewSource(0))

	clusters, err := Cluster(dataset, minDensity, epsilon, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedClusters := [][]Numbers{
		{1, 2, 3},
		{11, 12, 13},
		{21, 22, 23},
		{100},
	}

	// Check if the clusters match the expected clusters (unordered)
	if len(clusters) != len(expectedClusters) {
		t.Fatalf("expected %d clusters, got %d", len(expectedClusters), len(clusters))
	}

	matched := make([]bool, len(expectedClusters))
	for _, cluster := range clusters {
		slices.Sort(cluster)
		found := false
		for i, expected := range expectedClusters {
			if !matched[i] && slices.Equal(cluster, expected) {
				matched[i] = true
				found = true
				break
			}
		}
		if !found {
			t.Errorf("unexpected cluster: %v", cluster)
		}
	}

	for i, matched := range matched {
		if !matched {
			t.Errorf("expected cluster %v not found", expectedClusters[i])
		}
	}
}

func TestClusterCoordinates(t *testing.T) {
	dataset := []Coordinates{
		{1, 2}, {2, 3}, {3, 4},
		{11, 12}, {12, 13}, {13, 14},
		{21, 22}, {22, 23}, {23, 24},
		{100, 200},
	}
	minDensity := 1
	epsilon := 4.0
	rng := rand.New(rand.NewSource(0))

	clusters, err := Cluster(dataset, minDensity, epsilon, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedClusters := [][]Coordinates{
		{{1, 2}, {2, 3}, {3, 4}},
		{{11, 12}, {12, 13}, {13, 14}},
		{{21, 22}, {22, 23}, {23, 24}},
		{{100, 200}},
	}

	// Check if the clusters match the expected clusters (unordered)
	if len(clusters) != len(expectedClusters) {
		t.Fatalf("expected %d clusters, got %d", len(expectedClusters), len(clusters))
	}

	matched := make([]bool, len(expectedClusters))
	for _, cluster := range clusters {
		slices.SortFunc(cluster, func(a, b Coordinates) int {
			return a[0] - b[0]
		})
		found := false
		for i, expected := range expectedClusters {
			if !matched[i] && slices.Equal(cluster, expected) {
				matched[i] = true
				found = true
				break
			}
		}
		if !found {
			t.Errorf("unexpected cluster: %v", cluster)
		}
	}

	for i, matched := range matched {
		if !matched {
			t.Errorf("expected cluster %v not found", expectedClusters[i])
		}
	}
}
