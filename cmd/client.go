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
	"context"
	"encoding/hex"
	"fmt"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
)

var chainHash, _ = hex.DecodeString("8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce")

var urls = []string{
	"https://api.drand.sh",
	"https://drand.cloudflare.com",
	"https://api.drand.sh",
	"https://api2.drand.sh",
	"https://api3.drand.sh",
	"https://drand.cloudflare.com",
}

func newClient() (client.Client, error) {
	c, err := client.New(
		client.From(http.ForURLs(urls, chainHash)...),
		client.WithChainHash(chainHash),
		client.WithAutoWatch(),
	)
	return c, err
}

func getRandomFromDrand(client client.Client, round int) client.Result {
	r, err := client.Get(context.Background(), 0)
	if err != nil {
		fmt.Println("err", err)
	}
	return r
}

func getWatcher(c client.Client) {
	fmt.Println("initializing watcher...")
	for {
		drandResult = <-c.Watch(context.Background())
	}
}
