package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/mrusme/symmetry454"
)

const usageText = `sym454 shows a Symmetry454 calendar and converts dates.

Usage:
  sym454 [flags] [[month] year]
  sym454 -c [flags] < dates
  sym454 --from-sym [flags] < dates

With no arguments it prints the current Symmetry454 month and highlights today.
A single argument is a year and prints all twelve months. Two arguments are a
month and a year, in that order.

Converter mode reads one date per line from stdin and accepts Unix epoch
seconds from "date +%s", RFC 1123 with an offset from "date -R", RFC 3339 from
"date -Is", the default output of "date", and plain YYYY-MM-DD. Epoch seconds
are read in the local zone unless --utc is given.

Calendar flags:
  -y, --year        print all twelve months of the year
  -3                print the previous, current and next month
  --across n        months per row in the year view (default 3)
  --gregorian       print the Gregorian span covered by the output
  --color mode      highlight today: auto, always or never (default auto)

Converter flags:
  -c, --convert     read dates from stdin and print Symmetry454 dates
  --from-sym        read Symmetry454 dates from stdin and print Gregorian dates
  --utc             read epoch seconds as UTC rather than local time
  --long            print a labeled block for each date
  --json            print one JSON object per date
  -f, --format s    Go template for each output line

Examples:
  sym454
  sym454 -y 2026
  sym454 8 2026
  date +%s | sym454 -c
  date -R | sym454 -c --long
  echo 2026-08-35 | sym454 --from-sym
`

type options struct {
	convert       bool
	fromSym       bool
	year          bool
	three         bool
	utc           bool
	long          bool
	asJSON        bool
	showGregorian bool
	format        string
	color         string
	perRow        int
	targetYear    int
	targetMonth   int
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

func colorEnabled(mode string, file *os.File) bool {
	switch mode {
	case "always":
		return true
	case "never":
		return false
	}

	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	return isTerminal(file)
}

func isNegativeNumber(token string) bool {
	if len(token) < 2 || token[0] != '-' {
		return false
	}

	for _, char := range token[1:] {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func parseArgs(
	argv []string,
	today symmetry454.Sym,
	errOut io.Writer,
) (*options, int) {
	opts := options{
		targetYear:  today.Year,
		targetMonth: today.Month,
	}

	flags := flag.NewFlagSet("sym454", flag.ContinueOnError)
	flags.SetOutput(errOut)
	flags.Usage = func() { io.WriteString(errOut, usageText) }

	flags.BoolVar(&opts.convert, "c", false, "read dates from stdin")
	flags.BoolVar(&opts.convert, "convert", false, "read dates from stdin")
	flags.BoolVar(&opts.fromSym, "from-sym", false,
		"read Symmetry454 dates from stdin")
	flags.BoolVar(&opts.year, "y", false, "print all twelve months")
	flags.BoolVar(&opts.year, "year", false, "print all twelve months")
	flags.BoolVar(&opts.three, "3", false,
		"print the previous, current and next month")
	flags.BoolVar(&opts.utc, "utc", false, "read epoch seconds as UTC")
	flags.BoolVar(&opts.long, "long", false, "print a labeled block per date")
	flags.BoolVar(&opts.asJSON, "json", false, "print one JSON object per date")
	flags.BoolVar(&opts.showGregorian, "gregorian", false,
		"print the Gregorian span covered by the output")
	flags.StringVar(&opts.format, "f", "", "Go template for each output line")
	flags.StringVar(&opts.format, "format", "",
		"Go template for each output line")
	flags.StringVar(&opts.color, "color", "auto",
		"highlight today: auto, always or never")
	flags.IntVar(&opts.perRow, "across", 3, "months per row in the year view")

	args := make([]string, 0, 2)
	rest := argv
	for {
		for len(rest) > 0 && isNegativeNumber(rest[0]) &&
			flags.Lookup(rest[0][1:]) == nil {
			args = append(args, rest[0])
			rest = rest[1:]
		}

		if err := flags.Parse(rest); err != nil {
			return nil, 2
		}

		rest = flags.Args()
		if len(rest) == 0 {
			break
		}

		args = append(args, rest[0])
		rest = rest[1:]
	}

	switch opts.color {
	case "auto", "always", "never":
	default:
		fmt.Fprintf(errOut,
			"sym454: invalid --color %q, want auto, always or never\n", opts.color)
		return nil, 2
	}

	if opts.perRow < 1 || opts.perRow > 12 {
		fmt.Fprintf(errOut,
			"sym454: invalid --across %d, want 1 to 12\n", opts.perRow)
		return nil, 2
	}

	switch len(args) {
	case 0:
	case 1:
		year, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(errOut, "sym454: invalid year %q\n", args[0])
			return nil, 2
		}
		opts.targetYear = year
		if !opts.three {
			opts.year = true
		}
	case 2:
		month, err := strconv.Atoi(args[0])
		if err != nil || month < 1 || month > 12 {
			fmt.Fprintf(errOut, "sym454: invalid month %q, want 1 to 12\n", args[0])
			return nil, 2
		}
		year, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(errOut, "sym454: invalid year %q\n", args[1])
			return nil, 2
		}
		opts.targetMonth = month
		opts.targetYear = year
	default:
		flags.Usage()
		return nil, 2
	}

	return &opts, 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	today := symmetry454.FromTime(time.Now())

	opts, status := parseArgs(argv, today, os.Stderr)
	if opts == nil {
		return status
	}

	if opts.convert || opts.fromSym {
		location := time.Local
		if opts.utc {
			location = time.UTC
		}

		return runConvert(os.Stdin, os.Stdout, os.Stderr, convertOptions{
			location: location,
			format:   opts.format,
			long:     opts.long,
			asJSON:   opts.asJSON,
			fromSym:  opts.fromSym,
		})
	}

	style := styler{enabled: colorEnabled(opts.color, os.Stdout)}

	switch {
	case opts.year:
		fmt.Print(renderYear(opts.targetYear, today, style, opts.perRow,
			opts.showGregorian))
	case opts.three:
		fmt.Print(renderThreeMonths(opts.targetYear, opts.targetMonth, today,
			style, opts.showGregorian))
	default:
		fmt.Print(renderSingleMonth(opts.targetYear, opts.targetMonth, today,
			style, opts.showGregorian))
	}

	return 0
}
