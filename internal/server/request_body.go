package server

type TreeRequestBody struct {
	SelectedFolderId int32 `json:"selected_folder_id"`
	UserId           int32 `json:"user_id"`
	OwnTree          bool  `json:"own_tree"`
}
