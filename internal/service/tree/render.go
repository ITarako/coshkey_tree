package tree

import (
	"fmt"
	"strconv"

	"github.com/ITarako/coshkey_tree/internal/model"
)

func (s Service) renderMain(folders map[int]model.Folder, user *model.User, folderId int) string {
	var html string
	sortedFolders := model.SetClassification(folders, user)
	for category, mainFolders := range sortedFolders {
		html += setCategoryTag(category)
		html += "<ul>"
		for _, mainFolder := range mainFolders {
			html += s.renderMainItem(mainFolder, user, folderId)
		}
		html += "</ul>"
	}

	return html
}

func (s Service) renderMainItem(folder model.Folder, user *model.User, folderId int) string {
	var html string
	if !user.IsAdmin && folder.IsActive {
		li := `<li
			id='jstree-node-%d'
			data-url='%s'
			data-idfolder='%d'
			data-jstree='{"icon":"item-type fa fa-folder fa-lg","opened":false, "selected":%t}'
		>`
		url := s.CoshkeyUrl + "/" + strconv.Itoa(folder.Id)
		html += fmt.Sprintf(li, folder.Id, url, folder.Id, folderId == folder.Id)

		span := "<span class='folders-tytle'>%s</span>"
		html += fmt.Sprintf(span, folder.Title)

		if len(folder.Children) > 0 {
			html += "<ul>"
			for _, child := range model.GetSortedSliceFromMap(folder.Children) {
				html += s.renderMainItem(child, user, folderId)
			}
			html += "</ul>"
		}

		html += "</li>"
	} else if user.IsAdmin {
		li := `<li
			id='jstree-node-%d'
			data-url='%s'
			data-idfolder='%d'
			data-jstree='{"icon":"item-type fa fa-folder fa-lg","opened":false, "selected":%t}'
		>`
		url := s.CoshkeyUrl + "/" + strconv.Itoa(folder.Id)
		html += fmt.Sprintf(li, folder.Id, url, folder.Id, folderId == folder.Id)

		var item string
		if !folder.IsActive {
			item = "<i class='fa fa-minus-square-o deleted' aria-hidden='true' title='Папка удалена'></i> "
		}
		span := "<span class='folders-tytle'>%s%s</span>"
		html += fmt.Sprintf(span, item, folder.Title)

		if len(folder.Children) > 0 {
			html += "<ul>"
			for _, child := range model.GetSortedSliceFromMap(folder.Children) {
				html += s.renderMainItem(child, user, folderId)
			}
			html += "</ul>"
		}

		html += "</li>"
	}

	return html
}

func (s Service) renderFavorite(folders map[int]model.Folder, user *model.User, folderId int) string {
	html := "<ul>"
	li := `<li
		id='jstree-favorite-root'
		data-url='favorite'
		data-idfolder='%d'
		data-jstree='{"icon":"fa fa-star fa-lg","opened":false, "selected":%t}'
>`
	html += fmt.Sprintf(li, model.FolderIdForTree, folderId == model.FolderIdForTree)
	html += "<span class='folders-tytle'>Избранное</span>"

	sortedFolders := model.SetClassification(folders, user)
	for category, favoriteFolders := range sortedFolders {
		html += setCategoryTag(category)
		html += "<ul>"
		for _, favoriteFolder := range favoriteFolders {
			html += s.renderFavoriteItem(favoriteFolder)
		}
		html += "</ul>"
	}

	html += "</li>"
	html += "</ul>"
	return html
}

func (s Service) renderFavoriteItem(folder model.Folder) string {
	li := `<li
		id='jstree-node-%d'
		class='favorite-tree'
		data-url='%s'
		data-idfolder='%d'
		data-jstree='{"icon":"fa fa-folder fa-lg","opened":false, "selected":false}'
>`
	url := s.CoshkeyUrl + "/" + strconv.Itoa(folder.Id)
	html := fmt.Sprintf(li, folder.Id, url, folder.Id)

	span := "<span class='folders-tytle'>%s</span>"
	html += fmt.Sprintf(span, folder.Title)

	if len(folder.Children) > 0 {
		html += "<ul>"
		for _, child := range model.GetSortedSliceFromMap(folder.Children) {
			if child.IsFavorite || child.CountFavoriteKeys > 0 || s.hasFavoriteChild(child) {
				html += s.renderFavoriteItem(child)
			}
		}
		html += "</ul>"
	}

	html += "</li>"
	return html
}

func setCategoryTag(category int) string {
	html := "<ul>"
	switch category {
	case model.PrivateFolder:
		html += "<li class='js__category_title' data-category='Личные папки'></li>"
	case model.PrivateProject:
		html += "<li class='js__category_project' data-category='Мои проекты'></li>"
	case model.SharedProject:
		html += "<li class='js__category_shared_project' data-category='Расшаренные проекты'></li>"
	case model.SharedFolder:
		html += "<li class='js__category_shared_folder' data-category='Расшаренные папки'></li>"
	}

	html += "</ul>"
	return html
}
