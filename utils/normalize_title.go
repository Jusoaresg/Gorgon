package utils

import "strings"

func NormalizeTitle(title string) string {
	replacer := strings.NewReplacer(
		"-", " ", "_", " ", ".", " ",
		"'", "", "\"", "", ":", "", "!", "", ",", "",
	)
	title = replacer.Replace(title)
	title = strings.ToLower(title)
	title = strings.Join(strings.Fields(title), " ") // remove duplicated spaces
	return title
}
