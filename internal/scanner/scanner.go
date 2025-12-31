package scanner

import "encoding/json"

type Scanner struct {
	encoder *json.Encoder
}
type Finding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Pattern  string `json:"pattern"`
	Severity string `json:"severity"`
	Match    string `json:"match"`
}

func NewScanner(encoder *json.Encoder) *Scanner {
	return &Scanner{encoder: encoder}
}
