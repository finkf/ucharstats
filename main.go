package main

import (
	"bufio"
	"fmt"
	"flag"
	"os"
	"sort"
	"unicode"
)

var longnames = map[string]string{
	"C":  "Other",
	"Cc": "Other, control",
	"L":  "Letter",
	"Ll": "Letter, lowercase",
	"Lu": "Letter, uppercase",
	"M":  "Mark",
	"Mn": "Mark, nonspacing",
	"N":  "Number",
	"Nd": "Number, decimal digit",
	"P":  "Punctuation",
	"Pc": "Punctuation, connector",
	"Pd": "Punctuation, dash",
	"Pi": "Punctuation, initial quote",
	"Pf": "Punctuation, final quote",
	"Po": "Punctuation, other",
	"Ps": "Punctuation, open",
	"S":  "Symbol",
	"Sm": "Symbol, math",
	"Z":  "Separator",
	"Zs": "Separator, space",
}

func main() {
	long   := flag.Bool("long", false, "print long unicode classes")
	script := flag.Bool("script", false, "print long contained script names")
	flag.Parse()
	
	reader := bufio.NewReader(os.Stdin)
	cats := make(map[string]int)
	scripts := make(map[string]int)

	for r, _, err := reader.ReadRune(); err == nil; r, _, err = reader.ReadRune() {
		countCategory(cats, r)
		if *script {
			countScript(scripts, r)
		}
	}
	keys := sortKeys(cats)
	for _, k := range keys {
		if *long {
			fmt.Printf("%-2s %5d %s\n", k, cats[k], longnames[k])
		} else {
			fmt.Printf("%-2s %5d\n", k, cats[k])
		}
	}
	if *script {
		keys = sortKeys(scripts)
		for _, k := range keys {
			fmt.Printf("%-21s %d\n", k, scripts[k])
		}
	}
}


func countCategory(counts map[string]int, r rune) {
	for name, rng := range unicode.Categories {
		if unicode.In(r, rng) {
			counts[name]++
		}
	}
}

func countScript(scripts map[string]int, r rune) {
	for name, rng := range unicode.Scripts {
		if unicode.In(r, rng) {
			scripts[name]++
		}
	}
}

func sortKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
