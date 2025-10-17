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

// ProductStatistics holds the schema definition for the ProductStatistics entity.
type ProductStatistics struct {
	ent.Schema
}

func (ProductStatistics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "product_statistics"},
		entsql.WithComments(true),
		schema.Comment("商品统计"),
	}
}

// Fields of the ProductStatistics.
func (ProductStatistics) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.String("tenant_code").Comment("租户code"),

		// 商品基础统计
		field.Int("total_products").Default(0).Comment("商品总数"),
		field.Int("active_products").Default(0).Comment("启用商品数"),

		// 库存统计
		field.Int("total_stock").Default(0).Comment("总库存数量"),
		field.Other("total_stock_value", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("总库存价值"),
		field.Int("low_stock_products").Default(0).Comment("低库存商品数"),

		// 入库统计
		field.Int("total_in_quantity").Default(0).Comment("总入库数量"),
		field.Other("total_in_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("总入库金额"),

		// 出库统计
		field.Int("total_out_quantity").Default(0).Comment("总出库数量"),
		field.Other("total_out_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("总出库金额"),

		// 销售统计
		field.Other("total_sales_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(18,4)",
				dialect.Postgres: "numeric(18,4)",
			}).Default(decimal.NewFromFloat(0)).Comment("总销售金额"),
		field.Int("total_sales_quantity").Default(0).Comment("总销售数量"),
	}
}

// Edges of the ProductStatistics.
func (ProductStatistics) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (ProductStatistics) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_code"),
	}
}
