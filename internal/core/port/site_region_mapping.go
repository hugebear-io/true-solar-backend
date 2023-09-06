package port

import "time"

type SiteRegionMappingRepoPort interface {
	GetAllSiteRegionMapping() ([]SiteRegionMapping, error)
	GetSiteRegionMappingWithPagination(limit int, offset int) ([]SiteRegionMapping, error)
	GetAreaNullCity() ([]SiteRegionMapping, error)
	CreateCity(siteRegionMapping SiteRegionMapping) error
	UpdateCity(id int, siteRegionMapping SiteRegionMapping) error
	UpdateCityToNullArea(area string) error
	DeleteCity(id int) error
	CreateBatchSiteRegionMappings(siteRegionMappings []SiteRegionMapping) error
	UpdateSiteRegionMapping(area string, cityCodeListString string) error
	TotalSiteRegionMapping() (int, error)
}

type SiteRegionMapping struct {
	ID        int       `json:"id"`
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Area      *string   `json:"area"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
