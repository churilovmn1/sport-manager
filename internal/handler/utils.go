package handler

import (
	"encoding/json"
	"net/http"
	// Если у вас есть система логирования, подключите ее здесь,
	// чтобы записывать внутренние ошибки кодирования
)

// writeJSONResponse - общая функция для записи JSON-ответа.
// Она устанавливает заголовок Content-Type: application/json.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// В случае, если кодирование JSON не удалось,
			// мы ничего не можем вернуть клиенту,
			// так как заголовки уже отправлены.
			// Здесь нужно только логировать ошибку на стороне сервера.
			// fmt.Printf("Internal error encoding response: %v\n", err)
			return
		}
	}
}

// writeErrorResponse - вспомогательная функция для записи ошибок в JSON-формате.
// Возвращает ответ в виде {"error": message}.
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	// Используем исправленную функцию writeJSONResponse
	writeJSONResponse(w, statusCode, map[string]string{"error": message})
}
