package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

func SendJSON(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": statusCode,
		"status":      status,
		"message":     message,
		"data":        data,
	})
}

func GetCoordinates(address string) (float64, float64, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&addressdetails=1", url.QueryEscape(address))
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "HospitalApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var results []struct {
		Lat     string `json:"lat"`
		Lon     string `json:"lon"`
		Address struct {
			Postcode string `json:"postcode"`
		} `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil || len(results) == 0 {
		return 0, 0, "", fmt.Errorf("location not found")
	}
	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)
	return lat, lon, results[0].Address.Postcode, nil
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
