package service

import (
	"encoding/json"
	"errors"

	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/apikey"
)

type ProjectService struct {
	projectRepo  *repository.ProjectRepo
	customerRepo *repository.CustomerRepo
}

func NewProjectService(pr *repository.ProjectRepo, cr *repository.CustomerRepo) *ProjectService {
	return &ProjectService{projectRepo: pr, customerRepo: cr}
}

type CreateProjectReq struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Modules     []string `json:"modules"`
}

func (s *ProjectService) Create(customerID uint, req CreateProjectReq) (*model.MiniappProject, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, errors.New("客户不存在")
	}
	count, _ := s.projectRepo.CountByCustomerID(customerID)
	if int(count) >= customer.MaxProjects {
		return nil, errors.New("已达到项目数量上限，请升级套餐")
	}

	apiKey, _ := apikey.GenerateKey()
	apiSecret, _ := apikey.GenerateSecret()

	// Ensure base module is always included
	modules := req.Modules
	hasBase := false
	for _, m := range modules {
		if m == "base" {
			hasBase = true
			break
		}
	}
	if !hasBase {
		modules = append([]string{"base"}, modules...)
	}

	modulesJSON, _ := json.Marshal(modules)
	brandConfigJSON, _ := json.Marshal(map[string]interface{}{
		"appName":      req.Name,
		"primaryColor": "#4caf50",
	})

	project := &model.MiniappProject{
		CustomerID:  customerID,
		Name:        req.Name,
		Description: req.Description,
		Modules:     string(modulesJSON),
		APIKey:      apiKey,
		APISecret:   apikey.HashSecret(apiSecret),
		BrandConfig: string(brandConfigJSON),
		Status:      "draft",
	}
	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}
	project.APISecret = apiSecret // Return plaintext secret (only this once)
	return project, nil
}

func (s *ProjectService) List(customerID uint) ([]model.MiniappProject, error) {
	return s.projectRepo.FindByCustomerID(customerID)
}

func (s *ProjectService) Get(id uint) (*model.MiniappProject, error) {
	return s.projectRepo.FindByID(id)
}

func (s *ProjectService) UpdateModules(id uint, modules []string, brandConfig map[string]interface{}) error {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if modules != nil {
		modulesJSON, _ := json.Marshal(modules)
		p.Modules = string(modulesJSON)
	}
	if brandConfig != nil {
		brandConfigJSON, _ := json.Marshal(brandConfig)
		p.BrandConfig = string(brandConfigJSON)
	}
	return s.projectRepo.Update(p)
}
