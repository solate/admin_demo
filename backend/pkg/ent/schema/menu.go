package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Menu holds the schema definition for the Menu entity.
type Menu struct {
	ent.Schema
}

func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_menus"},
		entsql.WithComments(true),
		schema.Comment("菜单"),
	}
}

// Fields of the Menu.
func (Menu) Fields() []ent.Field {
	return []ent.Field{
		// 基础字段
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		// 菜单特有字段
		field.String("menu_id").Unique().Comment("菜单ID"),
		field.String("code").Unique().NotEmpty().Comment("菜单code"),
		field.String("parent_id").Default("").Comment("父菜单ID"),
		field.String("name").NotEmpty().Comment("菜单名称"),
		field.String("path").Default("").Comment("路由路径"),
		field.String("component").Default("").Comment("前端组件路径"),
		field.String("redirect").Default("").Comment("重定向路径"),
		field.String("icon").Default("").Comment("菜单图标"),
		field.Int("sort").Default(0).Comment("排序号"),
		field.String("type").Default("").Comment("菜单类型 dir/menu/button"),
		field.Int("status").Default(1).Comment("状态 1:启用 2:禁用"),
	}
}

// Edges of the Menu.
func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		// 自关联边，用于构建菜单树
		// edge.To("children", Menu.Type).From("parent").Unique().Field("parent_id"),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("menu_id"),
		index.Fields("parent_id"),
		// index.Fields("tenant_code", "menu_id").Unique(),
	}
}
