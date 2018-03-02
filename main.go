package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"unicode"
	"unicode/utf8"
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
	"Pe": "Punctuation, close",
	"Pi": "Punctuation, initial quote",
	"Pf": "Punctuation, final quote",
	"Po": "Punctuation, other",
	"Ps": "Punctuation, open",
	"S":  "Symbol",
	"Sm": "Symbol, math",
	"Z":  "Separator",
	"Zs": "Separator, space",
}

////////////////////////////////////////////////////////////////////////////////
func main() {
	scripts := flag.Bool("scripts", false, "print contained script names")
	cats := flag.Bool("cats", false, "print unicode categories")
	nochars := flag.Bool("no-chars", false, "do not print character information")
	u8 := flag.Bool("utf8", false, "print utf8 codes")
	flag.Parse()

	var reader *bufio.Reader
	args := flag.Args()
	if len(args) > 0 {
		is, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer func() { _ = is.Close() }()
		reader = bufio.NewReader(is)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	chars := make(map[rune]int)
	for r, _, err := reader.ReadRune(); err == nil; r, _, err = reader.ReadRune() {
		chars[r]++
	}
	if !*nochars {
		printChars(chars, *u8)
	}
	if *cats {
		printCategories(chars)
	}
	if *scripts {
		printScripts(chars)
	}
}

////////////////////////////////////////////////////////////////////////////////
func printChars(chars map[rune]int, u8 bool) {
	var buf [utf8.UTFMax]byte
	keys := sortRuneKeys(chars)
	var p [utf8.UTFMax]byte
	for _, key := range keys {
		char := key
		if !unicode.IsPrint(key) || unicode.IsMark(key) {
			char = ' '
		}
		var cat string
		for name, rng := range unicode.Categories {
			if unicode.In(key, rng) && len(name) > 1 {
				cat = name
			}
		}
		utf8.EncodeRune(p[:], key)
		if u8 {
			fmt.Printf("%c (%-2s %U 0x%x", char, cat, key, p)
			n := utf8.EncodeRune(buf[:], key)
			for i := 0; i < n; i++ {
				fmt.Printf(" 0x%x", buf[i])
			}
			fmt.Printf(") %d\n", chars[key])
		} else {
			fmt.Printf("%c (%-2s %U 0x%x) %d\n", char, cat, key, p, chars[key])
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
func printCategories(chars map[rune]int) {
	categories := make(map[string]int)
	for char, count := range chars {
		for name, rng := range unicode.Categories {
			if unicode.In(char, rng) {
				categories[name] += count
			}
		}
	}
	keys := sortStrKeys(categories)
	fmt.Println("// Categories")
	for _, key := range keys {
		fmt.Printf("%-2s %-21s %d\n", key, longnames[key], categories[key])
	}
}

////////////////////////////////////////////////////////////////////////////////
func printScripts(chars map[rune]int) {
	scripts := make(map[string]int)
	for char, count := range chars {
		for name, rng := range unicode.Scripts {
			if unicode.In(char, rng) {
				scripts[name] += count
			}
		}
	}
	keys := sortStrKeys(scripts)
	fmt.Println("// Scripts")
	for _, key := range keys {
		fmt.Printf("%-24s %d\n", key, scripts[key])
	}
}

////////////////////////////////////////////////////////////////////////////////
func sortStrKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

////////////////////////////////////////////////////////////////////////////////
func sortRuneKeys(m map[rune]int) []rune {
	keys := make([]rune, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
