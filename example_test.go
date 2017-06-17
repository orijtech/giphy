// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package giphy_test

import (
	"fmt"
	"log"

	"github.com/orijtech/giphy/v1"
)

func Example_client_Trending() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Trending(&giphy.Request{
		MaxPageNumber: 4,
	})
	if err != nil {
		log.Fatal(err)
	}

	for page := range res.Pages {
		if page.Err != nil {
			log.Printf("#%d: err: %v\n", page.PageNumber, page.Err)
			continue
		}

		log.Printf("PageNumber: %#d\n", page.PageNumber)
		for i, giph := range page.Giphs {
			log.Printf("\t%d: %#v\tSizes: %#v\n", i, giph, giph.Sizes)
		}
	}
}

func Example_client_Search() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Search(&giphy.Request{
		MaxPageNumber: 4,

		Query: "Milly Rock",
	})
	if err != nil {
		log.Fatal(err)
	}

	for page := range res.Pages {
		if page.Err != nil {
			fmt.Printf("#%d: err: %v\n", page.PageNumber, page.Err)
			continue
		}

		fmt.Printf("PageNumber: %#d\n", page.PageNumber)
		for _, giph := range page.Giphs {
			fmt.Printf("\tDescription: %v\n", giph)
			for size, detail := range giph.Sizes {
				fmt.Printf("\t\t%s: %#v\n", size, detail)
			}
		}
		fmt.Printf("\n\n")
	}
}

func Example_client_RandomGIF() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	giph, err := client.RandomGIF(&giphy.Request{
		Tag:    "netflix",
		Rating: giphy.RatingPG,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Your giph: %#v\n", giph)
}

func Example_client_GIFByID() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	giph, err := client.GIFByID("3ohze2UfcItWPUFqbm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The retrieved giph: %#v\n", giph)
}

func Example_client_TrendingStickers() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.TrendingStickers(&giphy.Request{
		MaxPageNumber: 4,
	})
	if err != nil {
		log.Fatal(err)
	}

	for page := range res.Pages {
		if page.Err != nil {
			log.Printf("#%d: err: %v\n", page.PageNumber, page.Err)
			continue
		}

		log.Printf("PageNumber: %#d\n", page.PageNumber)
		for i, giph := range page.Giphs {
			log.Printf("\t%d: %#v\tSizes: %#v\n", i, giph, giph.Sizes)
		}
	}
}

func Example_client_SearchStickers() {
	client, err := giphy.NewClientFromEnvOrDefault()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.SearchStickers(&giphy.Request{
		MaxPageNumber: 4,
		Query:         "Gotham City",
		Language:      giphy.LangChineseTraditional,
	})
	if err != nil {
		log.Fatal(err)
	}

	for page := range res.Pages {
		if page.Err != nil {
			fmt.Printf("#%d: err: %v\n", page.PageNumber, page.Err)
			continue
		}

		fmt.Printf("PageNumber: %#d\n", page.PageNumber)
		for _, giph := range page.Giphs {
			fmt.Printf("\tDescription: %v\n", giph)
			for size, detail := range giph.Sizes {
				fmt.Printf("\t\t%s: %#v\n", size, detail)
			}
		}
		fmt.Printf("\n\n")
	}
}
