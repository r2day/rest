package rest

// CreateResponse 返回
type CreateResponse struct {
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

	// ProductId 商品id
	ProductId string `form:"product_id" json:"product_id" xml:"product_id"`

	// BrandId 品牌id
	BrandId string `form:"brand_id" json:"brand_id" xml:"brand_id"`
}
