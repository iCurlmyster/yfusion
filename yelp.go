package yfusion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURL    = "https://api.yelp.com/v3"
	busDetails = "/businesses"
	busSearch  = busDetails + "/search"
)

// YelpFusion - Object to interact with Yelp's Fusion v3 API
type YelpFusion struct {
	client *http.Client
	apiKey string
}

// NewYelpFusion - Generate a new YelpFusion object with a given API key
func NewYelpFusion(key string) *YelpFusion {
	return &YelpFusion{
		client: http.DefaultClient,
		apiKey: key,
	}
}

// NewYelpFusionWithClient - Generate a new YelpFusion object with a given API key and http client object
func NewYelpFusionWithClient(key string, client *http.Client) *YelpFusion {
	if client == nil {
		client = http.DefaultClient
	}
	return &YelpFusion{
		client: client,
		apiKey: key,
	}
}

func (yf *YelpFusion) getRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", yf.apiKey))
	if ctx != nil {
		req.WithContext(ctx)
	}
	return req, nil
}

// SearchBusiness - Use the Business Search route with the given BussinesSearchParams options
// returns the parsed BusinessSearchData object
func (yf *YelpFusion) SearchBusiness(bus *BusinessSearchParams) (*BusinessSearchData, error) {
	resp, err := yf.SearchBusinessResponse(nil, bus)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var b *BusinessSearchData
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// SearchBusinessWithContext - Use the Business Search route with the given BussinesSearchParams options and context
// returns the parsed BusinessSearchData object
func (yf *YelpFusion) SearchBusinessWithContext(ctx context.Context, bus *BusinessSearchParams) (*BusinessSearchData, error) {
	resp, err := yf.SearchBusinessResponse(ctx, bus)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var b *BusinessSearchData
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// SearchBusinessResponse - Use the Business Search route with the given BusinessSearchParams options
// returns the Response from the request
func (yf *YelpFusion) SearchBusinessResponse(ctx context.Context, bus *BusinessSearchParams) (*http.Response, error) {
	params, err := bus.Params()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s?%s", baseURL, busSearch, params)
	req, err := yf.getRequest(ctx, "GET", urlStr)
	if err != nil {
		return nil, err
	}
	return yf.client.Do(req)
}

// SearchBusinessDetails - Query details about a business, given its ID
// returns the parsed DetailedBusinessInfo object
func (yf *YelpFusion) SearchBusinessDetails(busID string) (*DetailedBusinessInfo, error) {
	return yf.SearchBusinessDetailsWithLocale(nil, busID, "")
}

// SearchBusinessDetailsWithLocale - Query details about a business, given its ID with context.
// With the option of specifing a locale. (An empty string for locale will leave the parameter off)
// returns the parsed DetailedBusinessInfo object
func (yf *YelpFusion) SearchBusinessDetailsWithLocale(ctx context.Context, busID, locale string) (*DetailedBusinessInfo, error) {
	resp, err := yf.SearchBusinessDetailsWithLocaleResponse(ctx, busID, locale)
	if err != nil {
		return nil, err
	}
	var b *DetailedBusinessInfo
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// SearchBusinessDetailsWithLocaleResponse - Query details about a business, given its ID.
// With the option of specifing a locale. (An empty string for locale will leave the parameter off)
// returns the Response from the request
func (yf *YelpFusion) SearchBusinessDetailsWithLocaleResponse(ctx context.Context, busID, locale string) (*http.Response, error) {
	urlStr := fmt.Sprintf("%s%s/%s", baseURL, busDetails, busID)
	if strings.TrimSpace(locale) != "" {
		urlStr = fmt.Sprintf("%s?locale=%s", urlStr, url.QueryEscape(locale))
	}
	req, err := yf.getRequest(ctx, "GET", urlStr)
	if err != nil {
		return nil, err
	}
	return yf.client.Do(req)
}
