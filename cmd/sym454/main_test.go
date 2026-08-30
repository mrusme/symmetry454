package main

import (
	"strings"
	"testing"

	"github.com/mrusme/symmetry454"
)

var referenceToday = symmetry454.Sym{Year: 2026, Month: 8, Day: 35}

func TestParseArgsDefaults(t *testing.T) {
	var errOut strings.Builder

	opts, status := parseArgs(nil, referenceToday, &errOut)
	if opts == nil {
		t.Fatalf(`parseArgs returned status %v, stderr %q`, status, errOut.String())
	}

	if opts.targetYear != 2026 || opts.targetMonth != 8 {
		t.Fatalf(`target = %v-%v, want 2026-8`, opts.targetYear, opts.targetMonth)
	}
	if opts.year || opts.three || opts.convert || opts.fromSym {
		t.Fatalf(`a mode was set without being asked for: %+v`, *opts)
	}
	if opts.color != "auto" || opts.perRow != 3 {
		t.Fatalf(`color = %q, perRow = %v, want "auto" and 3`,
			opts.color, opts.perRow)
	}
}

func TestParseArgsPositional(t *testing.T) {
	tests := []struct {
		argv      []string
		wantYear  int
		wantMonth int
		yearMode  bool
	}{
		{[]string{}, 2026, 8, false},
		{[]string{"2027"}, 2027, 8, true},
		{[]string{"12", "2026"}, 2026, 12, false},
		{[]string{"1", "-121"}, -121, 1, false},
		{[]string{"-y", "2027"}, 2027, 8, true},
		{[]string{"-3", "2027"}, 2027, 8, false},
	}

	for _, test := range tests {
		var errOut strings.Builder

		opts, status := parseArgs(test.argv, referenceToday, &errOut)
		if opts == nil {
			t.Fatalf(`parseArgs(%v) returned status %v, stderr %q`,
				test.argv, status, errOut.String())
		}

		if opts.targetYear != test.wantYear || opts.targetMonth != test.wantMonth {
			t.Fatalf(`parseArgs(%v) target = %v-%v, want %v-%v`,
				test.argv, opts.targetYear, opts.targetMonth,
				test.wantYear, test.wantMonth)
		}
		if opts.year != test.yearMode {
			t.Fatalf(`parseArgs(%v) year mode = %v, want %v`,
				test.argv, opts.year, test.yearMode)
		}
	}
}

func TestParseArgsFlagsAfterPositional(t *testing.T) {
	tests := [][]string{
		{"12", "2026", "--gregorian"},
		{"--gregorian", "12", "2026"},
		{"12", "--gregorian", "2026"},
	}

	for _, argv := range tests {
		var errOut strings.Builder

		opts, status := parseArgs(argv, referenceToday, &errOut)
		if opts == nil {
			t.Fatalf(`parseArgs(%v) returned status %v, stderr %q`,
				argv, status, errOut.String())
		}

		if !opts.showGregorian {
			t.Fatalf(`parseArgs(%v) did not set --gregorian`, argv)
		}
		if opts.targetYear != 2026 || opts.targetMonth != 12 {
			t.Fatalf(`parseArgs(%v) target = %v-%v, want 2026-12`,
				argv, opts.targetYear, opts.targetMonth)
		}
	}
}

func TestParseArgsFlagValuesAfterPositional(t *testing.T) {
	var errOut strings.Builder

	opts, status := parseArgs([]string{"2026", "--across", "4", "--color", "never"},
		referenceToday, &errOut)
	if opts == nil {
		t.Fatalf(`parseArgs returned status %v, stderr %q`, status, errOut.String())
	}

	if opts.perRow != 4 {
		t.Fatalf(`perRow = %v, want 4`, opts.perRow)
	}
	if opts.color != "never" {
		t.Fatalf(`color = %q, want "never"`, opts.color)
	}
	if opts.targetYear != 2026 || !opts.year {
		t.Fatalf(`target year = %v, year mode = %v, want 2026 and true`,
			opts.targetYear, opts.year)
	}
}

func TestParseArgsLongAndShortForms(t *testing.T) {
	tests := []struct {
		argv  []string
		check func(*options) bool
	}{
		{[]string{"-c"}, func(o *options) bool { return o.convert }},
		{[]string{"--convert"}, func(o *options) bool { return o.convert }},
		{[]string{"-y"}, func(o *options) bool { return o.year }},
		{[]string{"--year"}, func(o *options) bool { return o.year }},
		{[]string{"-3"}, func(o *options) bool { return o.three }},
		{[]string{"--from-sym"}, func(o *options) bool { return o.fromSym }},
		{[]string{"--utc"}, func(o *options) bool { return o.utc }},
		{[]string{"--long"}, func(o *options) bool { return o.long }},
		{[]string{"--json"}, func(o *options) bool { return o.asJSON }},
		{[]string{"-f", "x"}, func(o *options) bool { return o.format == "x" }},
		{[]string{"--format", "x"}, func(o *options) bool { return o.format == "x" }},
	}

	for _, test := range tests {
		var errOut strings.Builder

		opts, status := parseArgs(test.argv, referenceToday, &errOut)
		if opts == nil {
			t.Fatalf(`parseArgs(%v) returned status %v, stderr %q`,
				test.argv, status, errOut.String())
		}
		if !test.check(opts) {
			t.Fatalf(`parseArgs(%v) did not apply the flag: %+v`, test.argv, *opts)
		}
	}
}

func TestParseArgsErrors(t *testing.T) {
	tests := [][]string{
		{"13", "2026"},
		{"0", "2026"},
		{"notayear"},
		{"8", "notayear"},
		{"--color", "purple"},
		{"--across", "0"},
		{"--across", "13"},
		{"1", "2", "3"},
		{"--nosuchflag"},
	}

	for _, argv := range tests {
		var errOut strings.Builder

		opts, status := parseArgs(argv, referenceToday, &errOut)
		if opts != nil {
			t.Fatalf(`parseArgs(%v) succeeded, want an error`, argv)
		}
		if status != 2 {
			t.Fatalf(`parseArgs(%v) status = %v, want 2`, argv, status)
		}
		if errOut.Len() == 0 {
			t.Fatalf(`parseArgs(%v) wrote nothing to stderr`, argv)
		}
	}
}
