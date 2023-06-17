package main

import (
	"flag"
	"fmt"
	api "http-login-packaged/pkg/api"
	"net/url"
	"os"
)

func main() {

	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)
	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.StringVar(&password, "password", "", "Password for access to API")
	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Validation Error, URL is not valid %s \n usage -h", err)
		flag.Usage()
		os.Exit(1)
	}
	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedURL.String())
	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("Error is :%s { HTTP Code : %d, Body : %s}\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error :%s \n ", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Println("No Response")
		os.Exit(1)
	}

	fmt.Print(res.GetResponse())
}
