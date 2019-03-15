package main

import (
	"log"
	"time"
)

var (
	maxDeviation    float64 = 10000
	minDeviation    float64 = 0
	maxConnectivity float64 = 10000
	minConnectivity float64 = 0
)

func setObjectivesMaxMinValues(img *Image) {

	width := len(*img)
	height := len((*img)[0])
	// Make image where all groups are in different segments

	connectedGroupsInDifferent := make([]map[uint64]bool, width*height)

	for x := range *img {
		for y := range (*img)[x] {
			idx := Expand(width, x, y)
			connectedGroupsInDifferent[idx] = map[uint64]bool{uint64(idx): true}
		}
	}
	maxConnectivity = connectivity(img, connectedGroupsInDifferent)

	connectedGroupsInSame := make([]map[uint64]bool, 1)
	connectedGroupsInSame[0] = make(map[uint64]bool)

	for x := range *img {
		for y := range (*img)[x] {
			idx := Expand(width, x, y)
			connectedGroupsInSame[0][uint64(idx)] = true
		}
	}

	maxDeviation = deviation(img, connectedGroupsInSame)

	// max deviation -> all pixels in one segment

	// min deviation -> all pixels in their own segment = 0

	// max connectivity -> all pixels in their own segment

	// min connectivity -> all pixels in one segment
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}


func deviation(img *Image, connectedGroups []map[uint64]bool) float64 {
	defer timeTrack(time.Now(), "deviation")
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

func connectivity(img *Image, connectedGroups []map[uint64]bool) float64 {
	defer timeTrack(time.Now(), "connectivity")

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