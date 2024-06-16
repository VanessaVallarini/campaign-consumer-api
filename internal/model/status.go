package model

type OwnerStatus string

const (
	ActiveOwner   OwnerStatus = "ACTIVE"
	InactiveOwner OwnerStatus = "INACTIVE"
)

type SlugStatus string

const (
	ActiveSlug   SlugStatus = "ACTIVE"
	InactiveSlug SlugStatus = "INACTIVE"
)

type RegionStatus string

const (
	ActiveRegion   RegionStatus = "ACTIVE"
	InactiveRegion RegionStatus = "INACTIVE"
)
