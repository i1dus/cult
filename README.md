# cult

Back-End часть парковочной системы.

## Инструкция запуска


Пре реквизиты:
- Go, не ниже версии 1.23.
- Docker (Postgresql).

Последовательность:
1. Поднять БД: ```make db-up```
2. Запустить миграции: ```make up```
3. Запустить приложение: ```make run```
4. Наслаждаться жизнью охранником парковки...

Приложение работает на хосте `8080`.

## Swagger

Подключение к Swagger'у осуществляется по URL:
`http://localhost:8080/swagger/index.html`





