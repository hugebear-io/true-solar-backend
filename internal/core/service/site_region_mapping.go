package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/xuri/excelize/v2"
)

type siteRegionMappingService struct {
	repo port.SiteRegionMappingRepoPort
}

func NewSiteRegionMappingService(repo port.SiteRegionMappingRepoPort) domain.SiteRegionMappingService {
	return &siteRegionMappingService{repo: repo}
}

func (s siteRegionMappingService) GetAllSiteRegionMapping() ([]port.SiteRegionMapping, error) {
	siteRegionMapping, err := s.repo.GetAllSiteRegionMapping()
	if err != nil {
		return []port.SiteRegionMapping{}, err
	}
	return siteRegionMapping, nil
}

func (s siteRegionMappingService) GetSiteRegionMappingWithPagination(limit int, offset int) ([]port.SiteRegionMapping, int, error) {
	total, err := s.repo.TotalSiteRegionMapping()
	if err != nil {
		return []port.SiteRegionMapping{}, 0, err
	}

	siteRegionMapping, err := s.repo.GetSiteRegionMappingWithPagination(limit, offset)
	if err != nil {
		return []port.SiteRegionMapping{}, total, err
	}
	return siteRegionMapping, total, nil
}

func (s siteRegionMappingService) GetRegion() (domain.Regions, error) {
	siteRegionMapping, err := s.repo.GetAllSiteRegionMapping()
	if err != nil {
		return domain.Regions{}, err
	}

	mapAreaCities := make(map[string][]port.SiteRegionMapping)
	for _, site := range siteRegionMapping {
		if site.Area != "" {
			mapAreaCities[site.Area] = append(mapAreaCities[site.Area], site)
		}
	}

	areaNullCity, err := s.repo.GetAreaNullCity()
	if err != nil {
		return domain.Regions{}, err
	}

	for _, site := range areaNullCity {
		mapAreaCities[site.Area] = append(mapAreaCities[site.Area], site)
	}

	var region domain.Regions
	for area, cities := range mapAreaCities {
		if len(cities) == 0 {
			cities = make([]port.SiteRegionMapping, 0)
		}

		region.Regions = append(region.Regions, domain.AreaWithCity{
			Area:   area,
			Cities: cities,
		})
	}

	return region, nil
}

func (s siteRegionMappingService) CreateCity(siteRegionMapping port.SiteRegionMapping) error {
	err := s.repo.CreateCity(siteRegionMapping)
	if err != nil {
		return err
	}
	return nil
}

func (s siteRegionMappingService) UpdateCity(id int, siteRegionMapping port.SiteRegionMapping) error {
	err := s.repo.UpdateCity(id, siteRegionMapping)
	if err != nil {
		return err
	}
	return nil
}

func (s siteRegionMappingService) UpdateRegion(regions domain.UpdateRegionsRequestBody) error {
	err := s.repo.UpdateCityToNullArea(strings.ToUpper(regions.Area))
	if err != nil {
		return err
	}

	if len(regions.Cities) == 0 {
		id := uuid.New()

		cityCode := strings.ToUpper(fmt.Sprintf("EMPTY-%s", id.String()))
		err = s.repo.CreateCity(port.SiteRegionMapping{
			Code: cityCode,
		})
		if err != nil {
			return err
		}

		regions.Cities = append(regions.Cities, cityCode)
	}

	cityCodeListString := strings.ToUpper(fmt.Sprintf("'%s'", strings.Join(regions.Cities, "','")))
	err = s.repo.UpdateSiteRegionMapping(regions.Area, cityCodeListString)
	if err != nil {
		return err
	}

	return nil
}

func (s siteRegionMappingService) DeleteCity(id int) error {
	err := s.repo.DeleteCity(id)
	if err != nil {
		return err
	}
	return nil
}

func (s siteRegionMappingService) DeleteArea(area string) error {
	err := s.repo.UpdateCityToNullArea(strings.ToUpper(area))
	if err != nil {
		return err
	}
	return nil
}

func (s siteRegionMappingService) ImportBMASiteRegionMapping(fileLocation string) error {
	f, err := excelize.OpenFile(fileLocation)
	if err != nil {
		return err
	}

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		err := errors.New("excel file must have at least one sheet")
		return err
	}

	firstSheet := sheets[0]
	rows, err := f.GetRows(firstSheet)
	if err != nil {
		return err
	}

	siteRegionMappings := make([]port.SiteRegionMapping, 0)
	for i, row := range rows {
		if i == 0 { // ignore header (Site_ID, NEW_BMA_ZONE_Y2021)
			continue
		}

		if len(row) == 2 { // BKA0001, BMA VI Central 2
			siteID := strings.TrimSpace(row[0])
			area := strings.ReplaceAll(strings.TrimSpace(row[1]), " ", "-")

			siteRegionMappings = append(siteRegionMappings, port.SiteRegionMapping{
				Code: siteID,
				Name: "กรุงเทพมหานคร",
				Area: area,
			})
		}
	}

	return s.repo.CreateBatchSiteRegionMappings(siteRegionMappings)
}
