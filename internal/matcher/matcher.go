package matcher

import "regexp"

type Matcher struct {
	re     *regexp.Regexp
	invert bool
}

func NewMatcher(pattern string, ignoreCase, invert bool) (*Matcher, error) {
	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Matcher{
		re:     re,
		invert: invert,
	}, nil
}

func (m *Matcher) Match(line []byte) bool {
	ok := m.re.Match(line)
	if m.invert {
		return !ok
	}
	return ok
}

func (m *Matcher) Regexp() *regexp.Regexp {
	return m.re
}
