package models

import "gorm.io/gorm"

type ProjectTree struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	AuthorID uint           `json:"author_id"`
	Type     string         `json:"type"`
	Method   string         `json:"method"`
	ParentID *uint          `json:"parent_id"`
	Children []*ProjectTree `json:"children"`
}

func (pt *ProjectTree) GetProjectTree(db *gorm.DB, projectID uint) ([]*ProjectTree, error) {
	var folders []ProjectFolder

	err := db.Debug().Where("project_id = ?", projectID).Find(&folders).Error
	if err != nil {
		return nil, err
	}

	var apis []ProjectApi
	err = db.Debug().Where("folder_id IN ?", getFolderIDs(folders)).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	folderMap := make(map[uint]*ProjectTree)
	for _, folder := range folders {
		folderMap[folder.ID] = &ProjectTree{
			ID:       folder.ID,
			Name:     folder.Name,
			AuthorID: folder.AuthorID,
			ParentID: folder.ParentID,
			Type:     "folder",
			Method:   "",
			Children: []*ProjectTree{},
		}
	}

	for _, api := range apis {
		folder := folderMap[api.FolderID]
		if folder != nil {
			folder.Children = append(folder.Children, &ProjectTree{
				ID:       api.ID,
				Name:     api.Name,
				AuthorID: api.AuthorID,
				Type:     "api",
				Method:   api.Method,
				ParentID: &api.FolderID,
				Children: nil, // APIs do not have children
			})
		}
	}

	var tree []*ProjectTree
	for _, folder := range folders {
		if folder.ParentID == nil {
			tree = append(tree, folderMap[folder.ID])
		} else {
			if parentFolder, exists := folderMap[*folder.ParentID]; exists {
				parentFolder.Children = append(parentFolder.Children, folderMap[folder.ID])
			}
		}
	}

	return tree, nil
}

func getFolderIDs(folders []ProjectFolder) []uint {
	var ids []uint
	for _, folder := range folders {
		ids = append(ids, folder.ID)
	}
	return ids
}
