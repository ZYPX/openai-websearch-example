package main

import (
	"fmt"
	"regexp"
	"strings"
)

func replaceTags(htmlString *string, elm string) {
	re, err := regexp.Compile(fmt.Sprintf(`</?%s[^>]*>`, elm))

	if err != nil {
		fmt.Println("Regex compilation error:", err)
		return
	}

	*htmlString = re.ReplaceAllString(*htmlString, "\n")
}

func removeTags(htmlString *string, tag string) {

	scriptRe, err := regexp.Compile(fmt.Sprintf(`(?is)<\s*%s\b[^>]*>.*?</\s*%s\s*>`, tag, tag))

	if err != nil {
		fmt.Println("Regex compilation error:", err)
		return
	}

	*htmlString = scriptRe.ReplaceAllString(*htmlString, "")
}

func removeWhitespace(htmlString *string) {
	re, err := regexp.Compile(`\s{4,}`)

	if err != nil {
		fmt.Println("Regex compilation error:", err)
		return
	}

	*htmlString = re.ReplaceAllString(*htmlString, "")
}

func extractInnerHTML(htmlString *string) {
	// Replace <br> tags with newlines
	brRe := regexp.MustCompile(`<br\s*/?>\s*`)
	*htmlString = brRe.ReplaceAllString(*htmlString, "\n")

	// Remove remaining HTML tags
	tagRe := regexp.MustCompile(`<[^>]+>`)
	*htmlString = tagRe.ReplaceAllString(*htmlString, "")

	// Trim spaces and newlines
	*htmlString = strings.TrimSpace(*htmlString)

	// Replace multiple newlines one newline
	newlineRe := regexp.MustCompile(`[ \t\n\f\r]{2,}`)
	*htmlString = newlineRe.ReplaceAllString(*htmlString, "\n")

	multiSpacesRe := regexp.MustCompile(` +`)
	*htmlString = multiSpacesRe.ReplaceAllString(*htmlString, " ")

}

func processHTMLString(removeTag []string, replaceTextTag []string, extractTextOnly bool, htmlString *string) *string {

	if len(removeTag) > 0 {
		for _, tag := range removeTag {
			removeTags(htmlString, tag)
		}
	}

	if len(replaceTextTag) > 0 {
		for _, tag := range replaceTextTag {
			replaceTags(htmlString, tag)
		}
	}

	removeWhitespace(htmlString)

	if extractTextOnly {
		extractInnerHTML(htmlString)
	}

	var htmlText = htmlString

	return htmlText
}
