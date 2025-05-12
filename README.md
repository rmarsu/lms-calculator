# Сервис подсчёта арифметических выражений. Программирование на Go | 24. Спринт 1
REST API для вычисления арифметических выражений.
> [!IMPORTANT]
>  Если у вас возникли проблемы - не ставьте мне ноль баллов , а лучше напишите сюда в тг -> @rmarsu

## Описание
Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Back-end часть состоит из 3 микросервисов:
Сервер (оркестратор) , который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
Вычислитель (агент), который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Авторизатор - регистрирует и логинит пользователей, отдавая jwt-токен. Общение между агентом и оркестратором происходит через grpc.

<div align="center">
  ..................
</div>

### Сервис поддерживает обработку таких символов как:
| Символ | Возможная ошибка | Описание |
| --- | --- | --- |
| Целое число ([int64](https://pkg.go.dev/builtin#int64)) | Число превосходит максимальное значение [int64](https://pkg.go.dev/builtin#int64) | [Просто целое число](https://ru.wikipedia.org/wiki/%D0%A6%D0%B5%D0%BB%D0%BE%D0%B5_%D1%87%D0%B8%D1%81%D0%BB%D0%BE)|
| Число с плавающей точкой ([float64](https://pkg.go.dev/builtin#float64))| Число превосходит максимальное значение [float64](https://pkg.go.dev/builtin#float64) | [Просто число с плавающей точкой](https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%BE_%D1%81_%D0%BF%D0%BB%D0%B0%D0%B2%D0%B0%D1%8E%D1%89%D0%B5%D0%B9_%D0%B7%D0%B0%D0%BF%D1%8F%D1%82%D0%BE%D0%B9)|
| + | - | Складывает числа. Минимальный приоритет|
| - | - | Вычитает числа. Минимальный приоритет |
| / | Деление на 0 | Деление чисел. Приоритет выше чем у + и - |
| * | - | Умножение чисел. Приоритет выше чем у + и - |
| () Группировка | Незакрытая скобка | Группирование действий. Повышает приоритет действия |

<div align="center">
  ..................
</div>

### Эндпоинты
| Эндпоинт | Допустимые методы |Нужна аутентификация?| Описание |
| --- | --- | --- | --- |
| /api/v1/calculate | *POST* |✅| Получает POST-запрос c телом запроса в формате [JSON](https://ru.wikipedia.org/wiki/JSON), содержащим выражение. |
| /api/v1/expressions | *GET* |❌| Получает все выражения |
| /api/v1/expressions/:id | *GET* |❌|Получает выражение по его UUID |
| /internal/task | *GET* |❌| Получает простое арифметическое выражение как "задачу" |
| /internal/task | *POST* |❌| Принимает результат задачи |
| /api/v1/register  | *POST* | ❌| Регистрирует пользователя в системе |
| /api/v1/login  | *POST* | ❌| Отдает jwt-токен |



## Для запуска
> [!TIP]
> Если случилась ошибка ,  <ins>убедитесь что установлена версия Go `1.24`</ins>.
> Последнюю версию можно установить [здесь](https://go.dev/dl/).

> ДЛЯ SQLITE поставьте gcc или mingw

Не забудьте установить зависимости командой:
```shell
$ go mod tidy
```
> [!IMPORTANT]
> Перед запуском создайте файл .env в корне проекта. Он должен содержать такие данные как:
> ```
> AUTH_GRPC_PORT=":50051"
> AUTH_REST_PORT=":8080"
> AUTH_SQLITE_PATH="./db/auth_db/auth.sqlite"
>
> ORCHESTRATOR_GRPC_PORT=":50052"
> ORCHESTRATOR_REST_PORT=":8081"
> ORCHESTRATOR_SQLITE_PATH="./db/orch_db/orchestrator.sqlite"
>
> JWT_SECRET="supersecretjwtkey"
> HASHER_SALT="randomsaltvalue123"
>
> COMPUTING_POWER=10
> HOST="localhost"
>
> TIME_ADDITION_MS=1000
> TIME_SUBTRACTION_MS=500
> TIME_MULTIPLICATION_MS=2
> TIME_DIVISION_MS=4
> ```

Для запуска запустите все 3 сервиса в разных терминалах
```shell
$ go run agent_service/cmd/main.go
$ go run auth_service/cmd/main.go
$ go run orchestrator_service/cmd/main.go
```
или запустите через docker-compose:
```
$ docker-compose up --build
```
!! При использовании docker-compose доступен nginx, работающий на localhost:80

## Примеры использования с cURL:

| cURL команда                                   | Ответ                                     |
|------------------------------------------------|-------------------------------------------|
| ```curl -XPOST -H 'Authorization: Bearer YOUR-TOKEN' -d '{ "expr" : "2+2*2"}' 'http://localhost:8081/api/v1/calculate'```  | ```{"id":"4", "status":"StatusPending", "result":0} ```|
| ```curl -XPOST -H 'Authorization: Bearer YOUR-TOKEN' -d '{ "expr" : "2+2@2"}' 'http://localhost:8081/api/v1/calculate'``` | ```{"code":2, "message":"invalid expression", "details":[]}```|
| ```curl -XPOST -H 'Authorization: Bearer YOUR-TOKEN' -d '{ "expr" : "2+2+2",}' 'http://localhost:8081/api/v1/calculate'``` | ```{    "code": 3,    "message": "invalid character '}' looking for beginning of object key string",    "details": []}```|
| ```curl --location 'localhost:8081/api/v1/expressions'``` | ```{ "expressions": [ { "id": <идентификатор выражения>, "status": <статус вычисления выражения>, "result": <результат выражения> }, { "id": <идентификатор выражения>, "status": <статус вычисления выражения>, "result": <результат выражения> } ] }``` |
| ```curl --location 'localhost:8081/api/v1/expressions/2'``` | ```{ "expression": { "id": <идентификатор выражения>, "status": <статус вычисления выражения>, "result": <результат выражения> } }```| 
| ```curl --location 'localhost:8081/api/v1/expressions/4184237'``` | ```{"code":2, "message":"expression not found", "details":[]}```|  
|```curl -XPOST -d '{ "username" : "silly_username", "password" : "12345678"}' 'http://localhost:8080/api/v1/register'```| ```{}``` |
|```curl -XPOST -d '{ "username" : "silly_username", "password" : "12345678"}' 'http://localhost:8080/api/v1/register'```| ```{"code":6,"message":"user already exists","details":[]}``` |
|```curl -XPOST -d '{ "username" : "silly_username", "password" : "123456"}' 'http://localhost:8080/api/v1/register'```| ```{"code":10,"message":"password is too short","details":[]}``` |
|```curl -XPOST -d '{ "username" : "silly_username", "password" : "12345678"}' 'http://localhost:8080/api/v1/login'```| ```{"token":"___"}``` | 200 |
|```curl -XPOST -d '{ "username" : "silly_username", "password" : "12345578"}' 'http://localhost:8080/api/v1/login'```| ```{"code":16,"message":"password is incorrect","details":[]}``` |
|```curl -XPOST -d '{ "username" : "silly_user", "password" : "12345578"}' 'http://localhost:8080/api/v1/login'```| ```{"code":5,"message":"user with this username not found","details":[]}``` |





> [!CAUTION]
> При использовании powershell или cmd могут возникуть проблемы с работой cURL , так как в них нельзя использовать одинарные кавычки. Можете воспользоваться аналогами , такими как : [Postman](https://www.postman.com/) или [WSL](https://en.wikipedia.org/wiki/Windows_Subsystem_for_Linux) . 

## Лицензия
[MIT](LICENSE)




