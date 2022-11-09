package server

type TreeRequestBody struct {
	SelectedFolderId int  `schema:"selected_folder_id"`
	UserId           int  `schema:"user_id"`
	OwnTree          bool `schema:"own_tree"`
}
