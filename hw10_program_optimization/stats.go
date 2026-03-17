package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	dec := json.NewDecoder(r)
	result := make(DomainStat)
	suffix := "." + domain

	for {
		var u User
		err := dec.Decode(&u)
		if errors.Is(err, io.EOF) { // исправлено
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}

		if u.Email == "" {
			continue
		}

		email := strings.ToLower(u.Email)
		if strings.HasSuffix(email, suffix) {
			emailSlice := strings.Split(email, "@")
			if len(emailSlice) == 2 {
				result[emailSlice[1]]++
			}
		}
	}

	return result, nil
}
