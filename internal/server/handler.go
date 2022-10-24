package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	internalErrors "github.com/ITarako/coshkey_tree/internal/pkg/errors"
	"github.com/ITarako/coshkey_tree/internal/service/tree"
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

	//router.HandleFunc("/tree", handler.getTree).Methods(http.MethodGet).Queries("user_id", "{user_id:[0-9]+}")
	router.HandleFunc("/tree", handler.getTree).Methods(http.MethodGet)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) getTree(w http.ResponseWriter, r *http.Request) {
	//args := r.URL.Query()
	//userId, _ := strconv.Atoi(args.Get("user_id"))
	requestBody := new(TreeRequestBody)

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		writeError(w, nil, http.StatusBadRequest, err, "json decode")
		return
	}

	user, err := h.treeService.UserService.GetUser(r.Context(), requestBody.UserId)
	if err != nil {
		if errors.Is(err, internalErrors.ErrNotFound) {
			writeError(w, nil, http.StatusNotFound, err, "user not found")
			return
		}

		writeError(w, nil, http.StatusBadRequest, err, "UserService.GetUser")
		return
	}

	userIsAdmin, err := h.treeService.RbacService.CheckRole(r.Context(), "admin", user.Id)
	if err != nil {
		writeError(w, nil, http.StatusBadRequest, err, "RbacService.CheckRole")
		return
	}
	user.IsAdmin = userIsAdmin

	res := h.treeService.Generate(user, requestBody.SelectedFolderId, requestBody.OwnTree)

	//_, _ = fmt.Fprintf(w, "requestBody: %+v\n", requestBody)
	writeSuccess(w, res)
}
