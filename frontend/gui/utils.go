package gui

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pulsone21/threattrack/lib/entities"
)

type backendData[t entities.Entity] struct {
	StatusCode int
	RequestUrl string
	Data       []t
}

func requestData[t entities.Entity](url string) (*backendData[t], error) {
	var data backendData[t]
	fmt.Println("Requesting data from backend", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			fmt.Println("No data found")
			return &data, nil
		}
		fmt.Println("Error with status code", res.StatusCode)
		return nil, fmt.Errorf("error with status code %d from backend", res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
