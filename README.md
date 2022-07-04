# COMMENTS SERVICE

* **GET /comment/{id}** - добавляет комментарий к новости по id новости.
* **POST /comment/{id}** - добавляет все комментари к новости по id новости.

### Микросервис работает в связке с другими сервисами:
* gateway: https://github.com/MarySmirnova/api_gateway
* сервис модерации комментариев: https://github.com/MarySmirnova/moderation_service
* сервис - парсер новостей: https://github.com/MarySmirnova/news_reader

### .env example:

    API_COMM_LISTEN=:8080
	API_COMM_READ_TIMEOUT=30s
	API_COMM_WRITE_TIMEOUT=30s

    PG_COMM_USER=
	PG_COMM_PASSWORD=
	PG_COMM_HOST=
	PG_COMM_PORT=
	PG_COMM_DATABASE=
	PG_COMM_TEST_DATABASE=