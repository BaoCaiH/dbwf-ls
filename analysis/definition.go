package analysis

import (
	"dbwf-ls/lsp"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type definition struct {
	defined, lastReferred lsp.Range
}

func findDefinition(document, item_name string, logger *log.Logger) definition {
	item := definition{}
	re, err := regexp.Compile(fmt.Sprintf("^.*[job_cluster_key|task_key]:\\s*\"?(%s)\"?\\s*/?/?.*$", item_name))
	if err != nil {
		logger.Println(err)
		return item
	}

	lines := strings.Split(document, "\n")
	documentLength := len(lines)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 2 {
			matchIndex := re.FindStringSubmatchIndex(line)
			var isNewCluster bool
			if i == documentLength-1 {
				isNewCluster = false
			} else {
				isNewCluster, err = regexp.MatchString("^\\s*new_cluster:\\s*/?/?.*$", lines[i+1])
				if err != nil {
					logger.Println(err)
					return item
				}
			}
			var isTaskRefer bool
			if i == 0 {
				isTaskRefer = false
			} else {
				isTaskRefer, err = regexp.MatchString("^\\s*depends_on:\\s*/?/?.*$", lines[i-1])
				if err != nil {
					logger.Println(err)
					return item
				}
			}
			if isNewCluster || !isTaskRefer {
				item.defined = lsp.LineRange(i, matchIndex[2], matchIndex[3])
			}
		}
	}

	return item
}
