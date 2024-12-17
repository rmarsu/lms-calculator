# Сервис подсчёта арифметических выражений // Программирование на Go | 24. Спринт 1
Веб сервис для вычисления арифметического выражения. Нужен , если у вас по какой-то причине нет калькулятора , но есть папка с этим проектом и компьютер. 

# Что поддерживает калькулятор
| Действие | Возможная ошибка | Описание |
| --- | --- | --- |
| + | - | Складывает числа. Минимальный приоритет|
| - | - | Вычитает числа. Минимальный приоритет |
| / | деление на 0 | Деление чисел. Приоритет выше чем у + и - |
| * | - | Умножение чисел. Приоритет выше чем у + и - |
| () Группировка | незакрытая скобка | Группирование действий. Повышает приоритет действия |


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
