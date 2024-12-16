# Calculator Web Service

Этот проект предоставляет веб-сервис для вычисления арифметических выражений, отправленных пользователем через HTTP. 

## Описание

Пользователь отправляет POST-запрос с арифметическим выражением на URL `/api/v1/calculate`, и в ответ получает результат вычисления. 

Программа поддерживает следующие возможности:
- Арифметические операции: `+`, `-`, `*`, `/`
- Целые и дробные числа
- Обработка приоритета операций с классическими арифметическими знаками, а также скобками ("(" и ")")
- Возвращает результат в случае успешного вычисления
- Обрабатывает ошибки с понятными HTTP-ответами

## Эндпоинт

### POST /api/v1/calculate

#### Запрос
Тело запроса должно содержать JSON с ключом `expression`, содержащим строку с арифметическим выражением:
```json
{
    "expression": "2+2*2"
}
```
#### Ответы
1. **Успешное вычисление (200 OK):**
```json
{
    "result": "6"
}
```

2. **Некорректное выражение (422 Unprocessable Entity):**
```json
{
    "error": "Expression is not valid"
}
```

3. **Внутренняя ошибка сервера (500 Internal Server Error):**
```json
{
    "error": "Internal server error"
}
```

## Примеры использования

### Пример cURL запроса
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Ответ:
```json
{
    "result": "6"
}
```

### Некорректное выражение
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*"
}'
```
Ответ:
```json
{
    "error": "Expression is not valid"
}
```

### Внутренняя ошибка сервера
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": ""
}'
```
Ответ:
```json
{
    "error": "Internal server error"
}
```

## Установка и запуск

1. Убедитесь, что у вас установлен Go версии 1.20 или выше.
2. Склонируйте репозиторий:
    ```bash
    git clone https://github.com/ваш-логин/calculator-service.git
    cd calculator-service
    ```
3. Запустите сервис:
    ```bash
    go run ./cmd/main.go
    ```

Сервис будет запущен по адресу `http://localhost` на порте, указанном в переменной среды "PORT". Значение порта по умолчанию - 1081.