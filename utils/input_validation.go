package utils

import (
	sErr "Engee-Server/stockErrors"
	"net/url"
)

func ValidateURL(endpoint string) error {
	if endpoint == "" {
		return &sErr.EmptyValueError{
			Field: "URL",
		}
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		return &sErr.EmptyValueError{
			Field: "Scheme",
		}
	}

	if u.Host == "" {
		return &sErr.EmptyValueError{
			Field: "Host",
		}
	}

	return nil
}
