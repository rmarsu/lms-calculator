# Сервис подсчёта арифметических выражений. Программирование на Go | 24. Спринт 1
REST API для вычисления арифметических выражений.

## Описание
Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления. Сервис реализует один endpoint, который обрабатывает *POST*-запросы с арифметическими выражениями.

В проекте присуствует:
- логирование
- тесты (вы можете их [запустить](internal/transport/http/calc_test.go) для проверки API)
- конфиг
- мейкфайл
  
<div align="center">
  ..................
</div>

> [!IMPORTANT]
> ### Сервис поддерживает обработку таких символов как:
> | Символ | Возможная ошибка | Описание |
> | --- | --- | --- |
> | Целое число ([int64](https://pkg.go.dev/builtin#int64)) | Число превосходит максимальное значение [int64](https://pkg.go.dev/builtin#int64) | [Просто целое число](https://ru.wikipedia.org/wiki/%D0%A6%D0%B5%D0%BB%D0%BE%D0%B5_%D1%87%D0%B8%D1%81%D0%BB%D0%BE)|
> | Число с плавающей точкой ([float64](https://pkg.go.dev/builtin#float64))| Число превосходит максимальное значение [float64](https://pkg.go.dev/builtin#float64) | [Просто число с плавающей точкой](https://ru.wikipedia.org/wiki/%D0%A7%D0%B8%D1%81%D0%BB%D0%BE_%D1%81_%D0%BF%D0%BB%D0%B0%D0%B2%D0%B0%D1%8E%D1%89%D0%B5%D0%B9_%D0%B7%D0%B0%D0%BF%D1%8F%D1%82%D0%BE%D0%B9)|
> | + | - | Складывает числа. Минимальный приоритет|
> | - | - | Вычитает числа. Минимальный приоритет |
> | / | Деление на 0 | Деление чисел. Приоритет выше чем у + и - |
> | * | - | Умножение чисел. Приоритет выше чем у + и - |
> | () Группировка | Незакрытая скобка | Группирование действий. Повышает приоритет действия |
<div align="center">
  ..................
</div>

> [!IMPORTANT]
> ### Эндпоинты
> | Эндпоинт | Допустимые методы | Описание |
> | --- | --- | --- |
> | /api/v1/calculate | *POST* | Получает POST-запрос c телом запроса в формате [JSON](https://ru.wikipedia.org/wiki/JSON), содержащим выражение. Отдает результат или ошибку в формате [JSON](https://ru.wikipedia.org/wiki/JSON) |



## Для запуска
> [!IMPORTANT]
> Если случилась ошибка ,  <ins>убедитесь что установлена версия Go `1.23.3`</ins>.
> Последнюю версию можно установить [здесь](https://go.dev/dl/).

Для запуска можете воспользоваться Makefile-ом
```shell
$ make run
```
или запустить вручную:
```shell
$ go run cmd/main.go
```

> [!TIP]
> Для смены порта сервера стоит воспользоваться [конфигом](configs/config.yaml). 

## Примеры использования с cURL:

| cURL команда                                   | Ответ                                     | *HTTP* код
|------------------------------------------------|-------------------------------------------| ----------------------------- |
| ```curl --location 'localhost:8080/api/v1/calculate' \ --header 'content-type: application/json' \ --data '{ "expression": "2+2*2" }'```  | ```{"result":6} ``` | 200 |
| ```curl --location 'localhost:8080/api/v1/calculate' \ --header 'content-type: application/json' \ --data '{ "expression": "2 -" }'``` | ```{"error":"это не выражение"}```|422|
| ```curl --request GET \ --url "http://localhost:8080/api/v1/calculate" \ --header "Content-Type: application/json" \ --data '{"expression":"1+1"}'``` | ```{"error":"Only POST method is allowed"}```|405|
| ```curl --location 'localhost:8080/api/v1/calculate' \ --header 'Content-Type: application/json' \ --data '{ "meow": "2 - 1" '``` | ```{"error":"Bad request"}```|400|

> [!CAUTION]
> При использовании powershell или cmd могут возникуть проблемы с работой cURL , так как в них нельзя использовать одинарные кавычки. Можете воспользоваться аналогами , такими как : [Postman](https://www.postman.com/) или [WSL](https://en.wikipedia.org/wiki/Windows_Subsystem_for_Linux) . 

## Лицензия
[MIT](LICENSE)




