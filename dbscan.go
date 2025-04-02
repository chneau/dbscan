package dbscan

import (
	"math"
	"math/rand"
)

// Observation is an interface that represents a point in the dataset.
type Observation interface {
	Coordinates() []float64
}

// distance computes the Euclidean distance between two observations.
func distance(a, b Observation) float64 {
	coordsA := a.Coordinates()
	coordsB := b.Coordinates()
	if len(coordsA) != len(coordsB) {
		panic("Observations have different dimensions")
	}
	sum := 0.0
	for i := range coordsA {
		diff := coordsA[i] - coordsB[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// getNeighborIndices returns the indices of points within epsilon distance from dataset[index].
func getNeighborIndices[T Observation](dataset []T, index int, epsilon float64) []int {
	var neighbors []int
	point := dataset[index]
	for i, other := range dataset {
		if distance(point, other) <= epsilon {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

// Cluster performs DBSCAN clustering on the given dataset.
func Cluster[T Observation](dataset []T, minDensity int, epsilon float64, rng *rand.Rand) ([][]T, error) {
	// Track visited points
	visited := make([]bool, len(dataset))
	// Store resulting clusters
	clusters := [][]T{}

	// Process all points
	for {
		// Collect indices of unvisited points
		var unvisited []int
		for i, v := range visited {
			if !v {
				unvisited = append(unvisited, i)
			}
		}
		// Exit if all points are visited
		if len(unvisited) == 0 {
			break
		}

		// Select a random unvisited point
		randIndex := rng.Intn(len(unvisited))
		i := unvisited[randIndex]
		visited[i] = true

		// Get neighbors (including the point itself)
		neighborIndices := getNeighborIndices(dataset, i, epsilon)
		if len(neighborIndices) < minDensity {
			// Point is noise; skip it (noise not returned)
			continue
		}

		// Start a new cluster
		var cluster []T
		cluster = append(cluster, dataset[i])
		// Use a queue to expand the cluster
		queue := make([]int, 0, len(neighborIndices))
		queue = append(queue, neighborIndices...)

		// Expand the cluster
		for len(queue) > 0 {
			// Dequeue the next point
			j := queue[0]
			queue = queue[1:]

			// Process only unvisited points
			if !visited[j] {
				visited[j] = true
				cluster = append(cluster, dataset[j])
				// Check if this point is a core point
				neighborIndicesJ := getNeighborIndices(dataset, j, epsilon)
				if len(neighborIndicesJ) >= minDensity {
					// Add all neighbors to queue (may include visited points, but skipped above)
					queue = append(queue, neighborIndicesJ...)
				}
			}
		}
		// Add completed cluster to clusters
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}
