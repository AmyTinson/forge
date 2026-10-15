package main

import (
	"errors"
	"regexp"
)

var likelyJWT = regexp.MustCompile(`\b[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\b`)

var likelyCredentialPatterns = []*regexp.Regexp{
	regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`),
	regexp.MustCompile(`(?i)\b(?:aws[_-]?)?(?:api[_-]?key|access[_-]?token|secret[_-]?(?:access[_-]?)?key|client[_-]?secret|password|passwd|token)\s*[:=]\s*["']?[A-Za-z0-9_./+=:-]{12,}`),
	regexp.MustCompile(`\b(?:AKIA|ASIA)[A-Z0-9]{16}\b`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}\b`),
	regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}\b`),
	regexp.MustCompile(`\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b`),
	regexp.MustCompile(`\bsk-ant-[A-Za-z0-9_-]{20,}\b`),
	regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/=-]{20,}\b`),
}

func rejectLikelyCredentials(value string) error {
	for _, pattern := range likelyCredentialPatterns {
		if pattern.MatchString(value) {
			return errors.New("content appears to contain a credential")
		}
	}
	if likelyJWT.MatchString(value) {
		return errors.New("content appears to contain a credential")
	}
	return nil
}
