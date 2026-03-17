package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat, 1000) // предвыделение памяти
	suffix := "." + domain
	emailPrefix := []byte(`"Email":"`)

	for scanner.Scan() {
		line := scanner.Bytes()
		// Ищем начало поля Email
		idx := bytes.Index(line, emailPrefix)
		if idx == -1 {
			continue
		}
		start := idx + len(emailPrefix)
		if start >= len(line) {
			continue
		}
		// Ищем закрывающую кавычку
		end := bytes.IndexByte(line[start:], '"')
		if end == -1 {
			continue
		}
		emailBytes := line[start : start+end]
		email := strings.ToLower(string(emailBytes))

		// Проверяем окончание на нужный домен
		if !strings.HasSuffix(email, suffix) {
			continue
		}
		atIndex := strings.LastIndexByte(email, '@')
		if atIndex == -1 {
			continue
		}
		domainName := email[atIndex+1:]
		result[domainName]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
