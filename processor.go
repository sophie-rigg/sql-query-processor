package sqlqueryprocessor

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/nleof/goyesql"
)

func GetQueries(queries embed.FS, queryDirectory string, queryPreProcessors map[goyesql.Tag]func(string) string) (goyesql.Queries, error) {
	sqlFiles, err := os.ReadDir(queryDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read queries directory: %w", err)
	}

	q := make(map[goyesql.Tag]string)
	for _, file := range sqlFiles {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		queryBytes, err := queries.ReadFile(queryDirectory + "/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read query file %s: %w", file.Name(), err)
		}
		fileQueries := goyesql.MustParseBytes(queryBytes)
		for tag, query := range fileQueries {
			if _, ok := q[tag]; ok {
				return nil, fmt.Errorf("duplicate tag found in queries, tag: %s", tag)
			}
			if queryPreProcessors != nil {
				if preProcessor, ok := queryPreProcessors[tag]; ok {
					query = preProcessor(query)
				}
			}

			q[tag] = query
		}
	}

	return q, nil
}
