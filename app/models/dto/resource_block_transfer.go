package dto

type ResourceBlock struct {
	ResourceUUID string `json:"resource_uuid" binding:"required"`
	BlockUUID    string `json:"block_uuid" binding:"required"`
}

type Collection struct {
	ResourceBlocks []*ResourceBlock `json:"resource_blocks" binding:"required"`
}
