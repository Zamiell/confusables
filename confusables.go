package confusables

// This file has a bunch of commented code because I was converting the "confusables" Python library
// and then later realized that I wanted to keep things simple and only do some basic normalization.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	ConfusablesFileName = "confusables.txt"
)

var (
	confusableMap map[rune]string
)

// Parse the "confusables.txt" file provided by The Unicode Consortium and turn it into a map.
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

		/*
			if utf8.RuneCountInString(char1) == 1 {
				char1Inverted := invertCase(char1)
				if char1Inverted != char1 {
					unicodeConfusableMap[char1] = append(unicodeConfusableMap[char1], char1Inverted)
					unicodeConfusableMap[char1Inverted] = append(unicodeConfusableMap[char1], char1)
				}
			}

			if utf8.RuneCountInString(char2) == 1 {
				char2Inverted := invertCase(char2)
				if char2Inverted != char2 {
					unicodeConfusableMap[char2] = append(unicodeConfusableMap[char1], char2Inverted)
					unicodeConfusableMap[char2Inverted] = append(unicodeConfusableMap[char1], char2)
				}
			}
		*/
	}

	/*
		for _, letter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" {
			char := string(letter)
			accentedCars := getAccentedCharacters(char)
			for _, accentedChar := range accentedCars {
				unicodeConfusableMap[char] = append(unicodeConfusableMap[char], accentedChar)
				unicodeConfusableMap[accentedChar] = append(unicodeConfusableMap[accentedChar], char)
			}
		}

		confusableMap := make(map[string][]string)
		for char := range unicodeConfusableMap {
			charGroup := getConfusableChars(char, unicodeConfusableMap, 0)
			sort.Strings(charGroup)
			confusableMap[char] = charGroup
		}

		var jsonString []byte
		if v, err := json.Marshal(confusableMap); err != nil {
			fmt.Println("Failed to marshal the JSON:", err)
			os.Exit(1)
		} else {
			jsonString = v
		}

		confusableMappingPath := path.Join("assets", ConfusableMappingFileName)
		if err := ioutil.WriteFile(confusableMappingPath, jsonString, 0644); err != nil {
			fmt.Println("Failed to write to \""+confusableMappingPath+"\":", err)
			os.Exit(1)
		}
	*/
}

/*
func invertCase(s string) string {
	if isUpper(s) {
		return strings.ToLower(s)
	}
	return strings.ToUpper(s)
}

// From: https://stackoverflow.com/questions/59293525
func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// getAccentedCharacters goes through every Unicode character and uses NFD (Normalization Form
// Canonical Decomposition) see if it is an accented version of the base character.
func getAccentedCharacters(char string) []string {
	accentedCharacters := make([]string, 0)
	for i := 0; i <= unicode.MaxRune; i++ {
		u := string(i)
		if u != char && asciify(u) == char {
			accentedCharacters = append(accentedCharacters, u)
		}
	}
	return accentedCharacters
}

// asciify turns a Unicode character in an ASCII one.
// It will return an empty string if there is no ASCII representation for the character.
// For example, "ȧ" will be converted to "a".
// For example, "Ƞ" will be converted to "".
func asciify(char string) string {
	// Normalize the character using NFD (Normalization Form Canonical Decomposition). NFD must be
	// used over NFC because we need to separate the base character from any accents.
	// https://en.wikipedia.org/wiki/Unicode_equivalence
	char = norm.NFD.String(char)

	// Strip out any runes (code points) that are non-printable ASCII characters.
	// https://www.ascii-code.com/
	// We explicitly avoid using the more-intuitive "for _, codePoint := range char {"
	// because that performs unnecessary rune conversions.
	// https://stackoverflow.com/questions/53069040/checking-a-string-contains-only-ascii-characters
	asciiChars := make([]byte, 0)
	for i := 0; i < len(char); i++ {
		codePoint := char[i]
		if codePoint >= 32 && codePoint <= 126 {
			asciiChars = append(asciiChars, codePoint)
		}
	}
	return string(asciiChars)
}

func getConfusableChars(
	character string,
	unicodeConfusableMap map[string][]string,
	depth int,
) []string {
	mappedChars := unicodeConfusableMap[character]

	group := []string{character}
	if depth <= MaxSimilarityDepth {
		for _, mappedChar := range mappedChars {
			mappedChars2 := getConfusableChars(mappedChar, unicodeConfusableMap, depth+1)
			for _, mappedChar2 := range mappedChars2 {
				if !stringInSlice(mappedChar2, group) {
					group = append(group, mappedChar2)
				}
			}
		}
	}
	return group
}

func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}
*/

// Normalize returns a copy of a string with common Unicode homoglyphs replaced with their
// more-standard versions.
func Normalize(s string) string {
	normalizedString := s
	for _, r := range s {
		if replacement, ok := confusableMap[r]; ok {
			normalizedString = strings.Replace(normalizedString, string(r), replacement, -1)
		}
	}

	return s
}
