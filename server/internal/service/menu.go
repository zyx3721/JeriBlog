package service

import (
	"errors"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"

	"gorm.io/gorm"
)

// MenuService 菜单服务
type MenuService struct {
	menuRepo    *repository.MenuRepository
	fileService *FileService
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuRepo *repository.MenuRepository, fileService *FileService) *MenuService {
	return &MenuService{
		menuRepo:    menuRepo,
		fileService: fileService,
	}
}

// Create 创建菜单
func (s *MenuService) Create(req *dto.MenuCreateRequest) (*dto.MenuResponse, error) {
	// 如果有父菜单，检查父菜单是否存在且类型一致
	if req.ParentID != nil && *req.ParentID > 0 {
		parent, err := s.menuRepo.Get(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父菜单不存在")
			}
			return nil, err
		}
		// 检查类型是否一致
		if parent.Type != req.Type {
			return nil, errors.New("子菜单类型必须与父菜单类型一致")
		}
		// 检查父菜单自己也不能是子菜单（限制只支持两级）
		if parent.ParentID != nil {
			return nil, errors.New("子菜单无法创建子菜单")
		}
		// 子菜单必须填写链接地址
		if req.URL == "" {
			return nil, errors.New("子菜单必须填写链接地址")
		}
	}

	menu := &model.Menu{
		Type:      req.Type,
		ParentID:  req.ParentID,
		Title:     req.Title,
		URL:       req.URL,
		Icon:      req.Icon,
		Sort:      req.Sort,
		IsEnabled: req.IsEnabled,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	// 标记图标为使用中
	if s.fileService != nil && req.Icon != "" {
		_ = s.fileService.MarkAsUsed(req.Icon)
	}

	return s.toMenuResponse(menu), nil
}

// Update 更新菜单
func (s *MenuService) Update(id uint, req *dto.MenuUpdateRequest) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		return nil, err
	}

	// 如果有父菜单，检查父菜单是否存在且类型一致
	if req.ParentID != nil && *req.ParentID > 0 {
		// 不能将自己设为父菜单
		if *req.ParentID == id {
			return nil, errors.New("不能将自己设为父菜单")
		}

		parent, err := s.menuRepo.Get(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父菜单不存在")
			}
			return nil, err
		}

		// 检查类型是否一致
		if parent.Type != req.Type {
			return nil, errors.New("子菜单类型必须与父菜单类型一致")
		}

		// 检查父菜单自己也不能是子菜单（限制只支持两级）
		if parent.ParentID != nil {
			return nil, errors.New("子菜单无法创建子菜单")
		}
	}

	// 如果当前菜单有子菜单，不能将它设为子菜单
	if req.ParentID != nil && *req.ParentID > 0 {
		hasChildren, err := s.menuRepo.HasChildren(id)
		if err != nil {
			return nil, err
		}
		if hasChildren {
			return nil, errors.New("该菜单有子菜单，不能设为子菜单")
		}
		// 子菜单必须填写链接地址
		if req.URL == "" {
			return nil, errors.New("子菜单必须填写链接地址")
		}
	}

	// 保存旧图标
	oldIcon := menu.Icon

	menu.Type = req.Type
	menu.ParentID = req.ParentID
	menu.Title = req.Title
	menu.URL = req.URL
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.IsEnabled = req.IsEnabled

	if err := s.menuRepo.Update(menu); err != nil {
		return nil, err
	}

	// 处理图标变化
	if s.fileService != nil && oldIcon != req.Icon {
		if oldIcon != "" {
			_ = s.fileService.MarkAsUnused(oldIcon)
		}
		if req.Icon != "" {
			_ = s.fileService.MarkAsUsed(req.Icon)
		}
	}

	return s.toMenuResponse(menu), nil
}

// Delete 删除菜单
func (s *MenuService) Delete(id uint, req *dto.MenuDeleteRequest) error {
	// 获取菜单信息，用于标记文件为未使用
	menu, err := s.menuRepo.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return err
	}

	// 检查是否有子菜单
	hasChildren, err := s.menuRepo.HasChildren(id)
	if err != nil {
		return err
	}

	if hasChildren {
		// 如果有子菜单，但前端没有指定处理方式，返回错误
		if req == nil || req.ChildrenAction == "" {
			return errors.New("该菜单下还有子菜单，请指定如何处理子菜单")
		}

		// 根据前端指定的方式处理子菜单
		switch req.ChildrenAction {
		case "delete":
			// 删除所有子菜单
			// 先获取子菜单列表，用于标记图标为未使用
			children, err := s.menuRepo.GetChildren(id)
			if err != nil {
				return err
			}

			// 删除子菜单
			if err := s.menuRepo.DeleteChildren(id); err != nil {
				return err
			}

			// 标记子菜单的图标为未使用
			if s.fileService != nil {
				for _, child := range children {
					if child.Icon != "" {
						_ = s.fileService.MarkAsUnused(child.Icon)
					}
				}
			}

		case "upgrade":
			// 将子菜单升级为主菜单
			if err := s.menuRepo.UpgradeChildren(id); err != nil {
				return err
			}

		default:
			return errors.New("无效的子菜单处理方式")
		}
	}

	// 删除父菜单
	if err := s.menuRepo.Delete(id); err != nil {
		return err
	}

	// 标记父菜单的图标为未使用
	if s.fileService != nil && menu.Icon != "" {
		_ = s.fileService.MarkAsUnused(menu.Icon)
	}

	return nil
}

// Get 根据ID获取菜单
func (s *MenuService) Get(id uint) (*dto.MenuResponse, error) {
	menu, err := s.menuRepo.GetWithParent(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		return nil, err
	}

	return s.toMenuResponse(menu), nil
}

// ListForWeb 获取菜单树（前台，树形结构，只返回启用的菜单）
func (s *MenuService) ListForWeb(menuType string) ([]dto.MenuTreeNode, error) {
	// 前台只返回启用的菜单
	enabled := true
	menus, err := s.menuRepo.GetTree(menuType, &enabled)
	if err != nil {
		return nil, err
	}

	nodes := make([]dto.MenuTreeNode, 0, len(menus))
	for _, menu := range menus {
		nodes = append(nodes, *s.toMenuTreeNode(&menu))
	}
	return nodes, nil
}

// List 获取菜单树（后台管理，树形结构，不过滤启用状态）
func (s *MenuService) List(menuType string) ([]dto.MenuTreeNode, error) {
	// 后台管理不过滤启用状态，显示所有菜单
	menus, err := s.menuRepo.GetTree(menuType, nil)
	if err != nil {
		return nil, err
	}

	nodes := make([]dto.MenuTreeNode, 0, len(menus))
	for _, menu := range menus {
		nodes = append(nodes, *s.toMenuTreeNode(&menu))
	}
	return nodes, nil
}

// toMenuResponse 转换为菜单响应（单条数据）
func (s *MenuService) toMenuResponse(menu *model.Menu) *dto.MenuResponse {
	return &dto.MenuResponse{
		ID:        menu.ID,
		Type:      menu.Type,
		ParentID:  menu.ParentID,
		Title:     menu.Title,
		URL:       menu.URL,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		IsEnabled: menu.IsEnabled,
	}
}

// toMenuTreeNode 转换为菜单树节点
func (s *MenuService) toMenuTreeNode(menu *model.Menu) *dto.MenuTreeNode {
	node := &dto.MenuTreeNode{
		ID:        menu.ID,
		Type:      menu.Type,
		ParentID:  menu.ParentID,
		Title:     menu.Title,
		URL:       menu.URL,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		IsEnabled: menu.IsEnabled,
	}

	if len(menu.Children) > 0 {
		node.Children = make([]dto.MenuTreeNode, 0, len(menu.Children))
		for _, child := range menu.Children {
			node.Children = append(node.Children, *s.toMenuTreeNode(&child))
		}
	}

	return node
}
