package analysis

import (
	"cmp"
	"dbwf-ls/lsp"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
)

type requiredKeys struct {
	name string
	Diag
	found bool
}

func diagnose(document string, logger *log.Logger) []lsp.Diagnostics {
	diagnostics := []lsp.Diagnostics{}
	required_keywords := map[string]requiredKeys{}
	for k, v := range Keywords {
		if v.diag.severity != 0 {
			required_keywords[k] = requiredKeys{
				Diag:  v.diag,
				found: false,
			}
		}
	}
	foundJobClusterChunk := false
	jobClusters := map[string]definition{}
	re, err := regexp.Compile("^.*job_cluster_key:\\s*\"?(\\w*)\"?\\s*/?/?.*$")
	if err != nil {
		logger.Println(err)
		return diagnostics
	}

	lines := strings.Split(document, "\n")
	documentLength := len(lines)

	for i, line := range lines {
		// required match
		for k, v := range required_keywords {
			if v.found {
				continue
			}
			matched, err := regexp.MatchString(fmt.Sprintf("^%s:\\s*/?/?.*$", k), line)
			if err != nil {
				logger.Println(err)
				return diagnostics
			}
			if matched {
				v.found = true
				required_keywords[k] = v
			}
		}
		// conditional match
		if !foundJobClusterChunk {
			matched, err := regexp.MatchString("^job_clusters:\\s*/?/?.*$", line)
			if err != nil {
				logger.Println(err)
				return diagnostics
			}
			if matched {
				foundJobClusterChunk = true
			}
		}
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
					return diagnostics
				}
			}
			current := jobClusters[matches[1]]
			if isNewCluster {
				jobClusters[matches[1]] = definition{
					defined:      lsp.LineRange(i, matchIndex[2], matchIndex[3]),
					lastReferred: current.lastReferred,
				}
			} else {
				jobClusters[matches[1]] = definition{
					defined:      current.defined,
					lastReferred: lsp.LineRange(i, matchIndex[2], matchIndex[3]),
				}
			}
		}
	}

	for _, v := range required_keywords {
		if !v.found {
			diagnostics = append(diagnostics, lsp.Diagnostics{
				Range:    lsp.LineRange(documentLength-1, 0, 0),
				Severity: v.severity,
				Source:   "dbwf-ls",
				Message:  v.help,
			})
		}
	}

	if len(jobClusters) > 0 && !foundJobClusterChunk {
		diagnostics = append(diagnostics, lsp.Diagnostics{
			Range:    lsp.LineRange(documentLength-1, 0, 0),
			Severity: 1,
			Source:   "dbwf-ls",
			Message:  "`job_cluster_key` is declared on task but no `job_clusters` chunk found. Hint: start by typing `cluster`",
		})
	} else {
		for k, v := range jobClusters {
			if v.defined != lsp.LineRange(0, 0, 0) && v.lastReferred == lsp.LineRange(0, 0, 0) {
				diagnostics = append(diagnostics, lsp.Diagnostics{
					Range:    v.defined,
					Severity: 2,
					Source:   "dbwf-ls",
					Message:  fmt.Sprintf("`%s` is declared but not used anywhere.", k),
				})
			} else if v.defined == lsp.LineRange(0, 0, 0) && v.lastReferred != lsp.LineRange(0, 0, 0) {
				diagnostics = append(diagnostics, lsp.Diagnostics{
					Range:    v.lastReferred,
					Severity: 1,
					Source:   "dbwf-ls",
					Message:  fmt.Sprintf("`%s` is not declared but not used in at least 1 task.", k),
				})
			}
		}
	}

	slices.SortFunc(diagnostics, func(a, b lsp.Diagnostics) int {
		if cmp.Compare(a.Severity, b.Severity) == 0 {
			return cmp.Compare(a.Message, b.Message)
		}
		return cmp.Compare(a.Severity, b.Severity)
	})

	return diagnostics
}
