package server

type TreeRequestBody struct {
	SelectedFolderId int  `json:"selected_folder_id"`
	UserId           int  `json:"user_id"`
	OwnTree          bool `json:"own_tree"`
}
