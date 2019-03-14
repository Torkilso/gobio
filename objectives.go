package main

var (
	maxDeviation    float64 = 10000
	minDeviation    float64 = 0
	maxConnectivity float64 = 10000
	minConnectivity float64 = 0
)

func setObjectivesMaxMinValues(img *Image) {
	// max deviation -> all pixels in one segment

	// min deviation -> all pixels in their own segment = 0

	// max connectivity -> all pixels in their own segment

	// min connectivity -> all pixels in one segment
}

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
