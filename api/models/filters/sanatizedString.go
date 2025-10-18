package filters

import (
	"api/sanatizer"
	"unicode/utf8"
)

type SanatizerString string

func (s SanatizerString) Sanatize() SanatizerString {
	return SanatizerString(sanatizer.SanitizeForHTML(string(s)))
}

func (s SanatizerString) SanatizeSQL() SanatizerString {
	return SanatizerString(sanatizer.SanitizeForSQL(string(s)))
}

func (s SanatizerString) SanatizedAll() SanatizerString {
	ss := s.Sanatize()
	ss = s.SanatizeSQL()
	return ss
}

func (s SanatizerString) WithMaxLength(maxLen int) SanatizerString {
	ss := string(s)
	if maxLen > 0 && utf8.RuneCountInString(ss) > maxLen {
		runes := []rune(ss)
		ss = string(runes[:maxLen])
	}

	return SanatizerString(ss)
}

func (s SanatizerString) String() string {
	return string(s)
}
