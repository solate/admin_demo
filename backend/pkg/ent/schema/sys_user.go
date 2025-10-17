package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type SysUser struct {
	ent.Schema
}

func (SysUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_users"},
		entsql.WithComments(true),
		schema.Comment("用户"),
	}
}

// Fields of the SysUser.
func (SysUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		field.String("user_id").Unique().Comment("用户ID"),
		field.String("user_name").NotEmpty().Default("").Comment("用户名"),
		field.String("pwd_hashed").NotEmpty().Default("").Comment("hash后的密码"),
		field.String("pwd_salt").NotEmpty().Default("").Comment("密码加盐"),
		field.String("token").Default("").Comment("登录后的token信息"),

		field.String("name").Default("").Comment("用户昵称"),
		field.String("avatar").Default("").Comment("头像"),
		field.String("phone").Default("").Comment("电话"),
		field.String("email").Default("").Comment("邮箱"),
		field.Int("sex").Default(0).Comment("性别: 1：男 2：女"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),

		// field.Uint64("role_id").Default(0).Comment("角色ID"),
		// field.Uint64("dept_id").Default(0).Comment("部门ID"),
		// field.Uint64("position_id").Default(0).Comment("岗位ID"),
	}
}

// Edges of the SysUser.
func (SysUser) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (SysUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone"),
		index.Fields("user_id"),
	}
}
