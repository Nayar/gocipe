package util

import (
	"fmt"
	"strconv"
	"strings"
)

//Postgres provides postgres compliant implementation of generated CRUD code
type Postgres struct {
	Entity
}

//SQLInsert returns SQL query for SQL Insert
func (s Postgres) SQLInsert() string {
	var (
		fields, placeholders []string
		index                = 1
	)

	for _, f := range s.Fields {
		fields = append(fields, f.schema.Field)
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`"%s"`, p.ThisID))
			placeholders = append(placeholders, "$"+strconv.Itoa(index))
			index++
		}
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		s.Table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)
}

//SQLGet returns SQL query for SQL Get
func (s Postgres) SQLGet() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`t."%s"`, p.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t WHERE t."id" = $1 ORDER BY t."id" ASC`,
		strings.Join(fields, ", "),
		s.Table,
	)
}

//SQLList returns SQL query for SQL List
func (s Postgres) SQLList() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne:
			fields = append(fields, fmt.Sprintf(`t."%s"`, p.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT %s FROM %s t`,
		strings.Join(fields, ", "),
		s.Table,
	)
}

//OrderList returns list of fields to be used for list statement
func (s Postgres) OrderList() string {
	var fields []string

	for _, f := range s.Fields {
		fields = append(fields, f.schema.Field)
	}

	return `"` + strings.Join(fields, `","`) + `"`
}

//SQLDeleteSingle returns SQL query for SQL Delete Single
func (s Postgres) SQLDeleteSingle() string {
	return fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1`,
		s.Table,
	)
}

//SQLDeleteMany returns SQL query for SQL Delete Many
func (s Postgres) SQLDeleteMany() string {
	return fmt.Sprintf(
		`DELETE FROM %s`,
		s.Table,
	)
}

//SQLDeleteManyJoin returns SQL query for SQL Delete Many Join
func (s Postgres) SQLDeleteManyJoin() []string {
	var query []string

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyManyOwner:
			query = append(query, fmt.Sprintf("DELETE FROM %s WHERE %s IN (SELECT id FROM %s", p.JoinTable, p.ThatID, s.Table))
		}
	}

	return query
}

//SQLUpdate returns SQL query for SQL Update
func (s Postgres) SQLUpdate() string {
	var (
		fields []string
		index  = 2
	)

	for _, f := range s.Fields {
		fields = append(fields, fmt.Sprintf(`"%s" = $%d`, f.schema.Field, index))
		index++
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne, RelationshipTypeOneOne:
			fields = append(fields, fmt.Sprintf(`"%s" = $%d`, p.ThisID, index))
			index++
		}
	}

	return fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = $1`,
		s.Table,
		strings.Join(fields, ", "),
	)
}

//SQLMerge returns SQL query for SQL Merge
func (s Postgres) SQLMerge() string {
	var (
		updates, inserts, placeholders []string
		index                          = 1
	)

	for _, f := range s.Fields {
		inserts = append(inserts, f.schema.Field)
		updates = append(updates, fmt.Sprintf(`"%s" = $%d`, f.schema.Field, index))
		placeholders = append(placeholders, "$"+strconv.Itoa(index))
		index++
	}

	for _, p := range s.Relationships {
		switch p.Type {
		case RelationshipTypeManyOne, RelationshipTypeOneOne:
			inserts = append(inserts, fmt.Sprintf(`"%s"`, p.ThisID))
			updates = append(updates, fmt.Sprintf(`"%s" = $%d`, p.ThisID, index))
			placeholders = append(placeholders, "$"+strconv.Itoa(index))
			index++
		}
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) 
		ON CONFLICT (id) DO UPDATE SET %s`,
		s.Table,
		strings.Join(inserts, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)
}

//SQLLoadManyMany returns SQL query for SQL Load many many
func (s Postgres) SQLLoadManyMany(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	return fmt.Sprintf(
		`SELECT j.%s, %s FROM %s t 
		INNER JOIN %s j ON t.id = j.%s
		WHERE j.%s IN`,
		rel.ThisID,
		strings.Join(fields, ", "),
		related.Table,
		rel.JoinTable,
		rel.ThatID,
		rel.ThisID,
	)
}

//SQLLoadManyOne returns SQL query for SQL Load many one
func (s Postgres) SQLLoadManyOne(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, rel := range related.Relationships {
		if rel.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`t."%s"`, rel.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT t."id", %s FROM %s t WHERE t."id" IN`,
		strings.Join(fields, ", "),
		related.Table,
	)
}

//SQLLoadOneMany returns SQL query for SQL Load one many
func (s Postgres) SQLLoadOneMany(rel Relationship) string {
	var fields []string

	related := rel.related

	for _, f := range related.Fields {
		fields = append(fields, fmt.Sprintf(`t."%s"`, f.schema.Field))
	}

	for _, rel := range related.Relationships {
		if rel.Type == RelationshipTypeManyOne {
			fields = append(fields, fmt.Sprintf(`t."%s"`, rel.ThisID))
		}
	}

	return fmt.Sprintf(
		`SELECT t."%s", %s FROM %s t WHERE t."%s" IN`,
		rel.ThisID,
		strings.Join(fields, ", "),
		related.Table,
		rel.ThisID,
	)
}

//SQLSaveManyManyOwnerDelete returns SQL query for SQL Save many many owner
func (s Postgres) SQLSaveManyManyOwnerDelete(rel Relationship) string {
	return fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		rel.JoinTable,
		rel.ThatID,
	)
}

//SQLSaveManyManyOwnerInsert returns SQL query for SQL Save many many owner
func (s Postgres) SQLSaveManyManyOwnerInsert(rel Relationship) string {
	return fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES ($1, $2)`,
		rel.JoinTable,
		rel.ThatID,
		rel.ThisID,
	)
}