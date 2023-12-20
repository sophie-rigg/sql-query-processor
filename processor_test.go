package sqlqueryprocessor

import (
	"embed"
	"fmt"
	"maps"
	"testing"

	"github.com/nleof/goyesql"
)

//go:embed queries/*.sql
var queries embed.FS

//go:embed duplicate-names/*.sql
var duplicateNames embed.FS

func TestGetQueries(t *testing.T) {
	type args struct {
		queries            embed.FS
		queryDirectory     string
		queryPreProcessors map[goyesql.Tag]func(string) string
	}
	tests := []struct {
		name    string
		args    args
		want    goyesql.Queries
		wantErr bool
	}{
		{
			name: "should return queries",
			args: args{
				queries:        queries,
				queryDirectory: "queries",
				queryPreProcessors: map[goyesql.Tag]func(string) string{
					"other-query": func(query string) string {
						return fmt.Sprintf(query, "test")
					},
				},
			},
			want: goyesql.Queries{
				"other-query":  "SELECT * FROM other_test WHERE id = $1;",
				"test-query-1": "SELECT * FROM test WHERE id = $1;",
				"test-query-2": "UPDATE test SET name = $1 WHERE id = $2;",
			},
			wantErr: false,
		},
		{
			name: "should return error when duplicate tag found",
			args: args{
				queries:        duplicateNames,
				queryDirectory: "duplicate-names",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQueries(tt.args.queries, tt.args.queryDirectory, tt.args.queryPreProcessors)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("GetQueries() got = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}
