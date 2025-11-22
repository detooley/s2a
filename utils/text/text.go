package text

import (
	"regexp"
)

// Return a string between two delimiters (stolen from AI)
func ExtractBetween(s, start, end string) string {
	re := regexp.MustCompile(regexp.QuoteMeta(start) + "(.*?)" + regexp.QuoteMeta(end))
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func RemoveBrackets(text string) string {
	// Regex engines
	rb := regexp.MustCompile(`\[.*?\]`)
	rt := regexp.MustCompile(`[[:punct:]\s]+$`)
	// Remove brackets
	brac := rb.ReplaceAllString(text, "")
	// Remove trailing white space and punctuation
	brac = rt.ReplaceAllString(brac, "")
	return brac
}

func SqlSanitize(text string) string {
	// Remove single quotes
	text = regexp.MustCompile(`'`).ReplaceAllString(text, "''")
	// Remove double quotes
	text = regexp.MustCompile(`"`).ReplaceAllString(text, "")
	// Remove backslashes
	text = regexp.MustCompile(`\\`).ReplaceAllString(text, "")
	// Remove semicolons
	text = regexp.MustCompile(`;`).ReplaceAllString(text, "")
	return text
}

func SearchSanitize(text string) string {
	// Remove single quotes
	text = regexp.MustCompile(`'`).ReplaceAllString(text, "")
	// Remove double quotes
	// text = regexp.MustCompile(`"`).ReplaceAllString(text, "")
	// Remove backslashes
	text = regexp.MustCompile(`\\`).ReplaceAllString(text, "")
	// Remove semicolons
	text = regexp.MustCompile(`;`).ReplaceAllString(text, "")
	return text
}

func StripHTMLTags(htmlString string) string {
	// Matches any HTML tag
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(htmlString, "")
}
