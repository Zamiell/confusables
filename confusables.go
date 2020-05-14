package confusables

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
)

const (
	ConfusablesFileName = "confusables.txt"
)

var (
	confusableMap map[rune]string
)

// When this package is initialized, parse the "confusables.txt" file provided by The Unicode
// Consortium and turn it into a map.
func init() {
	confusablesPath := path.Join("assets", ConfusablesFileName)
	if _, err := os.Stat(confusablesPath); os.IsNotExist(err) {
		fmt.Println("\"" + confusablesPath + "\" does not exist.")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Failed to check to see if \""+confusablesPath+"\" exists:", err)
		os.Exit(1)
	}

	var confusableLines []string
	if fileContents, err := ioutil.ReadFile(confusablesPath); err != nil {
		fmt.Println("Failed to read the \""+confusablesPath+"\" file:", err)
		os.Exit(1)
	} else {
		confusablesString := string(fileContents)
		confusableLines = strings.Split(confusablesString, "\n")
	}

	confusableMap = make(map[rune]string)

	for i, line := range confusableLines {
		// Ignore the first line, which should just be a comment of "# confusables.txt". This line
		// will also start with an invisible byte order mark to signify that this text file contains
		// Unicode.
		// https://en.wikipedia.org/wiki/Byte_order_mark
		lineNumber := i + 1
		if lineNumber == 1 {
			continue
		}

		// Ignore empty lines.
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Ignore comments.
		if strings.HasPrefix(line, "#") {
			continue
		}

		// The format used in the confusables file is:
		// 1D5A4 ;	0045 ;	MA	# ( 𝖤 → E ) MATHEMATICAL SANS-SERIF CAPITAL E → LATIN CAPITAL LETTER E	#
		mapping := strings.Split(line, ";")

		// Get the first character (e.g. the confusing character). This is represented as a hex
		// string (e.g. "2FA1D"). It is always one rune, so we don't have to worry about splitting
		// on spaces.
		char1String := "0x" + strings.TrimSpace(mapping[0])
		var char1Int int64
		if v, err := strconv.ParseInt(char1String, 0, 64); err != nil {
			fmt.Println("Failed to convert \""+char1String+"\" to an integer on line "+
				strconv.Itoa(lineNumber)+":", err)
			os.Exit(1)
		} else {
			char1Int = v
		}
		char1 := rune(char1Int)

		// Get the second character (e.g. the character that the confusing character looks like).
		// This is represented as one or more hex strings (e.g. "2A600", "0028 0072 006E 0029").
		char2String := strings.TrimSpace(mapping[1])
		char2StringArray := strings.Split(char2String, " ")
		char2Array := make([]rune, 0)
		for _, hexStr := range char2StringArray {
			hexStr = "0x" + hexStr
			var charInt int64
			if v, err := strconv.ParseInt(hexStr, 0, 64); err != nil {
				fmt.Println("Failed to convert \""+hexStr+"\" to an integer on line "+
					strconv.Itoa(lineNumber)+":", err)
				os.Exit(1)
			} else {
				charInt = v
			}
			char2Array = append(char2Array, rune(charInt))
		}
		char2 := string(char2Array)

		// See: https://staticcheck.io/docs/checks#S1036
		if _, ok := confusableMap[char1]; ok {
			fmt.Println("Failed to parse \"" + ConfusablesFileName + "\". There is a duplicate " +
				"rune on line " + strconv.Itoa(lineNumber) + ":")
			fmt.Println(line)
			os.Exit(1)
		}
		confusableMap[char1] = char2
	}
}

// Normalize returns a copy of a string with common Unicode homoglyphs replaced with their
// more-standard versions.
func Normalize(s string) string {
	// First, normalize the string with NFD.
	// https://blog.golang.org/normalization
	// We need to use NFD instead of NFC because we need to separate accents from the base
	// character. Otherwise, we wouldn't be able to find the match in "confusables.txt". For an
	// example, of this, see "TestNormalizeAccentAndNFD()".
	normalizedString := norm.NFD.String(s)

	// Second, replace homoglyphs (as reported by "confusables.txt")
	for _, r := range s {
		if replacement, ok := confusableMap[r]; ok {
			normalizedString = strings.ReplaceAll(normalizedString, string(r), replacement)
		}
	}

	return normalizedString
}
