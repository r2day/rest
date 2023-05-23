package rest

// SimpleResponse 返回内容
type SimpleResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// FilterRequest Binding from JSON
// 常用的过滤字段
type FilterRequest struct {
	// id 列表id
	Id []string `form:"id" json:"id" xml:"id"`

	// Status 状态
	Status string `form:"status" json:"status" xml:"status"`

	// CategoryId 分类id
	CategoryId string `form:"category_id" json:"category_id" xml:"category_id"`
}

// CommonFilter 常用的过滤器
// 可以根据最常使用的一些过滤条件作为基本
// 后续不断更新
// 一般分类应该按照从大到小依次编排，以便选择合适的使用
type CommonFilter struct {
	// Status 状态
	Status bool `form:"status" json:"status" xml:"status"`

	// Company 公司
	Company string `form:"company" json:"company" xml:"company"`

	// Organization 组织
	Organization string `form:"organization_id" json:"organization_id" xml:"organization_id"`

	// CategoryId 分类
	CategoryID string `form:"category_id" json:"category_id" xml:"category_id"`

	// CategoryName 分类
	CategoryName string `form:"category_name" json:"category_name" xml:"category_name"`
	// GroupID 分组
	GroupID string `form:"group_id" json:"group_id" xml:"group_id"`

	// Tag 标签
	Tag string `form:"tag" json:"tag" xml:"tag"`

	// BrandId 品牌id
	BrandID string `form:"brand_id" json:"brand_id" xml:"brand_id"`

	// ProductId 商品id
	ProductID string `form:"product_id" json:"product_id" xml:"product_id"`

	// Gender 性别
	Gender string `form:"gender" json:"gender" xml:"gender"`

	// From 来源
	From string `form:"from" json:"from" xml:"from"`

	// OrderCategory 订单分类
	OrderCategory string `form:"order_category" json:"order_category" xml:"order_category"`
	// Level 等级
	Level string `form:"level" json:"level" xml:"level"`
}
