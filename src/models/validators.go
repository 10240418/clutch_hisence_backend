package models

import "regexp"

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

func ValidateMobile(mobile string) bool {
	pattern := `^(\+?(\d{1,3}))?(\d{10,12})$`
	return regexp.MustCompile(pattern).MatchString(mobile)
}

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func ValidateUsername(username string) bool {
	return len(username) >= 3
}

func ValidateUrl(url string) bool {
	pattern := `^http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?$`
	return regexp.MustCompile(pattern).MatchString(url)
}
