// Command chbs generates xkcd-style passphrases with the same options as
// correcthorsebatterystaple.net, so agents, scripts, and humans can drive it
// from a shell instead of importing a library.
package main

import (
	"crypto/rand"
	"embed"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Word data from correcthorsebatterystaple.net (MIT, (c) John Van Der Loo).
//
//go:embed data/wordlist.txt data/jargon.txt data/science-terms.txt
var dictFS embed.FS

var dictFiles = map[string]string{
	"jargon":  "data/jargon.txt",
	"science": "data/science-terms.txt",
}

func main() {
	minWords := flag.Int("n", 4, "minimum number of words")
	minLength := flag.Int("l", 20, "minimum phrase length, separators included")
	sep := flag.String("s", "-", "separator; with more than one character one is drawn per gap")
	upper := flag.Bool("upper", true, "capitalize the first letter of every word")
	digits := flag.Bool("digits", true, "append one random digit (0-9)")
	dicts := flag.String("dict", "", "extra wordlists on top of the primary one, comma separated: jargon,science")
	flag.Parse()

	if *minWords < 1 {
		fmt.Fprintln(os.Stderr, "chbs: -n must be >= 1")
		os.Exit(2)
	}
	if *minLength < 1 {
		fmt.Fprintln(os.Stderr, "chbs: -l must be >= 1")
		os.Exit(2)
	}

	words, err := loadWords(*dicts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chbs:", err)
		os.Exit(2)
	}

	fmt.Println(generate(words, *minWords, *minLength, *sep, *upper, *digits))
}

// loadWords reads the always-on primary list plus any extra dictionaries.
func loadWords(extra string) ([]string, error) {
	files := []string{"data/wordlist.txt"}
	for d := range strings.SplitSeq(extra, ",") {
		switch d = strings.TrimSpace(strings.ToLower(d)); d {
		case "":
		case "jargon", "science":
			files = append(files, dictFiles[d])
		default:
			return nil, fmt.Errorf("unknown dictionary %q (have: jargon, science)", d)
		}
	}

	var words []string
	for _, f := range files {
		b, err := dictFS.ReadFile(f)
		if err != nil {
			return nil, err
		}
		// Comma separated; trailing commas and newlines fall out as empty fields.
		for w := range strings.FieldsFuncSeq(string(b), func(r rune) bool { return r == ',' || r == '\n' || r == '\r' }) {
			words = append(words, w)
		}
	}
	return words, nil
}

// generate mirrors the site: draw minWords first, then keep adding words until
// the phrase joined with the separator's first character reaches minLength.
// The final phrase draws a separator per gap, and the digit is appended
// without counting toward minLength.
func generate(words []string, minWords, minLength int, sep string, upper, digits bool) string {
	pick := func() string {
		w := words[randInt(len(words))]
		if upper {
			w = strings.ToUpper(w[:1]) + w[1:]
		}
		return w
	}

	seps := []rune(sep)
	testSep := ""
	if len(seps) > 0 {
		testSep = string(seps[0])
	}

	picked := make([]string, 0, minWords+4)
	testLen := 0
	add := func() {
		w := pick()
		if len(picked) > 0 {
			testLen += len(testSep)
		}
		testLen += len(w)
		picked = append(picked, w)
	}
	for range minWords {
		add()
	}
	for testLen < minLength {
		add()
	}

	var b strings.Builder
	for i, w := range picked {
		b.WriteString(w)
		if i < len(picked)-1 && len(seps) > 0 {
			b.WriteRune(seps[randInt(len(seps))])
		}
	}
	if digits {
		if len(seps) > 0 {
			b.WriteRune(seps[randInt(len(seps))])
		}
		b.WriteString(strconv.Itoa(randInt(10)))
	}
	return b.String()
}

// randInt returns a uniform int in [0, n). A failed random source aborts the
// process: emitting a phrase from a broken source would be worse than no output.
func randInt(n int) int {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "chbs: random source:", err)
		os.Exit(1)
	}
	return int(v.Int64())
}
