# Простой планировщик задач (TODO List)

## Описание проекта:

Этот проект представляет собой простейший планировщик задач (TODO List), позволяющий пользователям добавлять, редактировать, перемещать и удалять задачи, а также просматривать все текущие задачи на одной странице. Проект создан в качестве учебного упражнения для ознакомления с Go, а также с различными аспектами веб-разработки.

## Стек технологий:

• Go 

• SQLite

• HTML, CSS, JavaScript

• Docker

## Инструкция по развёртыванию:

Системные требования:

• Go 1.18+

• Docker

## Шаги по запуску:

1. Установка зависимостей: go mod tidy
  
2. Создание образа Docker:  docker build -t task-app .
  
3. Запуск образа Docker:  docker run -p 7540:7540 task-app
  
4. Доступ к приложению:  Откройте веб-браузер и перейдите по адресу: http://localhost:7540/

## Запуск тестов:
  
  go test ./tests
  
## Планы по доработке проекта:

• Добавить возможность сортировки задач по различным критериям (дата, приоритет, статус).

• Реализовать аутентификацию и авторизацию пользователей.

• Внедрить систему уведомлений о предстоящих задачах.

• Добавить интеграцию с внешними сервисами (например, Google Calendar).
