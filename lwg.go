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
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	margin = 20
	width  = 800
	height = 600
)

type lwgTown struct {
	Name         string `json:"name"`
	HasRiver     string `json:"isThereRiver"`
	WaterOnNorth bool   `json:"waterOnTop"`
	WaterOnEast  bool   `json:"waterOnRight"`
	WaterOnSouth bool   `json:"waterOnBottom"`
	WaterOnWest  bool   `json:"waterOnLeft"`
}

func loadLWGTowns(path string) []lwgTown {
	var towns []lwgTown

	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(f, &towns)

	return towns
}

func (t *lwgTown) saneName() string {
	sanitize := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		}
		return '_'
	}
	return strings.Map(sanitize, t.Name)
}

func drawLWGTown(dir string, town lwgTown) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	gc := draw2dimg.NewGraphicContext(img)

	forbiddenAreas := make([]image.Rectangle, 0)

	// keep some distance from map sides
	forbiddenAreas = append(forbiddenAreas, image.Rect(0, 0, margin, height))            // left
	forbiddenAreas = append(forbiddenAreas, image.Rect(0, 0, width, margin))             // top
	forbiddenAreas = append(forbiddenAreas, image.Rect(0, height-margin, width, height)) // bottom
	forbiddenAreas = append(forbiddenAreas, image.Rect(width-margin, 0, width, height))  // right

	// coast
	if town.WaterOnNorth {
		waterSize := rand.Intn(20) + margin
		water := image.Rect(0, 0, width, waterSize)
		forbiddenAreas = append(forbiddenAreas, water)
		drawRect(gc, water)
	}
	if town.WaterOnEast {
		waterSize := rand.Intn(20) + margin
		water := image.Rect(width-waterSize, 0, width, height)
		forbiddenAreas = append(forbiddenAreas, water)
		drawRect(gc, water)
	}
	if town.WaterOnSouth {
		waterSize := rand.Intn(20) + margin
		water := image.Rect(0, height-waterSize, width, height)
		forbiddenAreas = append(forbiddenAreas, water)
		drawRect(gc, water)
	}
	if town.WaterOnWest {
		waterSize := rand.Intn(20) + margin
		water := image.Rect(0, 0, waterSize, height)
		forbiddenAreas = append(forbiddenAreas, water)
		drawRect(gc, water)
	}

	startingRoad := Road{
		start:  image.Pt(width/2, height/2),
		length: 200,
		angle:  0,
		size:   3,
	}

	roads := genRoads(&img.Rect, forbiddenAreas, startingRoad)

	//sometimes we got empty because of bad random seed
	for len(roads) <= 50 {
		roads = genRoads(&img.Rect, forbiddenAreas, startingRoad)
	}

	for _, road := range roads {
		drawRoad(gc, road)
	}

	f, err := os.Create(filepath.Join(dir, town.saneName()+".png"))
	if err != nil {
		return err
	}
	defer f.Close()
	png.Encode(f, img)

	return nil
}

func drawLWGtowns(path string) error {
	towns := loadLWGTowns(path)

	dir := filepath.Join(filepath.Dir(path), "towns")

	err := os.MkdirAll(dir, os.ModePerm) // will have 777 perm.
	if err != nil {
		return err
	}

	for _, t := range towns {
		err := drawLWGTown(dir, t)
		if err != nil {
			return err
		}
	}

	return nil
}
