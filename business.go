package yfusion

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// CategoriesInfo - Category data returned on Business data
type CategoriesInfo struct {
	Alias string
	Title string
}

// Coords - Latitude and Longitude data
type Coords struct {
	Latitude  float64
	Longitude float64
}

// Loc - Location information including Address, City, Country, DisplayAddresses, State, and ZipCode
type Loc struct {
	Address1       string
	Address2       string
	Address3       string
	City           string
	Country        string
	DisplayAddress []string `json:"display_address"`
	State          string
	ZipCode        string `json:"zip_code"`
	CrossStreets   string `json:"cross_streets"`
}

// BusinessMigratedError - Represents the error returned from an HTTP 301 response on certain requests
type BusinessMigratedError struct {
	Code          string
	Description   string
	NewBusinessID string `json:"new_business_id"`
}

// GeneralBusinessInfo - Data about a Business
type GeneralBusinessInfo struct {
	Categories  []CategoriesInfo
	Coordinates *Coords
	// DisplayPhone is a user friendly version of the phone number to display
	DisplayPhone string `json:"display_phone"`
	Distance     float64
	ID           string
	Alias        string
	ImageURL     string `json:"image_url"`
	// IsClosed - Whether a business has been permanently closed
	IsClosed bool `json:"is_closed"`
	Location *Loc
	Name     string
	// Price - value is either $, $$, $$$, or $$$$
	Price        string
	Rating       float64
	ReviewCount  int `json:"review_count"`
	URL          string
	Transactions []string
}

// DetailedBusinessInfo - Data returned from a Business Details request
type DetailedBusinessInfo struct {
	GeneralBusinessInfo
	Phone        string
	Photos       []string
	Hours        []HoursInfo
	Transactions []string
	IsClaimed    bool `json:"is_claimed"`
	// Attributes is only visible for Yelp Fusion VIP clients
	Attributes map[string]interface{}
	Error      *BusinessMigratedError
}

// HoursInfo - Data about hours for the business
type HoursInfo struct {
	HoursType string `json:"hours_type"`
	Open      []OpenInfo
	IsOpenNow bool `json:"is_open_now"`
}

// OpenInfo - Data about when the business is open
type OpenInfo struct {
	IsOvernight bool `json:"is_overnight"`
	// End and Start are 24 hour clocks
	End   string
	Day   int
	Start string
}

// BusinessSearchData - The data returned from the Business Search route
type BusinessSearchData struct {
	Total      int
	Businesses []GeneralBusinessInfo
	Region     map[string]interface{}
}

// BusinessSearchParams - Options to use when sending a request to the Business Search route.
// Location is Mandatory if Latitude and Longitude are not specified.
// Latitude and Longitude are required if Location is not specified.
// All other fields are optional.
type BusinessSearchParams struct {
	Term       *string
	Location   *string
	Latitude   *float64
	Longitude  *float64
	Radius     *int
	Categories *string
	Locale     *string
	Limit      *int
	Offset     *int
	SortBy     *string
	Price      *string
	OpenNow    *bool
	OpenAt     *int
	Attributes *string
}

func getLocOrLatLong(bus *BusinessSearchParams) (string, error) {
	if bus.Location == nil && bus.Latitude == nil && bus.Longitude == nil {
		return "", errors.New("error missing required fields: Location or (Latitude and Longitude)")
	}
	sb := &strings.Builder{}
	if bus.Location != nil {
		sb.WriteString(fmt.Sprintf("location=%s", url.QueryEscape(*bus.Location)))
	}
	if bus.Latitude != nil && bus.Longitude != nil {
		if sb.Len() > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(fmt.Sprintf("latitude=%f&longitude=%f", *bus.Latitude, *bus.Longitude))
	}
	return sb.String(), nil
}

func getOpenAtOrNow(bus *BusinessSearchParams) (string, error) {
	if bus.OpenAt != nil && bus.OpenNow != nil {
		return "", errors.New("cannot set both open_at and open_now parameters")
	}
	if bus.OpenNow != nil {
		return fmt.Sprintf("open_now=%v", *bus.OpenNow), nil
	}
	if bus.OpenAt != nil {
		return fmt.Sprintf("open_at=%d", *bus.OpenAt), nil
	}
	return "", nil
}

// Params - Return the set BusinessSearchParams fields in a query param string
func (bs *BusinessSearchParams) Params() (string, error) {
	sb := &strings.Builder{}
	locString, err := getLocOrLatLong(bs)
	if err != nil {
		return "", err
	}
	sb.WriteString(locString)
	if bs.Term != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("term=%s", url.QueryEscape(*bs.Term)))
	}
	if bs.Radius != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("radius=%d", *bs.Radius))
	}
	if bs.Categories != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("categories=%s", url.QueryEscape(*bs.Categories)))
	}
	if bs.Locale != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("locale=%s", url.QueryEscape(*bs.Locale)))
	}
	if bs.Limit != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("limit=%d", *bs.Limit))
	}
	if bs.Offset != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("offset=%d", *bs.Offset))
	}
	if bs.SortBy != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("sort_by=%s", url.QueryEscape(*bs.SortBy)))
	}
	if bs.Price != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("price=%s", url.QueryEscape(*bs.Price)))
	}
	open, err := getOpenAtOrNow(bs)
	if err != nil {
		return "", err
	}
	if open != "" {
		sb.WriteString("&")
		sb.WriteString(open)
	}
	if bs.Attributes != nil {
		sb.WriteString("&")
		sb.WriteString(fmt.Sprintf("attributes=%s", url.QueryEscape(*bs.Attributes)))
	}
	return sb.String(), nil
}

// SetTerm - Set the terms to query for
func (bs *BusinessSearchParams) SetTerm(s string) {
	bs.Term = new(string)
	*bs.Term = s
}

// SetLocation - Set the location to focus on
func (bs *BusinessSearchParams) SetLocation(s string) {
	bs.Location = new(string)
	*bs.Location = s
}

// SetLatitude - Set the latitude to query for
func (bs *BusinessSearchParams) SetLatitude(i float64) {
	bs.Latitude = new(float64)
	*bs.Latitude = i
}

// SetLongitude - Set the longitude to query for
func (bs *BusinessSearchParams) SetLongitude(i float64) {
	bs.Longitude = new(float64)
	*bs.Longitude = i
}

// SetRadius - Set how wide the search radius should be.
//
// The max is 40000 meters or about 25 miles.
func (bs *BusinessSearchParams) SetRadius(i int) {
	bs.Radius = new(int)
	*bs.Radius = i
}

// SetCategories - Set the categories to filter on
//
// Get full list of categories here: https://www.yelp.com/developers/documentation/v3/all_category_list
func (bs *BusinessSearchParams) SetCategories(s string) {
	bs.Categories = new(string)
	*bs.Categories = s
}

// SetLocale - Set the locale.
func (bs *BusinessSearchParams) SetLocale(s string) {
	bs.Locale = new(string)
	*bs.Locale = s
}

// SetLimit - Set the limit of returned businesses
func (bs *BusinessSearchParams) SetLimit(i int) {
	bs.Limit = new(int)
	*bs.Limit = i
}

// SetOffset - Set the offset starting point in the list of businesses
func (bs *BusinessSearchParams) SetOffset(i int) {
	bs.Offset = new(int)
	*bs.Offset = i
}

// SetSortBy - Set how the return values should be sorted
//
// options:
// - best_match
// - rating
// - review_count
// - distance
func (bs *BusinessSearchParams) SetSortBy(s string) {
	bs.SortBy = new(string)
	*bs.SortBy = s
}

// SetPrice - Set the price to filter on
//
// options:
// - 1 = $
// - 2 = $$
// - 3 = $$$
// - 4 = $$$$
//
// You can also combine like so: "1, 2, 3" which filters results for "$, $$, $$$"
func (bs *BusinessSearchParams) SetPrice(s string) {
	bs.Price = new(string)
	*bs.Price = s
}

// SetOpenNow - Set the flag to only show if businesses are open now.
// True means to only show businesses that are open currently
// False means to show all businesses
func (bs *BusinessSearchParams) SetOpenNow(b bool) {
	bs.OpenNow = new(bool)
	*bs.OpenNow = b
}

// SetOpenAt - Set the filter to show business that open at a certain time
func (bs *BusinessSearchParams) SetOpenAt(i int) {
	bs.OpenAt = new(int)
	*bs.OpenAt = i
}

// SetAttributes - Set additional attributes to filter on.
//
// options:
// - hot_and_new (popular businesses that recently joined yelp)
// - request_a_quote (businesses that actively reply to request a quote inquries)
// - reservation (businesses with yelp reservations)
// - waitlist_reservation (businesses with yelp waitlist bookings)
// - cashback (businesses that offer yelp cash back)
// - deals (businesses offering yelp deals)
// - gender_neutral_restrooms (businesses with gender netural restrooms)
// - open_to_all (businesses which are open to all)
func (bs *BusinessSearchParams) SetAttributes(s string) {
	bs.Attributes = new(string)
	*bs.Attributes = s
}
