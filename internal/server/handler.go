package server

import (
	"net/http"

	"github.com/ITarako/coshkey_tree/internal/pkg/errors"
	"github.com/ITarako/coshkey_tree/internal/server/helper"
	"github.com/ITarako/coshkey_tree/internal/service/tree"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

type Handler struct {
	router      *mux.Router
	treeService tree.Service
}

func NewHandler(treeService tree.Service) *Handler {
	router := mux.NewRouter()
	handler := &Handler{
		router:      router,
		treeService: treeService,
	}

	router.HandleFunc("/tree", handler.getTree).Methods(http.MethodGet)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) getTree(w http.ResponseWriter, r *http.Request) {
	requestData := new(TreeRequestBody)

	if err := schema.NewDecoder().Decode(requestData, r.URL.Query()); err != nil {
		serverhelper.WriteError(w, nil, http.StatusBadRequest, err, "json decode")
		return
	}

	user, err := h.treeService.UserService.GetUser(r.Context(), requestData.UserId)
	if err != nil {
		if errors.Is(err, internalerrors.ErrNotFound) {
			serverhelper.WriteError(w, nil, http.StatusNotFound, err, "user not found")
			return
		}

		serverhelper.WriteError(w, nil, http.StatusBadRequest, err, "UserService.GetUser")
		return
	}

	userIsAdmin, err := h.treeService.RbacService.CheckRole(r.Context(), "admin", user.Id)
	if err != nil {
		serverhelper.WriteError(w, nil, http.StatusBadRequest, err, "RbacService.CheckRole")
		return
	}
	user.IsAdmin = userIsAdmin

	res := h.treeService.Generate(r.Context(), user, requestData.SelectedFolderId, requestData.OwnTree)

	serverhelper.WriteSuccess(w, res)
}
