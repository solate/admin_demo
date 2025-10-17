package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CasbinRule holds the schema definition for the CasbinRule entity.
type CasbinRule struct {
	ent.Schema
}

// Fields of the CasbinRule.
func (CasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("Ptype").Default("").Comment("策略类型：'p'(权限策略) 或 'g'(角色关系)"),
		field.String("V0").Default("").Comment("用户/角色"),
		field.String("V1").Default("").Comment("角色/资源"),
		field.String("V2").Default("").Comment("域/动作"),
		field.String("V3").Default("").Comment("其他属性"),
		field.String("V4").Default("").Comment("其他属性"),
		field.String("V5").Default("").Comment("其他属性"),
	}
}

// Edges of the CasbinRule.
func (CasbinRule) Edges() []ent.Edge {
	return nil
}

func (CasbinRule) Index() []ent.Index {
	return []ent.Index{
		index.Fields("Ptype", "V0", "V1", "V2", "V3", "V4", "V5").Unique(),
	}
}
