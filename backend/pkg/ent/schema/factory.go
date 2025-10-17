package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Factory holds the schema definition for the Factory entity.
type Factory struct {
	ent.Schema
}

func (Factory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "factories"},
		entsql.WithComments(true),
		schema.Comment("工厂"),
	}
}

// Fields of the Factory.
func (Factory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		field.String("factory_id").Unique().Comment("工厂ID"),
		field.String("factory_name").NotEmpty().Comment("工厂名称"),
		field.String("address").Default("").Comment("工厂地址"),
		field.String("contact_phone").Default("").Comment("联系电话"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),
	}
}

// Edges of the Factory.
func (Factory) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Factory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("factory_name"),
		index.Fields("tenant_code"),
	}
}
