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

func genRoads(mapRect *image.Rectangle) []Road {

	mapCenter := image.Pt(mapRect.Min.X+mapRect.Size().X/2, mapRect.Min.Y+mapRect.Size().Y/2)

	rand.Seed(time.Now().UTC().UnixNano())

	road := roadWithMetadata{0, Road{mapCenter, 10, 0, 3}}

	pq := make(priorityQueue, 1)
	pq[0] = &item{value: road, priority: 0, index: 0}
	heap.Init(&pq)

	segments := make([]Road, 0)

	for pq.Len() > 0 {

		road := heap.Pop(&pq).(*item).value

		roadSegment, err := localConstraints(segments, road.data, mapRect)

		if err == nil {
			segments = append(segments, roadSegment)
			newRoads := genNewRoads(road, mapCenter)
			for _, r := range newRoads {
				heap.Push(&pq, &item{value: r, priority: r.branchDelay})
			}
		}

	}

	return segments
}

func localConstraints(segments []Road, d Road, mapRect *image.Rectangle) (Road, error) {

	if !d.start.In(*mapRect) {
		return d, errors.New("road not starting on map")
	}

	for !d.roadEnd().In(*mapRect) {
		if d.length <= 10 {
			return d, errors.New("road not ending on map")
		}
		d.length -= 10
	}

	// check for overlaping with other roads

	for _, otherRoad := range segments {
		if d.isHorizontal() == otherRoad.isHorizontal() {
			if d.getRect().Overlaps(otherRoad.getRect()) {
				return d, errors.New("road overlaps with another road")
			}
		}
	}

	return d, nil
}

func genNewRoads(r roadWithMetadata, mapCenter image.Point) []roadWithMetadata {
	roads := make([]roadWithMetadata, 0)

	branchDelay := r.branchDelay + 1

	densityVector := r.data.roadEnd().Sub(mapCenter)
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
