package server

import (
	"fmt"
	"net/http"
)

func getTreeHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("Запрашиваемый урл: %s\n", req.URL.Path)
	//if req.URL.Path == "/search" {
	//	// разбираем аргументы
	//	args := req.URL.Query()
	//	query := args.Get("query")
	//	limit, err := strconv.Atoi(args.Get("limit"))
	//	if err != nil {
	//		panic("bad limit") // по-хорошему нужно возвращать HTTP 400
	//	}
	//	// выполняем бизнес-логику
	//	results, err := DoBusinessLogicRequest(query, limit)
	//	if err != nil {
	//		resp.WriteHeader(404)
	//		return
	//	}
	//	// устанавливаем заголовки ответа
	//	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	resp.WriteHeader(200)
	//	// сериализуем и записываем тело ответа
	//	json.NewEncoder(resp).Encode(results)
	//}
}

func createRestServer(restAddr string) *http.Server {
	restServer := &http.Server{
		Addr:    restAddr,
		Handler: http.HandlerFunc(getTreeHandler),
		//ReadTimeout:    10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}

	return restServer
}
