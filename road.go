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
	"image"
	"math"
)

const (
	density = 10
)

//
type roadWithMetadata struct {
	branchDelay int
	data        Road
}

// Road type contains roads segment data.
type Road struct {
	start  image.Point
	length int
	angle  int
	size   int
}

func (d *Road) roadEnd() image.Point {
	fi := float64(d.angle) * math.Pi / 180
	xdiff := int(float64(d.length) * math.Cos(fi))
	ydiff := int(float64(d.length) * math.Sin(fi))

	return d.start.Add(image.Pt(xdiff, ydiff))
}

func (d *Road) isHorizontal() bool {
	if d.start.X == d.roadEnd().X {
		return false
	}
	return true
}

func (d *Road) getRect() image.Rectangle {
	x1, y1 := d.start.X, d.start.Y
	roadEnd := d.roadEnd()
	x2, y2 := roadEnd.X, roadEnd.Y

	if d.isHorizontal() {
		return image.Rect(x1, y1-density, x2, y2+density)
	}
	return image.Rect(x1-density, y1, x2+density, y2)

}
