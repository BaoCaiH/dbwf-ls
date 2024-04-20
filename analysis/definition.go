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
	re, err := regexp.Compile(fmt.Sprintf("^[\\s-]*(job_cluster_key|task_key):\\s*\"?(%s)\"?\\s*#*.*$", item_name))
	if err != nil {
		logger.Println(err)
		return item
	}

	lines := strings.Split(document, "\n")
	documentLength := len(lines)

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 3 {
			matchIndex := re.FindStringSubmatchIndex(line)
			var isNewCluster bool
			if i == documentLength-1 {
				isNewCluster = false
			} else {
				isNewCluster, err = regexp.MatchString("^.*new_cluster:\\s*#*.*$", lines[i+1])
				if err != nil {
					logger.Println(err)
					return item
				}
			}
			var isTaskDefinition bool
			if i == documentLength-1 {
				isTaskDefinition = false
			} else {
				isTaskDefinition, err = regexp.MatchString("^.*description:\\s*#*.*$", lines[i+1])
				if err != nil {
					logger.Println(err)
					return item
				}
			}
			if isNewCluster || isTaskDefinition {
				item.defined = lsp.LineRange(i, matchIndex[4], matchIndex[5])
			}
		}
	}

	return item
}
