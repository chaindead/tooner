package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pythonToJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		errFn assert.ErrorAssertionFunc
	}{
		{
			name:  "python repaired",
			input: `[{'schema_name': 'public', 'schema_owner': 'pg_database_owner', 'active': True, 'meta': None}]`,
			errFn: assert.NoError,
		},
		{
			name:  "python nested dict",
			input: `{'users': [{'id': 1, 'name': 'Alice', 'verified': True}, {'id': 2, 'name': 'Bob', 'verified': False}], 'count': 2}`,
			errFn: assert.NoError,
		},
		{
			name:  "python with None values",
			input: `[{'host': 'localhost', 'port': 5432, 'ssl': None, 'timeout': None}]`,
			errFn: assert.NoError,
		},
		{
			name:  "python trailing comma",
			input: `[{'a': 1111, 'b': 2,'c': 3,'d': 4,},]`,
			errFn: assert.NoError,
		},
		{
			name: "markdown heading with table",
			input: `# Table
| col1 | col2 |
|------|------|
| a    | b    |`,
			errFn: assert.Error,
		},
		{
			name: "markdown list",
			input: `- item 1
  - sub item A
  - sub item B
- item 2
- item 3`,
			errFn: assert.Error,
		},
		{
			name:  "markdown bold and links",
			input: `**Bold text** and [a link](https://example.com) with *italic*`,
			errFn: assert.Error,
		},
		{
			name:  "xml format",
			input: `<?xml version="1.0"?><root><item key="val"/></root>`,
			errFn: assert.Error,
		},
		{
			name: "docker ps format",
			input: `CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
abc123def456   nginx     "/docker-entrypoint.…"   2 hours ago   Up 2 hours   80/tcp    web`,
			errFn: assert.NoError,
		},
		{
			name:  "python pg format",
			input: `[{'id': UUID('d755ae27-48ae-4881-b012-41002340be53'), 'email': 'alice@example.com', 'name': 'Alice', 'created_at': datetime.datetime(2026, 2, 27, 23, 41, 13, 440912, tzinfo=zoneinfo.ZoneInfo(key='Etc/UTC')), 'updated_at': datetime.datetime(2026, 2, 27, 23, 41, 13, 440912, tzinfo=zoneinfo.ZoneInfo(key='Etc/UTC'))}, {'id': UUID('ef607c8c-7055-4ad3-bf1e-11f4aeb9e619'), 'email': 'bob@example.com', 'name': 'Bob', 'created_at': datetime.datetime(2026, 2, 27, 23, 41, 13, 440912, tzinfo=zoneinfo.ZoneInfo(key='Etc/UTC')), 'updated_at': datetime.datetime(2026, 2, 27, 23, 41, 13, 440912, tzinfo=zoneinfo.ZoneInfo(key='Etc/UTC'))}]`,
			errFn: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rjson, err := jsonRepair(tt.input)
			tt.errFn(t, err)

			fmt.Println(rjson)
		})
	}
}
