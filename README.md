# Task Scheduler

## Установка

```
git clone https://github.com/ваш-репозиторий/task-scheduler.git
cd task-scheduler
go mod tidy
```

## Настройка
Приложение использует два параметра окружения:

TODO_PORT: Порт, на котором будет работать сервер. По умолчанию используется порт 7540.
TODO_DBFILE: Путь к файлу базы данных SQLite. По умолчанию используется ./scheduler.db.

Запуск приложения осуществляется командой 'go run main.go'

## Структура проекта

- main.go: Главный файл приложения, точка входа сервера.
- task/: Пакет для инициализации базы данных, который также содержит логику обработки задач и вычисления следующей даты выполнения.
- handlers/: Пакет с обработчиками API запросов.
- tests/: В директории находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.
- web/: Директория содержит файлы фронтенда.

### Контейнер
Для запуска приложения в контейнере необходимо выполнить команду:
- создание: sudo docker build --tag task-scheduler:v1 .
- запуск: sudo docker run -p 7540:7540 task-scheduler:v1
### TODO_PASSWORD
Пароль задается командой export TODO_PASSWORD="Ваш пароль" (текущий пароль 2303), снимается unset TODO_PASSWORD
