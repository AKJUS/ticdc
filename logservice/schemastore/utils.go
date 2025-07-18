// Copyright 2024 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package schemastore

import (
	"strings"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/pkg/meta/model"
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/format"
	"go.uber.org/zap"
)

// transform ddl query based on sql mode.
func transformDDLJobQuery(job *model.Job) (string, error) {
	p := parser.New()
	// We need to use the correct SQL mode to parse the DDL query.
	// Otherwise, the parser may fail to parse the DDL query.
	// For example, it is needed to parse the following DDL query:
	//  `alter table "t" add column "c" int default 1;`
	// by adding `ANSI_QUOTES` to the SQL mode.
	p.SetSQLMode(job.SQLMode)
	stmts, _, err := p.Parse(job.Query, job.Charset, job.Collate)
	if err != nil {
		return "", errors.Trace(err)
	}
	var result string
	buildQuery := func(stmt ast.StmtNode) (string, error) {
		var sb strings.Builder
		// translate TiDB feature to special comment
		restoreFlags := format.RestoreTiDBSpecialComment
		// escape the keyword
		restoreFlags |= format.RestoreNameBackQuotes
		// upper case keyword
		restoreFlags |= format.RestoreKeyWordUppercase
		// wrap string with single quote
		restoreFlags |= format.RestoreStringSingleQuotes
		// remove placement rule
		restoreFlags |= format.SkipPlacementRuleForRestore
		// force disable ttl
		restoreFlags |= format.RestoreWithTTLEnableOff
		if err = stmt.Restore(format.NewRestoreCtx(restoreFlags, &sb)); err != nil {
			return "", errors.Trace(err)
		}
		return sb.String(), nil
	}
	if len(stmts) > 1 {
		results := make([]string, 0, len(stmts))
		for _, stmt := range stmts {
			query, err := buildQuery(stmt)
			if err != nil {
				return "", errors.Trace(err)
			}
			results = append(results, query)
		}
		result = strings.Join(results, ";")
	} else {
		result, err = buildQuery(stmts[0])
		if err != nil {
			return "", errors.Trace(err)
		}
	}

	log.Info("transform ddl query to result",
		zap.String("DDL", job.Query),
		zap.String("charset", job.Charset),
		zap.String("collate", job.Collate),
		zap.String("result", result))

	return result, nil
}

// isSplitable returns whether the table is eligible for split in all sinks
// Only the table with pk and no uk can be splitted in all sinks.
func isSplitable(tableInfo *model.TableInfo) bool {
	if tableInfo.GetPkColInfo() == nil {
		return false
	}

	indices := tableInfo.Indices
	for _, index := range indices {
		if index.Primary {
			continue
		}
		if index.Unique {
			return false
		}
	}
	return true
}
