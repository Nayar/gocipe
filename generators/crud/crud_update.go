package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateUpdate produces code for database update of entity (UPDATE)
func generateUpdate(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before, after, sqlfields, structFields []string
		count                                  = 2
	)
	related := make(map[string]string)
	structFields = append(structFields, fmt.Sprintf("&entity.%s", "ID"))

	for _, field := range entity.Fields {
		if field.Property.Name == "UpdatedAt" {
			before = append(before, "entity.UpdatedAt = ptypes.TimestampNow()")
		}

		sqlfields = append(sqlfields, fmt.Sprintf(`"%s" = $%d`, field.Schema.Field, count))

		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("%s, _ := ptypes.Timestamp(entity.%s)", prop, field.Property.Name))
			structFields = append(structFields, fmt.Sprintf("%s", prop))
		} else {
			structFields = append(structFields, fmt.Sprintf("entity.%s", field.Property.Name))
		}
		count++
	}

	for _, rel := range entity.Relationships {
		// No SaveRelated needed:
		// RelationshipTypeManyManyOwner, RelationshipTypeOneMany, RelationshipTypeManyOne, RelationshipTypeOneOne
		if rel.Type == util.RelationshipTypeManyMany || rel.Type == util.RelationshipTypeManyManyOwner {
			related[rel.Name] = fmt.Sprintf("repo.Save%s(ctx, tx, false, entity.ID, entity.%s...)", util.RelFuncName(rel), rel.Name)
		} else if rel.Type == util.RelationshipTypeManyOne || rel.Type == util.RelationshipTypeOneOne {
			sqlfields = append(sqlfields, fmt.Sprintf(`"%s" = $%d`, rel.ThisID, count))
			structFields = append(structFields, fmt.Sprintf("entity.%sID", rel.Name))
			count++
		}
	}

	for _, ref := range entity.References {
		// IDField
		sqlfields = append(sqlfields, fmt.Sprintf(`"%s" = $%d`, ref.IDField.Schema.Field, count))
		structFields = append(structFields, fmt.Sprintf("entity.%s", ref.IDField.Property.Name))
		count++

		// IDType
		sqlfields = append(sqlfields, fmt.Sprintf(`"%s" = $%d`, ref.TypeField.Schema.Field, count))
		structFields = append(structFields, fmt.Sprintf("entity.%s", ref.TypeField.Property.Name))
		count++
	}

	return util.ExecuteTemplate("crud/partials/update.go.tmpl", struct {
		Before        []string
		After         []string
		Related       map[string]string
		EntityName    string
		HasPostHook   bool
		HasPreHook    bool
		SQLFields     string
		StructFields  string
		Table         string
		Relationships []relationship
	}{
		EntityName:    entity.Name,
		Table:         entity.Table,
		Before:        before,
		Related:       related,
		After:         after,
		SQLFields:     strings.Join(sqlfields, ", "),
		StructFields:  strings.Join(structFields, ", "),
		HasPreHook:    entity.CrudHooks.PreSave,
		HasPostHook:   entity.CrudHooks.PostSave,
		Relationships: nil,
	})
}
