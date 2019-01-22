package yfusion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURL      = "https://api.yelp.com/v3"
	busDetails   = "/businesses"
	busSearch    = busDetails + "/search"
	phoneSearch  = busSearch + "/phone"
	reviewSearch = "/reviews"
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
	return yf.SearchBusinessWithContext(nil, bus)
}

// SearchBusinessWithContext - Use the Business Search route with the given BussinesSearchParams options and context
// returns the parsed BusinessSearchData object
func (yf *YelpFusion) SearchBusinessWithContext(ctx context.Context, bus *BusinessSearchParams) (*BusinessSearchData, error) {
	resp, err := yf.SearchBusinessResponse(ctx, bus)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
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
//
// With the option of specifing a locale. (An empty string for locale will leave the parameter off)
// returns the parsed DetailedBusinessInfo object
func (yf *YelpFusion) SearchBusinessDetailsWithLocale(ctx context.Context, busID, locale string) (*DetailedBusinessInfo, error) {
	resp, err := yf.SearchBusinessDetailsWithLocaleResponse(ctx, busID, locale)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	var b *DetailedBusinessInfo
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// SearchBusinessDetailsWithLocaleResponse - Query details about a business, given its ID.
//
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

// SearchBusinessesByPhone - Query Businesses by a phone number.
//
// The phone number must start with a "+" and the country code.
func (yf *YelpFusion) SearchBusinessesByPhone(phoneNumber string) (*BusinessSearchData, error) {
	return yf.SearchBusinessesByPhoneWithContext(nil, phoneNumber)
}

// SearchBusinessesByPhoneWithContext - Query Businesses by a phone number.
//
// The phone number must start with a "+" and the country code.
func (yf *YelpFusion) SearchBusinessesByPhoneWithContext(ctx context.Context, phoneNumber string) (*BusinessSearchData, error) {
	resp, err := yf.SearchBusinessesByPhoneResponse(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	var b *BusinessSearchData
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// SearchBusinessesByPhoneResponse - Query Businesses by a phone number.
//
// The phone number must start with a "+" and the country code.
// Returns the response from the request
func (yf *YelpFusion) SearchBusinessesByPhoneResponse(ctx context.Context, phoneNumber string) (*http.Response, error) {
	if strings.TrimSpace(phoneNumber) == "" {
		return nil, errors.New("phone number is required")
	}
	urlStr := fmt.Sprintf("%s%s?phone=%s", baseURL, phoneSearch, url.QueryEscape(phoneNumber))
	req, err := yf.getRequest(ctx, "GET", urlStr)
	if err != nil {
		return nil, err
	}
	return yf.client.Do(req)
}

// SearchBusinessReviews - Query for reviews for a particular business
// The error field on the ReviewsData object will only be populated if an HTTP 301 status code is returned
// In which case you can resend the request with the NewBusinessID from the error field on the ReviewsData object.
// Otherwise error should be nil.
func (yf *YelpFusion) SearchBusinessReviews(busID string) (*ReviewsData, error) {
	return yf.SearchBusinessReviewsWithLocale(nil, busID, "")
}

// SearchBusinessReviewsWithLocale - Query for reviews for a particular business
// The error field on the ReviewsData object will only be populated if an HTTP 301 status code is returned
// In which case you can resend the request with the NewBusinessID from the error field on the ReviewsData object.
// Otherwise error should be nil.
//
// The locale defaults to en_US if left blank.
func (yf *YelpFusion) SearchBusinessReviewsWithLocale(ctx context.Context, busID, locale string) (*ReviewsData, error) {
	resp, err := yf.SearchBusinessReviewsWithLocaleResponse(ctx, busID, locale)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMovedPermanently {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	var rd *ReviewsData
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&rd); err != nil {
		return nil, err
	}
	return rd, nil
}

// SearchBusinessReviewsWithLocaleResponse - Query for reviews for a particular business
// The error field on the ReviewsData object will only be populated if an HTTP 301 status code is returned
// In which case you can resend the request with the NewBusinessID from the error field on the ReviewsData object.
// Otherwise error should be nil.
//
// The locale defaults to en_US if left blank.
// Returns the response from the request
func (yf *YelpFusion) SearchBusinessReviewsWithLocaleResponse(ctx context.Context, busID, locale string) (*http.Response, error) {
	urlStr := fmt.Sprintf("%s%s/%s%s", baseURL, busDetails, busID, reviewSearch)
	if strings.TrimSpace(locale) != "" {
		urlStr = fmt.Sprintf("%s?locale=%s", urlStr, url.QueryEscape(locale))
	}
	req, err := yf.getRequest(ctx, "GET", urlStr)
	if err != nil {
		return nil, err
	}
	return yf.client.Do(req)
}
