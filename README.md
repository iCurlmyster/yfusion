# yfusion
Golang client for Yelp Fusion API. https://www.yelp.com/fusion

## Current Status

Supported requests:
- Business Search
- Business Details
- Business Search by Phone number
- Business Reviews

Currently this has limited support of the API, because these are the requests that I needed.
However I do plan to flesh out the rest of the API.

## Requirements

- Go v1.10+

## Install

```bash
go get github.com/iCurlmyster/yfusion
```

## Sample Code

```golang
package main

import (
  "fmt"

  "github.com/iCurlmyster/yfusion"
)

const (
  key = "<api-key>"
)

func main() {
  yelp := yfusion.NewYelpFusion(key)
  bs := &yfusion.BusinessSearchParams{}
  bs.SetLocation("Austin, TX")
  bs.SetTerm("food")
  result, err := yelp.SearchBusiness(bs)
  if err != nil {
    panic(err)
  }
  for _, b := range result.Businesses {
    fmt.Printf("Name: %s, Price: %s, Distance: %f, Rating: %f\n", b.Name, b.Price, b.Distance, b.Rating)
  }
}
```

