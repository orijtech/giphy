# giphy
Giphy API client in Go

## Samples
### Preamble
```go
package main

import (
	"fmt"
	"log"

	"github.com/orijtech/giphy/v1"
)
```

* Trending gifs
```go
func latestTrending() {
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
```

* Search for gifs
```go
func searchForGIFS() {
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
```

* Random GIF
```go
func randomGIF() {
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
```

* GIF by ID
```go
func gifByID() {
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
```

* Trending stickers
```go
func trendingStickers() {
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
```

* Search for stickers
```go
func searchStickers() {
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
```
