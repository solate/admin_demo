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

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

func (Product) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "bus_products"},
		entsql.WithComments(true),
		schema.Comment("商品"),
	}
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		field.String("product_id").Unique().Comment("商品ID"),
		field.String("product_name").NotEmpty().Comment("商品名称"),
		field.String("unit").Default("").Comment("单位"),

		// 价格相关字段 - 使用decimal类型
		field.Other("purchase_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("采购价格"),
		field.Other("sale_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("销售价格"),

		// 库存相关字段
		field.Int("current_stock").Default(0).Comment("当前库存"),
		field.Int("min_stock").Default(0).Comment("最小库存预警"),

		// 状态字段
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),

		// 关联工厂
		field.String("factory_id").Optional().Comment("所属工厂ID"),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_name"),
		index.Fields("factory_id"),
		index.Fields("tenant_code"),
	}
}
