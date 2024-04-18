package analysis

import (
	"dbwf-ls/lsp"
	"log"
	"regexp"
)

func wordAtCursor(line string, position lsp.Position, logger *log.Logger) (string, error) {
	re, err := regexp.Compile("\\W")
	if err != nil {
		logger.Printf("Regexp Compile %s", err)
		return "", err
	}

	// Because the flocking cursor is 1 step ahead of the line while typing
	// So this can fail, quietly, damn.
	char := position.Character
	if char == len(line) {
		char--
	}

	if loc := re.FindStringIndex(line[char : char+1]); loc != nil {
		return "", nil
	}

	start, end := 0, 0
	if locs := re.FindAllStringIndex(line, -1); locs != nil {
		for _, loc := range locs {
			if loc[0] > position.Character {
				end = loc[0]
			}
			if loc[1] <= position.Character {
				start = loc[1]
			}
			if end != 0 {
				break
			}
		}
	}
	if end == 0 {
		end = len(line)
	}

	return line[start:end], nil
}

func leadingSpaces(line string, logger *log.Logger) (string, error) {
	re, err := regexp.Compile("^\\s*")
	if err != nil {
		logger.Panicf("Regexp Compile %s", err)
		return "", err
	}

	return re.FindString(line), nil
}
