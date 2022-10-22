package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ITarako/coshkey_tree/internal/service/tree"
)

type Handler struct {
	router  *mux.Router
	service tree.Service
}

func NewHandler(treeService tree.Service) *Handler {
	router := mux.NewRouter()
	handler := &Handler{
		router:  router,
		service: treeService,
	}
	router.HandleFunc("/tree", handler.getTree).Methods(http.MethodGet).Queries("id", "{id:[0-9]+}")
	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) getTree(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	id, err := strconv.Atoi(args.Get("id"))
	if err != nil {
		log.Error().Err(err).Msg("error strconv")
		return
	}

	folder, err := h.service.GetFolder(r.Context(), int32(id))
	if err != nil {
		log.Error().Err(err).Msg("error get folder")
		return
	}

	_, _ = fmt.Fprintf(w, "Folder: %+v", folder)

	//requestBody := createBookRequestBody{}
	//if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
	//	writeError(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//book, err := h.service.CreateBook(requestBody.Title)
	//if err != nil {
	//	writeError(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//writeSuccess(w, book)
}

func getTree(resp http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(resp, "Запрашиваемый урл: %s\n", req.URL.Path)
	if req.URL.Path == "/tree" {
		_, _ = fmt.Fprintln(resp, "Начинаем строить дерево")
		// разбираем аргументы
		//args := req.URL.Query()
		//query := args.Get("query")
		//limit, err := strconv.Atoi(args.Get("limit"))
		//if err != nil {
		//	panic("bad limit") // по-хорошему нужно возвращать HTTP 400
		//}
		//// выполняем бизнес-логику
		//results, err := DoBusinessLogicRequest(query, limit)
		//if err != nil {
		//	resp.WriteHeader(404)
		//	return
		//}
		//// устанавливаем заголовки ответа
		//resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		//resp.WriteHeader(200)
		//// сериализуем и записываем тело ответа
		//json.NewEncoder(resp).Encode(results)
	}
}
