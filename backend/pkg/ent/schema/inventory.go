package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
)

// Inventory holds the schema definition for the Inventory entity.
type Inventory struct {
	ent.Schema
}

func (Inventory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "inventories"},
		entsql.WithComments(true),
		schema.Comment("库存记录"),
	}
}

// Fields of the Inventory.
func (Inventory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		field.String("inventory_id").Unique().Comment("库存记录ID"),
		field.String("product_id").NotEmpty().Comment("商品ID"),
		field.String("operation_type").NotEmpty().Comment("操作类型: in-入库, out-出库"),
		field.Int("quantity").Comment("操作数量"),
		field.Other("unit_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("单价"),
		field.Other("total_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("总金额"),
		field.String("operator_id").Comment("操作人ID"),
		field.String("remark").Default("").Comment("备注"),
		field.Int64("operation_time").Comment("操作时间"),

		// 库存快照
		field.Int("before_stock").Default(0).Comment("操作前库存"),
		field.Int("after_stock").Default(0).Comment("操作后库存"),
	}
}

// Edges of the Inventory.
func (Inventory) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Inventory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_id"),
		index.Fields("operation_type"),
		index.Fields("operator_id"),
		index.Fields("operation_time"),
	}
}
