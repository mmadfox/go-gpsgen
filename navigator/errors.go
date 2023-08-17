package navigator

import "errors"

// ErrNoRoutes indicates that no routes are available.
var ErrNoRoutes = errors.New("gpsgen/navigator: no routes")

// ErrTrackNotFound indicates that a specified track was not found.
var ErrTrackNotFound = errors.New("gpsgen/navigator: track not found")

// ErrRouteNotFound indicates that a specified route was not found.
var ErrRouteNotFound = errors.New("gpsgen/navigator: route not found")

// ErrInvalidRoutePath indicates that a route has an invalid or too short path.
var ErrInvalidRoutePath = errors.New("gpsgen/navigator: route is too short")
