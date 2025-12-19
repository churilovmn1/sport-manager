package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSONResponse — универсальный вспомогательный метод для отправки JSON-ответов.
// Централизация этого процесса позволяет гарантировать наличие правильных заголовков во всем API.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Устанавливаем заголовок типа контента перед записью статус-кода
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		// Кодируем структуру в JSON и пишем напрямую в поток ответа
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Ошибка на этом этапе означает проблему с данными внутри сервера (например, циклическая ссылка).
			// Так как заголовки уже ушли, мы можем только зафиксировать инцидент в логах.
			log.Printf("CRITICAL: Failed to encode JSON response: %v", err)
			return
		}
	}
}

// writeErrorResponse — обертка для стандартизации сообщений об ошибках.
// Помогает фронтенду всегда ожидать объект вида {"error": "описание"}.
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	// Мы инкапсулируем создание карты (map) здесь, чтобы не дублировать код в хендлерах
	writeJSONResponse(w, statusCode, map[string]string{
		"error": message,
	})
}
