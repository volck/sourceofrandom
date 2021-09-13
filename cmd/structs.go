/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"github.com/drand/drand/client"
)

type drawInput struct {
	Min       int  `json:"Min"`
	Max       int  `json:"Max"`
	StartDraw int  `json:"StartDraw"`
	EndDraw   int  `json:"EndDraw"`
	PutBack   bool `json:"Putback"`
}

type newSeed struct {
	Drand         []byte `json:"Drand"`
	CalculateSeed int    `json:"CalculatedSeed"`
}

type Metadata struct {
	CreatedAt    int64         `json:"CreatedAt"`
	DrawNumber   int           `json: "DrawNumber"`
	DrandSeed    client.Result `json: "DrandSeed"`
	ComputedSeed int           `json: "computedSeed"`
}

type drawResults struct {
	Request  drawInput `json: "DrawInput"`
	Metadata Metadata  `json: "Metadata"`
	Numbers  []int     `json: "Numbers"`
}

var drandResult client.Result

var drawHistory []drawResults
