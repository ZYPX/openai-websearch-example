package main

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"log"
	"net/url"
	"regexp"
)

func getPage(endpoint string) (responseBody string, err error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header = http.Header{
		"accept": {"text/html,application/xhtml+xml"},
		//"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cache-control":             {"no-cache"},
		"sec-ch-ua":                 {`"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"empty"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-site":            {"same-origin"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"cache-control",
			"if-none-match",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"sec-fetch-user",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(readBytes), nil
}

func searchWeb(query string) (results []string, err error) {
	queryURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	var res = []string{"", ""}

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL, nil)
	if err != nil {
		log.Println(err)
		return res, err
	}

	req.Header = http.Header{
		"accept": {"text/html,application/xhtml+xml"},
		//"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cache-control":             {"no-cache"},
		"sec-ch-ua":                 {`"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"empty"},
		"sec-fetch-mode":            {"cors"},
		"sec-fetch-site":            {"same-origin"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"cache-control",
			"if-none-match",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"sec-fetch-user",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res, err
	}

	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return res, err
	}

	htmlString := string(readBytes)

	pattern := `<a\s+class="result__url"\s+href="([^"]+)"`

	// Compile the regex
	re := regexp.MustCompile(pattern)

	// Find the first match
	matches := re.FindAllStringSubmatch(htmlString, -1)

	//fmt.Println(matches)
	var queryResults []string
	var i uint8 = 0
	for _, match := range matches {
		if i == 10 {
			break
		}
		href := match[1]

		// Parse the URL to extract the query parameter 'uddg'
		u, err := url.Parse(href)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return res, err
		}

		// Extract the 'uddg' query parameter
		uddg := u.Query().Get("uddg")
		if uddg == "" {
			fmt.Println("No 'uddg' parameter found")
			return res, err
		}

		// Decode the URL
		decodedURL, err := url.QueryUnescape(uddg)
		if err != nil {
			fmt.Println("Error decoding URL:", err)
			return res, err
		}
		queryResults = append(queryResults, decodedURL)
		i++
	}

	s, err := getPage(queryResults[0])
	if err != nil {
		log.Fatal(err)
	}

	tags := []string{"head", "footer", "style", "script", "header", "nav", "navbar"}

	blockElements := []string{"div", "p", "h1", "h2", "h3", "h4", "h5", "h6", "pre", "blockquote", "span"}

	finalOutput := processHTMLString(tags, blockElements, true, &s)

	parsedURL, err := url.Parse(queryResults[0])
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}

	finalResult := []string{*finalOutput, parsedURL.String()}

	return finalResult, nil
}
