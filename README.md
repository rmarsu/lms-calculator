# Сервис подсчёта арифметических выражений. Программирование на Go | 24. Спринт 1
REST API для вычисления арифметических выражений.

## Описание
Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Сервис реализует один endpoint, который обрабатывает *POST*-запросы с арифметическими выражениями.

<div align="center">
  ..................
</div>

> [!IMPORTANT]
> ### Сервис поддерживает обработку таких символов как:
> | Символ | Возможная ошибка | Описание |
> | --- | --- | --- |
> | Целое число (int64) | Число превосходит максимальное значение int64 | [Просто целое число](https://ru.wikipedia.org/wiki/%D0%A6%D0%B5%D0%BB%D0%BE%D0%B5_%D1%87%D0%B8%D1%81%D0%BB%D0%BE)|
> | Число с плавающей точкой (float64) | Число превосходит максимальное значение float64 | [Просто число с плавающей точкой](https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%BE_%D1%81_%D0%BF%D0%BB%D0%B0%D0%B2%D0%B0%D1%8E%D1%89%D0%B5%D0%B9_%D0%B7%D0%B0%D0%BF%D1%8F%D1%82%D0%BE%D0%B9)|
> | + | - | Складывает числа. Минимальный приоритет|
> | - | - | Вычитает числа. Минимальный приоритет |
> | / | деление на 0 | Деление чисел. Приоритет выше чем у + и - |
> | * | - | Умножение чисел. Приоритет выше чем у + и - |
> | () Группировка | незакрытая скобка | Группирование действий. Повышает приоритет действия |
<div align="center">
  ..................
</div>

> [!IMPORTANT]
> ### Сервис поддерживает обработку таких символов как:
> | Символ | Возможная ошибка | Описание |
> | --- | --- | --- |
> | Целое число [int64](https://pkg.go.dev/builtin#int64) | Число превосходит максимальное значение [int64](https://pkg.go.dev/builtin#int64) | [Просто целое число](https://ru.wikipedia.org/wiki/%D0%A6%D0%B5%D0%BB%D0%BE%D0%B5_%D1%87%D0%B8%D1%81%D0%BB%D0%BE)|
> | Число с плавающей точкой (float64) | Число превосходит максимальное значение float64 | [Просто число с плавающей точкой](https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%BE_%D1%81_%D0%BF%D0%BB%D0%B0%D0%B2%D0%B0%D1%8E%D1%89%D0%B5%D0%B9_%D0%B7%D0%B0%D0%BF%D1%8F%D1%82%D0%BE%D0%B9)|
> | + | - | Складывает числа. Минимальный приоритет|
> | - | - | Вычитает числа. Минимальный приоритет |
> | / | деление на 0 | Деление чисел. Приоритет выше чем у + и - |
> | * | - | Умножение чисел. Приоритет выше чем у + и - |
> | () Группировка | незакрытая скобка | Группирование действий. Повышает приоритет действия |


# Для запуска
Для запуска можете воспользоваться Makefile-ом
```bash
$ make run
```
или запустить вручную
```bash
$ go run cmd/main.go
```
> [!IMPORTANT]
> Убедитесь что установлена версия Go 1.23.3

# Примеры использования с cURL:
cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Ожидаемый ответ сервера:
```json
// код http 200
{"result":6}
```

cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2 -"
}'
```
Ожидаемый ответ сервера:
```json
// код http 422
{"error":"Expression is not valid"}
```

cURL:
```bash
curl --request GET \
--url "http://localhost:8080/api/v1/calculate" \
--header "Content-Type: application/json" \
--data '{"expression":"1+1"}'
```
Ожидаемый ответ сервера:
```json
// код http 405
{"error":"Only POST method is allowed"}
```

cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "meow": "2 - 1"
'
```
Ожидаемый ответ сервера:
```json
// код http 400
{"error":"Bad request"}
```





![images](https://github.com/user-attachments/assets/09b0393e-ed77-4a61-8ea4-8b057ffb07c1)



# Для запуска
Для запуска можете воспользоваться Makefile-ом
```bash
$ make run
```
или запустить вручную
```bash
$ go run cmd/main.go
```
> [!IMPORTANT]
> Убедитесь что установлена версия Go 1.23.3

# Примеры использования с cURL:
cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Ожидаемый ответ сервера:
```json
// код http 200
{"result":6}
```

cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2 -"
}'
```
Ожидаемый ответ сервера:
```json
// код http 422
{"error":"Expression is not valid"}
```

cURL:
```bash
curl --request GET \
--url "http://localhost:8080/api/v1/calculate" \
--header "Content-Type: application/json" \
--data '{"expression":"1+1"}'
```
Ожидаемый ответ сервера:
```json
// код http 405
{"error":"Only POST method is allowed"}
```

cURL:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "meow": "2 - 1"
'
```
Ожидаемый ответ сервера:
```json
// код http 400
{"error":"Bad request"}
```

