package main

var (
	maxDeviation    float64 = 100
	minDeviation    float64 = 0
	maxConnectivity float64 = 100
	minConnectivity float64 = 0
)

func deviation(img *Image, connectedGroups []map[uint64]bool) float64 {

	var dist float64
	width := len(*img)
	for _, group := range connectedGroups {
		centroid := Centroid(img, group)

		for k := range group {
			x, y := Flatten(width, int(k))
			dist += ColorDist(&(*img)[x][y], centroid)
		}
	}
	return dist
}

func connectiviy(img *Image, connectedGroups []map[uint64]bool) float64 {

	var dist float64

	for _, group := range connectedGroups {
		for k := range group {
			intK := int(k)
			for j, neighbour := range GetTargets(img, intK) {

				if _, ok := group[uint64(neighbour)]; ok { // To nothing
				} else {
					dist += 1.0 / (float64(j) + 1.0)
				}
			}
		}
	}
	return dist
}
