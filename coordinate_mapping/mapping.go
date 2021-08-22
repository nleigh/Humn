package coordinate_mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CoordinateMapper struct {
	ApiToken      string
	OutputChannel chan CoordinatePostcodeOutput
}

func NewCoordinateMapper(apiToken string, outputChan chan CoordinatePostcodeOutput) CoordinateMapper {
	return CoordinateMapper{ApiToken: apiToken, OutputChannel: outputChan}

}
func (cm CoordinateMapper) GetPostcodeDataForCoordinatesAndWriteToOutput(coordinate Coordinate) error {
	mapResponse, err := getCoordinateData(coordinate, cm.ApiToken)
	if err != nil {
		return err
	}

	output, err := processMapData(coordinate, mapResponse)
	if err != nil {
		return err
	}

	cm.OutputChannel <- output
	return nil
}

func processMapData(coordinate Coordinate, mapResponse MapBoxApiResponse) (CoordinatePostcodeOutput, error) {
	// the relevant field you should obtain is the text field from the single returned Feature.
	features := mapResponse.Features

	if len(features) == 0 {
		return CoordinatePostcodeOutput{}, errors.New("error: response object does not contain expected json property 'features'")
	}

	postcode := features[0].Text

	//TODO  sanitize postcode
	// The postcode returned from mapbox should be sanitised to ensure to conforms to the expected postcode format.

	output := CoordinatePostcodeOutput{
		Lat:      coordinate.Latitude,
		Lng:      coordinate.Longitude,
		Postcode: postcode,
	}

	return output, nil
}

func getCoordinateData(coordinate Coordinate, apiToken string) (MapBoxApiResponse, error) { //
	// curl "https://api.mapbox.com/geocoding/v5/mapbox.places/<long,lat>.json?types=postcode&limit=1&access_token=YOUR_MAPBOX_ACCESS_TOKEN"

	url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%f,%f.json?types=postcode&limit=1&access_token=%s",
		coordinate.Longitude,
		coordinate.Latitude,
		apiToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MapBoxApiResponse{}, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return MapBoxApiResponse{}, err
	}

	if response.StatusCode != http.StatusOK {
		return MapBoxApiResponse{}, errors.New("error: response status not 200/OK")
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return MapBoxApiResponse{}, err
	}

	var mapResponse MapBoxApiResponse
	err = json.Unmarshal(bodyBytes, &mapResponse)
	if err != nil {
		return MapBoxApiResponse{}, err
	}

	return mapResponse, nil

}

type Coordinate struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type CoordinatePostcodeOutput struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Postcode string  `json:"postcode"`
}

type MapBoxApiResponse struct {
	Type     string    `json:"type"`
	Query    []float64 `json:"query"`
	Features []struct {
		Id         string   `json:"id"`
		Type       string   `json:"type"`
		PlaceType  []string `json:"place_type"`
		Relevance  int      `json:"relevance"`
		Properties struct {
		} `json:"properties"`
		Text      string    `json:"text"`
		PlaceName string    `json:"place_name"`
		Bbox      []float64 `json:"bbox"`
		Center    []float64 `json:"center"`
		Geometry  struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Context []struct {
			Id        string `json:"id"`
			Wikidata  string `json:"wikidata"`
			Text      string `json:"text"`
			ShortCode string `json:"short_code,omitempty"`
		} `json:"context"`
	} `json:"features"`
	Attribution string `json:"attribution"`
}
