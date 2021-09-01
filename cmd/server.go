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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "enables web server API",
	Long:  `starts web server`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Println(err)
		}
		client, err := newClient()
		if err != nil {
			fmt.Printf("failed to init client: %v\n", err)
		}
		go getWatcher(client)
		for {
			if drandResult != nil {
				serve(port)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().Int("port", 1337, "Port to serve web server")

}

func serve(port int) {
	fmt.Printf("[*] Going to serve at %d[*]\n", port)
	http.HandleFunc("/seed", seed)
	http.HandleFunc("/info", info)
	http.HandleFunc("/getRound", getRound)
	http.HandleFunc("/doDraw", doDraw)
	serveString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(serveString, nil)
}

func getRound(w http.ResponseWriter, req *http.Request) {

	msg := fmt.Sprintf("latest getRound is %v", drandResult.Round())

	fmt.Fprintf(w, msg)
}

func info(w http.ResponseWriter, req *http.Request) {

	msg := fmt.Sprintf("latest random is %v", drandResult.Randomness())

	fmt.Fprintf(w, msg)
}

func doDraw(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		reqBody, _ := ioutil.ReadAll(req.Body)
		var draw drawInput
		err := json.Unmarshal(reqBody, &draw)

		if err != nil {
			fmt.Println(err)
		}
		theDraw := makeDrawRange(draw)

		msg := fmt.Sprintf("%v\n", theDraw)
		fmt.Fprintf(w, msg)

	}

}

func seed(w http.ResponseWriter, r *http.Request) {

	newUuid := makeUuid()
	calculatedSeed := calculateSeed(drandResult, newUuid)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	newSeed := make(map[string]interface{})
	newSeed["drandPublic"] = drandResult.Randomness()
	newSeed["calculatedSeed"] = calculatedSeed
	jsonResp, err := json.Marshal(newSeed)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
