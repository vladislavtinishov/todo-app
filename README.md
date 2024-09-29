Рекомендованные утилиты для работы с приложением:
- migrate

Установка migrate на Windows:
- irm get.scoop.sh | iex
- scoop install migrate

Создание файлов миграций в migrate:
- migrate create -ext sql -dir /schema -seq {имя миграций}

Запуск миграций
- migrate -path ./schema -database 'mysql://{username}:{password}@tcp({host}:{port})/{db}' up