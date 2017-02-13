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

/*
TODO building types:
* military
* religious
* rulers
* buisness
*/

package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
)

func drawRoad(gc *draw2dimg.GraphicContext, road Road) {

	roadEnd := road.roadEnd()

	gc.SetLineWidth(float64(road.size))
	gc.SetStrokeColor(color.Black)

	// colorful debug
	// gc.SetStrokeColor(color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255})

	gc.MoveTo(float64(road.start.X), float64(road.start.Y))
	gc.LineTo(float64(roadEnd.X), float64(roadEnd.Y))
	gc.Close()
	gc.FillStroke()
}

func drawRect(gc *draw2dimg.GraphicContext, rect image.Rectangle) {

	gc.SetStrokeColor(color.RGBA{0, 0, 255, 255})
	gc.SetLineWidth(float64(rect.Size().Y))

	gc.MoveTo(float64(rect.Min.X), float64(rect.Min.Y+rect.Size().Y/2))
	gc.LineTo(float64(rect.Max.X), float64(rect.Min.Y+rect.Size().Y/2))
	gc.Close()
	gc.FillStroke()
}
