package util

import (
	"regexp"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Name struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type MatchResult struct {
	Matched bool

	FirstNameMatch  bool
	MiddleNameMatch bool
	LastNameMatch   bool
}

func VerifyName(user Name, verified Name) MatchResult {
	userFirst := normalize(user.FirstName)
	userMiddle := normalize(user.MiddleName)
	userLast := normalize(user.LastName)

	verifiedFirst := normalize(verified.FirstName)
	verifiedMiddle := normalize(verified.MiddleName)
	verifiedLast := normalize(verified.LastName)

	firstMatch := fuzzy.RankMatchNormalizedFold(
		userFirst,
		verifiedFirst,
	) >= 0

	lastMatch := fuzzy.RankMatchNormalizedFold(
		userLast,
		verifiedLast,
	) >= 0

	middleMatch := true

	if userMiddle != "" && verifiedMiddle != "" {
		middleMatch = fuzzy.RankMatchNormalizedFold(
			userMiddle,
			verifiedMiddle,
		) >= 0
	}

	return MatchResult{
		Matched:         firstMatch && lastMatch,
		FirstNameMatch:  firstMatch,
		LastNameMatch:   lastMatch,
		MiddleNameMatch: middleMatch,
	}
}

func normalize(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))

	reg := regexp.MustCompile(`[^a-z0-9\s]`)
	name = reg.ReplaceAllString(name, "")

	return strings.Join(strings.Fields(name), " ")
}
