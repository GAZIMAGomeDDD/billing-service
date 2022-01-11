# Billing Service


## Запуск

1. Вводим команду для создания docker образов:
```bash
docker-compose build
```
2. Как только образы будут собраны, запускаем контейнеры командой::
```bash
docker-compose up -d
```
3. Смотрим логи:
```bash
docker-compose logs -f
```
4. Для приостановки docker контейнеров используйте команду:
```bash
docker-compose down
```

### Примеры запросов
Изменение(увеличение или уменьшение) баланса по uuid пользователя(По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). 
Данные о балансе появляются при первом зачислении денег.). 
```
$ curl --location --request POST 'localhost:8080/changeBalance' \
    --header 'Content-Type: application/json' \
    --data-raw '{"id": "b91a95a4-078f-4afd-b11c-4850eb65e782", "money": 2300.212}'

```

Узнать баланс по uuid пользователя (базовая валюта которая хранится на балансе у пользователя всегда рубль, чтоб ухнать баланс в других валютах, то нужно добавить к методу получения баланса доп. параметр: currency. Пример: ?currency=USD.):
```
$ curl --location --request POST 'localhost:8080/getBalance?currency=USD' \
    --header 'Content-Type: application/json' \
    --data-raw '{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"}'
```

Перевод средств от пользователя к пользователю:
```
$ curl --location --request POST 'localhost:8080/moneyTransfer' \
    --header 'Content-Type: application/json' \
    --data-raw '{"to_id": "b91a95a4-078f-4afd-b11c-4850eb65e782", "from_id": "b91a95a4-078f-4afd-b11c-4850eb65e784", "money": 2564.52}'
```

Получение списка транзакций по uuid пользователя. Предусмотрена пагинация и сортировка по сумме и дате. Для сортировки по дате и сумме надо использовать параметр sort (date_asc, money_asc, date_desc, money_desc)
```
$ curl --location --request POST 'localhost:8080/listOfTransactions?page=1&sort=money_desc' \
    --header 'Content-Type: application/json' \
    --data-raw '{"user_id": "b91a95a4-078f-4afd-b11c-4850eb65e784", "limit": 7}'
```

### Документация доступна по адресу [http://localhost:8080/swagger/]


## Tests
Запустить юнит тесты:
```bash
    go test ./... -cover
```
Запустить интнграционные тесты:
```bash
    go test -tags=integration ./...
```
