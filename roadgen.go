/*
   Copyright 2017 Hubert Jarosz

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"container/heap"
	"errors"
	"image"
	"math"
	"math/rand"
	"time"
)

func genRoads(mapRect *image.Rectangle, forbiddenAreas []image.Rectangle, startingRoad Road) []Road {

	rand.Seed(time.Now().UTC().UnixNano())

	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &item{value: roadWithMetadata{0, startingRoad}, priority: 0})

	segments := make([]Road, 0)

	for pq.Len() > 0 {

		road := heap.Pop(&pq).(*item).value

		roadSegment, err := localConstraints(segments, road.data, mapRect, forbiddenAreas)

		if err == nil {
			segments = append(segments, roadSegment)
			newRoads := genNewRoads(road, startingRoad.start)
			for _, r := range newRoads {
				heap.Push(&pq, &item{value: r, priority: r.branchDelay})
			}
		}

	}

	return segments
}

func roadStartOK(r Road, mapRect *image.Rectangle, forbiddenAreas []image.Rectangle) bool {

	if !r.start.In(*mapRect) {
		return false
	}

	for _, area := range forbiddenAreas {
		if r.start.In(area) {
			return false
		}
	}

	return true
}

func roadEndOK(r Road, mapRect *image.Rectangle, forbiddenAreas []image.Rectangle) bool {

	if !r.roadEnd().In(*mapRect) {
		return false
	}

	for _, area := range forbiddenAreas {
		if r.roadEnd().In(area) {
			return false
		}
	}

	return true
}

func localConstraints(segments []Road, r Road, mapRect *image.Rectangle, forbiddenAreas []image.Rectangle) (Road, error) {

	if !roadStartOK(r, mapRect, forbiddenAreas) {
		return r, errors.New("road not starting on map")
	}

	for !roadEndOK(r, mapRect, forbiddenAreas) {
		if r.length <= 10 {
			return r, errors.New("road not ending on map")
		}
		r.length -= 10
	}

	// check for overlaping with other roads

	for _, otherRoad := range segments {
		if r.isHorizontal() == otherRoad.isHorizontal() {
			if r.getRect().Overlaps(otherRoad.getRect()) {
				return r, errors.New("road overlaps with another road")
			}
		}
	}

	return r, nil
}

func genNewRoads(r roadWithMetadata, townCenter image.Point) []roadWithMetadata {
	roads := make([]roadWithMetadata, 0)

	branchDelay := r.branchDelay + 1

	densityVector := r.data.roadEnd().Sub(townCenter)
	density := int(math.Sqrt(float64(densityVector.X*densityVector.X + densityVector.Y*densityVector.Y)))

	if density < 50 {
		density = 50
	}

	if density > 500 {
		density = 500
	}

	for i := 0; i < 3; i++ {

		tmp := 0

		if density > 200 {
			tmp = rand.Intn(10) - 2
		}

		if tmp >= 0 {
			d := r.data

			d.start = r.data.roadEnd()

			d.angle += 90 * rand.Intn(4)
			if d.angle >= 360 {
				d.angle -= 360
			}

			//	if  < 200 {
			d.length = rand.Intn(density) - density/4
			// } else {
			// 	d.length = rand.Intn(100) + 50
			// }

			roads = append(roads, roadWithMetadata{branchDelay + tmp, d})
		}
	}

	return roads
}
