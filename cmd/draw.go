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
	"fmt"
	"math/rand"
	"time"

	"github.com/drand/drand/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// drawCmd represents the draw command
var diehardheader bool

var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Make sourceofrandom draw randomness",
	Long:  `Enters sourceofrandom into drawmode, allows for making files with uint32s for verifying randomness in e.g. dieharder`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("filename")
		maxint, err := cmd.Flags().GetInt("max")
		maxsize, err := cmd.Flags().GetString("maxsize")
		if err != nil {
			fmt.Println("error init flag", err)
		}
		client, err := newClient()
		if err != nil {
			fmt.Printf("failed to init client: %v\n", err)
		}

		go getWatcher(client)
		for {
			if drandResult != nil {
				fmt.Println("Going to draw random")
				makeASCIIinputforDieHarder(filename, maxint, maxsize)
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(drawCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	drawCmd.Flags().StringP("filename", "f", "diehard.output", "define filename which sourceofrandom will output to")
	drawCmd.Flags().Int("max", 0, "maximum number of draws sourceofrandom will do")
	drawCmd.Flags().String("maxsize", "", "maximum size of file generated")
	drawCmd.Flags().BoolVar(&diehardheader, "noheader", false, "inserts die hard header for 202 generator (file input)")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// drawCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func calculateSeed(dRandSeed client.Result, randomUUID uuid.UUID) (calculatedSeed int64) {
	newUuid := makeUuid()
	marshalBinary, err := newUuid.MarshalBinary()
	if err != nil {
		fmt.Println("marshalBinary failed")
	}
	nano := time.Now().UnixNano()
	calculatedSeed = BytesToInt64(dRandSeed.Randomness()) + BytesToInt64(marshalBinary) + nano
	return calculatedSeed
}

func makeUuid() uuid.UUID {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}
	return newUUID

}

func makeDrawRange(drawRequest drawInput) drawResults {
	var drawres drawResults

	drawres.Request = drawRequest
	drawres.Metadata.DrandSeed = drandResult
	for drawres.Request.StartDraw <= drawres.Request.EndDraw {
		calculatedSeed := calculateSeed(drawres.Metadata.DrandSeed, makeUuid())
		rand.Seed(calculatedSeed)
		newrand := rand.Intn(drawres.Request.Max-drawres.Request.Min) + drawRequest.Min
		if drawres.Request.PutBack {
			if !drawInResultsAlready(drawres.Numbers, newrand) {
				drawres.Numbers = append(drawres.Numbers, newrand)
			} else if drawInResultsAlready(drawres.Numbers, newrand) {
				drawres.Request.StartDraw = drawres.Request.StartDraw - 1

			}
		} else {
			fmt.Printf("draw(%d): putback not defined. we put every result in here \n", drawRequest.StartDraw)
			drawres.Numbers = append(drawres.Numbers, newrand)
		}
		drawres.Request.StartDraw = drawres.Request.StartDraw + 1

	}
	drawres.Metadata.CreatedAt = time.Now().Unix()
	return drawres
}

func makeASCIIinputforDieHarder(filename string, maxint int, maxsize string) {
	fmt.Println(filename, maxint, maxsize)
	dieharderFile := createandDefer(filename)
	fmt.Printf("max size is %s\n", maxsize)
	if maxint != 0 {
		dieharderFile.WriteString(insertDiehardHeader(maxint))
	}
	i := 0
	for {
		checkFileSizeAndFileLength(maxsize, filename, int64(maxint), i)
		draw := draw()
		dieharderFile.WriteString(draw)
		printStats(i, draw, filename)
		i = i + 1
	}
}

func draw() string {
	calculatedSeed := calculateSeed(drandResult, makeUuid())
	rand.Seed(calculatedSeed)
	mynewrandInt := rand.Uint32()
	return fmt.Sprintf("%d\n", mynewrandInt)
}

func drawInResultsAlready(results []int, randno int) bool {
	inList := false

	for i := range results {
		if i == randno {
			inList = true
		}
	}
	return inList

}
