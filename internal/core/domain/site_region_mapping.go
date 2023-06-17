package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type SiteRegionMappingService interface {
	GetAllSiteRegionMapping() ([]port.SiteRegionMapping, error)
	GetSiteRegionMappingWithPagination(limit int, offset int) ([]port.SiteRegionMapping, int, error)
	GetRegion() (Regions, error)
	CreateCity(siteRegionMapping port.SiteRegionMapping) error
	UpdateCity(id int, siteRegionMapping port.SiteRegionMapping) error
	UpdateRegion(regions UpdateRegionsRequestBody) error
	DeleteCity(id int) error
	DeleteArea(area string) error
	ImportBMASiteRegionMapping(fileLocation string) error
}

type Regions struct {
	Regions []AreaWithCity `json:"regions"`
}

type AreaWithCity struct {
	Area   string                   `json:"area"`
	Cities []port.SiteRegionMapping `json:"cities"`
}

type UpdateRegionsRequestBody struct {
	Area   string   `json:"area" binding:"required"`
	Cities []string `json:"cities" binding:"required"`
}
