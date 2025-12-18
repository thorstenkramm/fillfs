package filenames

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRandomDirectoryName(t *testing.T) {
	name := RandomDirectoryName()
	if len(name) < 7 {
		t.Fatalf("name too short: %q", name)
	}

	prefix := name[:3]
	if _, err := strconv.Atoi(prefix); err != nil {
		t.Fatalf("prefix is not numeric: %q", prefix)
	}
	if !strings.ContainsRune(separators, rune(name[3])) {
		t.Fatalf("separator missing after prefix: %q", name)
	}

	sepCount := 0
	for _, r := range name[3:] {
		if strings.ContainsRune(separators, r) {
			sepCount++
		}
	}
	if sepCount < 2 {
		t.Fatalf("expected at least two separators in %q, got %d", name, sepCount)
	}
}

func TestRandomFileNamesIncludeOptionalDateEveryThird(t *testing.T) {
	generators := []struct {
		name string
		fn   func() string
	}{
		{"document", RandomDocumentFileName},
		{"spreadsheet", RandomSpreadsheetFileName},
		{"image", RandomImageFileName},
		{"sound", RandomSoundFileName},
		{"powerpoint", RandomPowerpointFileName},
	}

	for _, g := range generators {
		t.Run(g.name, func(t *testing.T) {
			nameCount = 0 // reset counter for deterministic third-call behavior

			withoutDate := regexp.MustCompile(`^.+[._\- \+=]v([1-9]|1[0-9]|2[0-5])$`)
			withDate := regexp.MustCompile(`^.+[._\- \+=]v([1-9]|1[0-9]|2[0-5])[._\- \+=][\(\{\[]?(\d{4}-\d{2}-\d{2})[\)\}\]]?$`)

			first := g.fn()
			second := g.fn()
			third := g.fn()

			if !withoutDate.MatchString(first) {
				t.Fatalf("first name missing version: %q", first)
			}
			if !withoutDate.MatchString(second) {
				t.Fatalf("second name missing version: %q", second)
			}

			matches := withDate.FindStringSubmatch(third)
			if len(matches) != 3 {
				t.Fatalf("third name missing date: %q", third)
			}

			version, err := strconv.Atoi(matches[1])
			if err != nil || version < 1 || version > 25 {
				t.Fatalf("invalid version in %q", third)
			}

			dateStr := matches[2]
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				t.Fatalf("invalid date %q: %v", dateStr, err)
			}
			if date.Before(startDate) || date.After(time.Now().UTC().Add(24*time.Hour)) { // allow small clock skew
				t.Fatalf("date %q out of expected range", dateStr)
			}
		})
	}
}
