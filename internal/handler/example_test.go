package handler

func ExampleServerHandler_Ping() {
	// Запрос:
	//GET /ping

	// Output:
	//Код состояния 200 OK - если подключение к базе данных активно и работает
	//Код состояния 500 Internal Server Error - если подключение к базе данных отсутствует или не работает
}

func ExampleServerHandler_Main() {
	// Запрос:
	// GET /

	// Output:
	//Всегда возвращает HTML-страницу с кодом 200 OK
}

func ExampleServerHandler_Value() {
	// Запрос:
	// GET /value/{type}/{name}, где:
	// {type} - тип метрики (gauge или counter)
	// {name} - имя метрики (например, Alloc, PollCount)

	// Output:
	//Код состояния 200 OK при успешном выполнении, возвращает значение метрики в виде простого текста
	//Код состояния 404 Not Found если метрика не найдена
	//Код состояния 400 Bad Request если указан неверный тип метрики
}

func ExampleServerHandler_ValueJSON() {
	// Запрос:
	// POST /value/ HTTP/1.1

	//Content-Type: application/json
	//{
	//  "id": "{name}",
	//  "type": "{type} "
	//}
	// {type} - тип метрики (gauge или counter)
	// {name} - имя метрики (например, Alloc, PollCount)

	// Output:
	//Код состояния 200 OK при успешном выполнении, возвращает значение метрики
	//	Content-Type: application/json {"id":"PollCount","type":"counter","delta":42}
	//Код состояния 404 Not Found если метрика не найдена
	//Код состояния 400 Bad Request если указан неверный тип метрики
}

func ExampleServerHandler_UpdateBatchJSON() {
	// Запрос:
	// POST /updates/
	// Content-Type: application/json
	//
	// [
	//   {
	//     "id": "{name}",
	//     "type": "{type}",
	//     "value": float64
	//   },
	//   {
	//     "id": "{name}",
	//     "type": "{type}",
	//     "delta": int64
	//   }
	// ]
	// {type} - тип метрики (gauge или counter)
	// {name} - имя метрики (например, Alloc, PollCount)
	// {value} - значение метрики (int64 для counter или float64 для gauge)

	// Output:
	// Код состояния 200 OK - при успешном обновлении всех метрик
	// Код состояния 400 Bad Request - при некорректных данных в запросе
	//   - Некорректный JSON
	//   - Неверный тип метрики
	//   - Некорректное значение метрики
}

func ExampleServerHandler_Update() {
	// Запрос:
	// POST /update/{type}/{name}/{value}
	// {type} - тип метрики (gauge или counter)
	// {name} - имя метрики (например, Alloc, PollCount)
	// {value} - значение метрики (int64 для counter или float64 для gauge)

	// Output:
	// Код состояния 200 OK - при успешном обновлении метрики
	// Код состояния 400 Bad Request - при некорректных данных в запросе
	//   - Неверный тип метрики
	//   - Некорректное значение метрики
	// Код состояния 405 Method Not Allowed - при использовании метода отличного от POST
}
