package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}
type Page struct {
	Name string `json:"page"`
}

type Occurrences struct {
	Words map[string]int `json:"words"`
}
type Words struct {
	PageName string   `json:"page"`
	Input    string   `json:"input"`
	Words    []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ","))
}
func (o Occurrences) GetResponse() string {
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s, %d", word, occurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ","))
}
func (a API) DoGetRequest(requestUrl string) (Response, error) {

	resp, err := a.Client.Get(requestUrl)

	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("Error is :%s { HTTP Code : %d, Body : %s}\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if !json.Valid(body) {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("Error is :%s { HTTP Code : %d, Body : %s}\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}

	}
	if err != nil {
		log.Fatal(err, body)
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}
	switch page.Name {
	case "occurrence":
		var occurrence Occurrences
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{HTTPCode: resp.StatusCode,
				Body: string(body),
				Err:  fmt.Sprintf("unmarshall Error: %s", err),
			}
		}
		fmt.Println("Unmarshalled Data is", occurrence)

		for key, value := range occurrence.Words {
			fmt.Printf("%s :  %d \n", key, value)
		}

		//search a Value
		if values, ok := occurrence.Words["words4"]; ok {
			fmt.Println(values, " is present")
		}
		return occurrence, nil
	case "words":
		var words Words
		fmt.Printf("Status code is %d ,\n and body is %s", resp.StatusCode, body)

		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("Unmarshall Error: %s", err)
		}
		fmt.Println("unmarshalled body is ", words)
		return words, nil
	default:
		return nil, RequestError{HTTPCode: resp.StatusCode,
			Body: string(body),
			Err:  fmt.Sprintf("unmarshall Error: %s", err),
		}

	}

}
