# yfusion
Golang client for Yelp Fusion API. https://www.yelp.com/fusion

## Current Status

Currently this library has limited support for the API. It started out as just a need for a couple of the requests but now I am
trying to flesh out the rest of the API. 

This library is just a wrapper to make interacting with the Yelp Fusion API easier and to provide returned structs with defined
fields instead of having just a map object.

[Current Docs](https://github.com/iCurlmyster/yfusion/wiki/Docs)

Supported requests:
- Business Search
- Business Details
- Business Search by Phone number
- Business Reviews

Requests TODO:
- Business Match
- Transaction Search
- Autocomplete

## Requirements

- Go v1.10+

## Install

```bash
go get github.com/jmatth11/yfusion
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

