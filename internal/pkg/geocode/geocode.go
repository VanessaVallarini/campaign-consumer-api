package geocode

import "math"

// Function to validate whether the region is within a 5 km radius of the address
func IsWithinRegion(addressLat, addressLon, regionLat, regionLon float64) bool {
	distance := haversine(addressLat, addressLon, regionLat, regionLon)
	const radiusKm = 5.0 // 5 kilometer radius
	return distance <= radiusKm
}

// Function to calculate the distance between two points (latitude and longitude) in kilometers
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Difference between coordinates
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Returns the distance in kilometers
	return R * c
}
