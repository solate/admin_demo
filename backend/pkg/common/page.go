package common

// PageBuilder 延迟分页查询器，用于分页查询
type PageBuilder[T any] struct {
	lastEntity *T
	callback   func(lastEntity *T, limit int) ([]*T, error)
}

// NewPageBuilder 创建一个分页查询器
func NewPageBuilder[T any](call func(lastEntity *T, limit int) ([]*T, error)) *PageBuilder[T] {
	return &PageBuilder[T]{
		callback: call,
	}
}

// Next 查询下一页数据，不支持并发调用
func (pb *PageBuilder[T]) Next(limit int) ([]*T, error) {
	data, err := pb.callback(pb.lastEntity, limit)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		pb.lastEntity = data[len(data)-1]
	}
	return data, nil
}

// Offset 计算offset
func Offset(page, pageSize int) int {
	offset := int(0)
	limit := pageSize
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	return offset
}
