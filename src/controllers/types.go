package controllers

type IDsField struct {
	IDs []int64 `json:"ids" form:"ids" uri:"ids" binding:"required"`
}

type IDField struct {
	ID int64 `json:"id" form:"id" uri:"id" binding:"required"`
}
