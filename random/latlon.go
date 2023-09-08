package random

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/mmadfox/go-gpsgen/geo"
)

var (
	// ErrCountryBoundsNotFound returns when bbox for country is not found.
	ErrCountryBoundsNotFound = errors.New("gpsgen/latlon: country bounds not found")

	// ErrCountryNotFound returns when country not found.
	ErrCountryNotFound = errors.New("gpsgen/latlon: country not found")
)

var countries = map[string]string{
	"AF": "Afghanistan",
	"AO": "Angola",
	"AL": "Albania",
	"AE": "United Arab Emirates",
	"AR": "Argentina",
	"AM": "Armenia",
	"AQ": "Antarctica",
	"TF": "French Southern Territories",
	"AU": "Australia",
	"AT": "Austria",
	"AZ": "Azerbaijan",
	"BI": "Burundi",
	"BE": "Belgium",
	"BJ": "Benin",
	"BF": "Burkina Faso",
	"BD": "Bangladesh",
	"BG": "Bulgaria",
	"BS": "Bahamas",
	"BA": "Bosnia and Herzegovina",
	"BY": "Belarus",
	"BZ": "Belize",
	"BO": "Bolivia",
	"BR": "Brazil",
	"BN": "Brunei",
	"BT": "Bhutan",
	"BW": "Botswana",
	"CF": "Central African Republic",
	"CA": "Canada",
	"CH": "Switzerland",
	"CL": "Chile",
	"CN": "China",
	"CI": "Ivory Coast",
	"CM": "Cameroon",
	"CD": "Congo (Kinshasa)",
	"CG": "Congo (Brazzaville)",
	"CO": "Colombia",
	"CR": "Costa Rica",
	"CU": "Cuba",
	"CY": "Cyprus",
	"CZ": "Czech Republic",
	"DE": "Germany",
	"DJ": "Djibouti",
	"DK": "Denmark",
	"DO": "Dominican Republic",
	"DZ": "Algeria",
	"EC": "Ecuador",
	"EG": "Egypt",
	"ER": "Eritrea",
	"ES": "Spain",
	"EE": "Estonia",
	"ET": "Ethiopia",
	"FI": "Finland",
	"FJ": "Fiji",
	"FK": "Falkland Islands",
	"FR": "France",
	"GA": "Gabon",
	"GB": "United Kingdom",
	"GE": "Georgia",
	"GH": "Ghana",
	"GN": "Guinea",
	"GM": "Gambia",
	"GW": "Guinea Bissau",
	"GQ": "Equatoral Guinea",
	"GR": "Greece",
	"GL": "Greenland",
	"GT": "Guatemala",
	"GY": "Guyana",
	"HN": "Honduras",
	"HR": "Croatia",
	"HT": "Haiti",
	"HU": "Hungary",
	"ID": "Indonesia",
	"IN": "India",
	"IE": "Ireland",
	"IR": "Iran",
	"IQ": "Iraq",
	"IS": "Iceland",
	"IL": "Israel",
	"IT": "Italy",
	"JM": "Jamaica",
	"JO": "Jordan",
	"JP": "Japan",
	"KZ": "Kazakhstan",
	"KE": "Kenya",
	"KG": "Kyrgyzstan",
	"KH": "Cambodia",
	"KR": "South Korea",
	"KW": "Kuwait",
	"LA": "Laos",
	"LB": "Lebanon",
	"LR": "Liberia",
	"LY": "Libya",
	"LK": "Sri Lanka",
	"LS": "Lesotho",
	"LT": "Lithuania",
	"LU": "Luxembourg",
	"LV": "Latvia",
	"MA": "Morocco",
	"MD": "Moldova",
	"MG": "Madagascar",
	"MX": "Mexico",
	"MK": "Macedonia",
	"ML": "Mali",
	"MM": "Myanmar",
	"ME": "Montenegro",
	"MN": "Mongolia",
	"MZ": "Mozambique",
	"MR": "Mauritania",
	"MW": "Malawi",
	"MY": "Malaysia",
	"NA": "Namibia",
	"NC": "New Caledonia",
	"NE": "Niger",
	"NG": "Nigeria",
	"NI": "Nicaragua",
	"NL": "Netherlands",
	"NO": "Norway",
	"NP": "Nepal",
	"NZ": "New Zealand",
	"OM": "Oman",
	"PK": "Pakistan",
	"PA": "Panama",
	"PE": "Peru",
	"PH": "Philippines",
	"PG": "Papua New Guinea",
	"PL": "Poland",
	"PR": "Puerto Rico",
	"KP": "North Korea",
	"PT": "Portugal",
	"PY": "Paraguay",
	"QA": "Qatar",
	"RO": "Romania",
	"RU": "Russia",
	"RW": "Rwanda",
	"SA": "Saudi Arabia",
	"SD": "Sudan",
	"SS": "South Sudan",
	"SN": "Senegal",
	"SB": "Solomon Islands",
	"SL": "Sierra Leone",
	"SV": "El Salvador",
	"SO": "Somalia",
	"RS": "Serbia",
	"SR": "Suriname",
	"SK": "Slovakia",
	"SI": "Slovenia",
	"SE": "Sweden",
	"SZ": "Swaziland",
	"SY": "Syria",
	"TD": "Chad",
	"TG": "Togo",
	"TH": "Thailand",
	"TJ": "Tajikistan",
	"TM": "Turkmenistan",
	"TL": "East Timor",
	"TT": "Trinidad and Tobago",
	"TN": "Tunisia",
	"TR": "Turkey",
	"TW": "Taiwan",
	"TZ": "Tanzania",
	"UG": "Uganda",
	"UA": "Ukraine",
	"UY": "Uruguay",
	"US": "United States",
	"UZ": "Uzbekistan",
	"VE": "Venezuela",
	"VN": "Vietnam",
	"VU": "Vanuatu",
	"PS": "West Bank",
	"YE": "Yemen",
	"ZA": "South Africa",
	"ZM": "Zambia",
	"ZW": "Zimbabwe",
}

var bboxByCountries = map[string][4]float64{
	"AF": {60.53, 29.32, 75.16, 38.49},
	"AO": {11.64, -17.93, 24.08, -4.44},
	"AL": {19.3, 39.62, 21.02, 42.69},
	"AE": {51.58, 22.5, 56.4, 26.06},
	"AR": {-73.42, -55.25, -53.63, -21.83},
	"AM": {43.58, 38.74, 46.51, 41.25},
	"AQ": {-180.0, -90.0, 180.0, -63.27},
	"TF": {68.72, -49.78, 70.56, -48.63},
	"AU": {113.34, -43.63, 153.57, -10.67},
	"AT": {9.48, 46.43, 16.98, 49.04},
	"AZ": {44.79, 38.27, 50.39, 41.86},
	"BI": {29.02, -4.5, 30.75, -2.35},
	"BE": {2.51, 49.53, 6.16, 51.48},
	"BJ": {0.77, 6.14, 3.8, 12.24},
	"BF": {-5.47, 9.61, 2.18, 15.12},
	"BD": {88.08, 20.67, 92.67, 26.45},
	"BG": {22.38, 41.23, 28.56, 44.23},
	"BS": {-78.98, 23.71, -77.0, 27.04},
	"BA": {15.75, 42.65, 19.6, 45.23},
	"BY": {23.2, 51.32, 32.69, 56.17},
	"BZ": {-89.23, 15.89, -88.11, 18.5},
	"BO": {-69.59, -22.87, -57.5, -9.76},
	"BR": {-73.99, -33.77, -34.73, 5.24},
	"BN": {114.2, 4.01, 115.45, 5.45},
	"BT": {88.81, 26.72, 92.1, 28.3},
	"BW": {19.9, -26.83, 29.43, -17.66},
	"CF": {14.46, 2.27, 27.37, 11.14},
	"CA": {-141.0, 41.68, -52.65, 73.23},
	"CH": {6.02, 45.78, 10.44, 47.83},
	"CL": {-75.64, -55.61, -66.96, -17.58},
	"CN": {73.68, 18.2, 135.03, 53.46},
	"CI": {-8.6, 4.34, -2.56, 10.52},
	"CM": {8.49, 1.73, 16.01, 12.86},
	"CD": {12.18, -13.26, 31.17, 5.26},
	"CG": {11.09, -5.04, 18.45, 3.73},
	"CO": {-78.99, -4.3, -66.88, 12.44},
	"CR": {-85.94, 8.23, -82.55, 11.22},
	"CU": {-84.97, 19.86, -74.18, 23.19},
	"CY": {32.26, 34.57, 34.0, 35.17},
	"CZ": {12.24, 48.56, 18.85, 51.12},
	"DE": {5.99, 47.3, 15.02, 54.98},
	"DJ": {41.66, 10.93, 43.32, 12.7},
	"DK": {8.09, 54.8, 12.69, 57.73},
	"DO": {-71.95, 17.6, -68.32, 19.88},
	"DZ": {-8.68, 19.06, 12.0, 37.12},
	"EC": {-80.97, -4.96, -75.23, 1.38},
	"EG": {24.7, 22.0, 36.87, 31.59},
	"ER": {36.32, 12.46, 43.08, 18.0},
	"ES": {-9.39, 35.95, 3.04, 43.75},
	"EE": {23.34, 57.47, 28.13, 59.61},
	"ET": {32.95, 3.42, 47.79, 14.96},
	"FI": {20.65, 59.85, 31.52, 70.16},
	"FJ": {-180.0, -18.29, 180.0, -16.02},
	"FK": {-61.2, -52.3, -57.75, -51.1},
	"FR": {-5.0, 42.5, 9.56, 51.15},
	"GA": {8.8, -3.98, 14.43, 2.33},
	"GB": {-7.57, 49.96, 1.68, 58.64},
	"GE": {39.96, 41.06, 46.64, 43.55},
	"GH": {-3.24, 4.71, 1.06, 11.1},
	"GN": {-15.13, 7.31, -7.83, 12.59},
	"GM": {-16.84, 13.13, -13.84, 13.88},
	"GW": {-16.68, 11.04, -13.7, 12.63},
	"GQ": {9.31, 1.01, 11.29, 2.28},
	"GR": {20.15, 34.92, 26.6, 41.83},
	"GL": {-73.3, 60.04, -12.21, 83.65},
	"GT": {-92.23, 13.74, -88.23, 17.82},
	"GY": {-61.41, 1.27, -56.54, 8.37},
	"HN": {-89.35, 12.98, -83.15, 16.01},
	"HR": {13.66, 42.48, 19.39, 46.5},
	"HT": {-74.46, 18.03, -71.62, 19.92},
	"HU": {16.2, 45.76, 22.71, 48.62},
	"ID": {95.29, -10.36, 141.03, 5.48},
	"IN": {68.18, 7.97, 97.4, 35.49},
	"IE": {-9.98, 51.67, -6.03, 55.13},
	"IR": {44.11, 25.08, 63.32, 39.71},
	"IQ": {38.79, 29.1, 48.57, 37.39},
	"IS": {-24.33, 63.5, -13.61, 66.53},
	"IL": {34.27, 29.5, 35.84, 33.28},
	"IT": {6.75, 36.62, 18.48, 47.12},
	"JM": {-78.34, 17.7, -76.2, 18.52},
	"JO": {34.92, 29.2, 39.2, 33.38},
	"JP": {129.41, 31.03, 145.54, 45.55},
	"KZ": {46.47, 40.66, 87.36, 55.39},
	"KE": {33.89, -4.68, 41.86, 5.51},
	"KG": {69.46, 39.28, 80.26, 43.3},
	"KH": {102.35, 10.49, 107.61, 14.57},
	"KR": {126.12, 34.39, 129.47, 38.61},
	"KW": {46.57, 28.53, 48.42, 30.06},
	"LA": {100.12, 13.88, 107.56, 22.46},
	"LB": {35.13, 33.09, 36.61, 34.64},
	"LR": {-11.44, 4.36, -7.54, 8.54},
	"LY": {9.32, 19.58, 25.16, 33.14},
	"LK": {79.7, 5.97, 81.79, 9.82},
	"LS": {27.0, -30.65, 29.33, -28.65},
	"LT": {21.06, 53.91, 26.59, 56.37},
	"LU": {5.67, 49.44, 6.24, 50.13},
	"LV": {21.06, 55.62, 28.18, 57.97},
	"MA": {-17.02, 21.42, -1.12, 35.76},
	"MD": {26.62, 45.49, 30.02, 48.47},
	"MG": {43.25, -25.6, 50.48, -12.04},
	"MX": {-117.13, 14.54, -86.81, 32.72},
	"MK": {20.46, 40.84, 22.95, 42.32},
	"ML": {-12.17, 10.1, 4.27, 24.97},
	"MM": {92.3, 9.93, 101.18, 28.34},
	"ME": {18.45, 41.88, 20.34, 43.52},
	"MN": {87.75, 41.6, 119.77, 52.05},
	"MZ": {30.18, -26.74, 40.78, -10.32},
	"MR": {-17.06, 14.62, -4.92, 27.4},
	"MW": {32.69, -16.8, 35.77, -9.23},
	"MY": {100.09, 0.77, 119.18, 6.93},
	"NA": {11.73, -29.05, 25.08, -16.94},
	"NC": {164.03, -22.4, 167.12, -20.11},
	"NE": {0.3, 11.66, 15.9, 23.47},
	"NG": {2.69, 4.24, 14.58, 13.87},
	"NI": {-87.67, 10.73, -83.15, 15.02},
	"NL": {3.31, 50.8, 7.09, 53.51},
	"NO": {4.99, 58.08, 31.29, 70.92},
	"NP": {80.09, 26.4, 88.17, 30.42},
	"NZ": {166.51, -46.64, 178.52, -34.45},
	"OM": {52.0, 16.65, 59.81, 26.4},
	"PK": {60.87, 23.69, 77.84, 37.13},
	"PA": {-82.97, 7.22, -77.24, 9.61},
	"PE": {-81.41, -18.35, -68.67, -0.06},
	"PH": {117.17, 5.58, 126.54, 18.51},
	"PG": {141.0, -10.65, 156.02, -2.5},
	"PL": {14.07, 49.03, 24.03, 54.85},
	"PR": {-67.24, 17.95, -65.59, 18.52},
	"KP": {124.27, 37.67, 130.78, 42.99},
	"PT": {-9.53, 36.84, -6.39, 42.28},
	"PY": {-62.69, -27.55, -54.29, -19.34},
	"QA": {50.74, 24.56, 51.61, 26.11},
	"RO": {20.22, 43.69, 29.63, 48.22},
	"RU": {-180.0, 41.15, 180.0, 81.25},
	"RW": {29.02, -2.92, 30.82, -1.13},
	"SA": {34.63, 16.35, 55.67, 32.16},
	"SD": {21.94, 8.62, 38.41, 22.0},
	"SS": {23.89, 3.51, 35.3, 12.25},
	"SN": {-17.63, 12.33, -11.47, 16.6},
	"SB": {156.49, -10.83, 162.4, -6.6},
	"SL": {-13.25, 6.79, -10.23, 10.05},
	"SV": {-90.1, 13.15, -87.72, 14.42},
	"SO": {40.98, -1.68, 51.13, 12.02},
	"RS": {18.83, 42.25, 22.99, 46.17},
	"SR": {-58.04, 1.82, -53.96, 6.03},
	"SK": {16.88, 47.76, 22.56, 49.57},
	"SI": {13.7, 45.45, 16.56, 46.85},
	"SE": {11.03, 55.36, 23.9, 69.11},
	"SZ": {30.68, -27.29, 32.07, -25.66},
	"SY": {35.7, 32.31, 42.35, 37.23},
	"TD": {13.54, 7.42, 23.89, 23.41},
	"TG": {-0.05, 5.93, 1.87, 11.02},
	"TH": {97.38, 5.69, 105.59, 20.42},
	"TJ": {67.44, 36.74, 74.98, 40.96},
	"TM": {52.5, 35.27, 66.55, 42.75},
	"TL": {124.97, -9.39, 127.34, -8.27},
	"TT": {-61.95, 10.0, -60.9, 10.89},
	"TN": {7.52, 30.31, 11.49, 37.35},
	"TR": {26.04, 35.82, 44.79, 42.14},
	"TW": {120.11, 21.97, 121.95, 25.3},
	"TZ": {29.34, -11.72, 40.32, -0.95},
	"UG": {29.58, -1.44, 35.04, 4.25},
	"UA": {22.09, 44.36, 40.08, 52.34},
	"UY": {-58.43, -34.95, -53.21, -30.11},
	"US": {-125.0, 25.0, -66.96, 49.5},
	"UZ": {55.93, 37.14, 73.06, 45.59},
	"VE": {-73.3, 0.72, -59.76, 12.16},
	"VN": {102.17, 8.6, 109.34, 23.35},
	"VU": {166.63, -16.6, 167.84, -14.63},
	"PS": {34.93, 31.35, 35.55, 32.53},
	"YE": {42.6, 12.59, 53.11, 19.0},
	"ZA": {16.34, -34.82, 32.83, -22.09},
	"ZM": {21.89, -17.96, 33.49, -8.24},
	"ZW": {25.26, -22.27, 32.85, -15.51},
}

// CountryName returns the full name of a country based on its two-letter country code.
func CountryName(code string) (string, error) {
	if len(code) != 2 {
		return "", ErrCountryNotFound
	}
	code = strings.ToUpper(code)
	name, ok := countries[code]
	if !ok {
		return "", ErrCountryNotFound
	}
	return name, nil
}

// EachCountry iterates through each country and executes a given function for each.
func EachCountry(fn func(code, name string)) {
	for code, name := range countries {
		fn(code, name)
	}
}

// BoundingBox retrieves the bounding box of a country based on its two-letter country code.
func BoundingBox(countryCode string) (geo.BBox, error) {
	if len(countryCode) != 2 {
		return geo.BBox{}, ErrCountryBoundsNotFound
	}
	countryCode = strings.ToUpper(countryCode)
	bbox, ok := bboxByCountries[countryCode]
	if !ok {
		return geo.BBox{}, ErrCountryBoundsNotFound
	}
	return geo.BBox{
		MinLon: bbox[0],
		MinLat: bbox[1],
		MaxLon: bbox[2],
		MaxLat: bbox[3],
	}, nil
}

// LatLonByCountry generates a random latitude and longitude within a specified country's bounds.
func LatLonByCountry(countryCode string) (geo.LatLonPoint, error) {
	if len(countryCode) != 2 {
		return geo.LatLonPoint{}, ErrCountryBoundsNotFound
	}
	countryCode = strings.ToUpper(countryCode)
	bbox, err := BoundingBox(countryCode)
	if err != nil {
		return geo.LatLonPoint{}, err
	}
	return geo.LatLonPoint{
		Lon: rand.Float64()*(bbox.MaxLon-bbox.MinLon) + bbox.MinLon,
		Lat: rand.Float64()*(bbox.MaxLat-bbox.MinLat) + bbox.MinLat,
	}, nil
}

// LatLon generates a random latitude and longitude anywhere in the world.
func LatLon() geo.LatLonPoint {
	return geo.LatLonPoint{
		Lat: lat(),
		Lon: lon(),
	}
}

func rnd() float64 {
	return rand.Float64() - 0.5
}

func lon() float64 {
	return rnd() * 360
}

func lat() float64 {
	return rnd() * 180
}
