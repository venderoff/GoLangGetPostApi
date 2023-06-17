package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, requestURL string, password string) (string, error) {
	loginRequest := LoginRequest{Password: password}
	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshall Error %s", err)
	}

	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("http post Error %s", err)
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Read All Error %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("invalid output HttpCode %d : %s \n", response.StatusCode, string(body))
	}

	if !json.Valid(resBody) {
		return "", RequestError{HTTPCode: response.StatusCode,
			Body: string(resBody),
			Err:  fmt.Sprintf("No Valid Json return %s", resBody),
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{HTTPCode: response.StatusCode,
			Body: string(body),
			Err:  fmt.Sprintf("Page unmarshall Error %s", err)}

	}
	if loginResponse.Token == "" {
		if err != nil {
			return "", RequestError{HTTPCode: response.StatusCode,
				Body: string(resBody),
				Err:  "Empty Token Replied"}
		}
	}
	return loginResponse.Token, nil

}
