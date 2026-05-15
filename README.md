# YaDro-Impulse

клиент-серверное приложение для управления DNS-серверами

В терминале
go build -o dns-server ./cmd/server
go build -o dns-client ./cmd/client
sudo ./dns-server

далее через cli: ./dns-client общаться с сервером

## cmd

### `server/main.go` - точка входа серверного приложения. Здесь происходит загрузка из .env,
инициализация структур из internal

### `cmd/client/main.go` - точка входа cli клиента. Поддерживает:
--help
add <dns>
del <dns>
list

## internal

### `handlers.go` - http слой с эндпоинтами:
POST /dns
GET /dns
DELETE /dns

этот слой отвечает за чтение json, вызов бизнес логики, отправка ответов

### `usecases.go` - слой бизнес логики, чтобы не перегружать слой с http зпросами.

здесь происходит только обработка dns и проверка на корректность. Добавляет в файл, обращаясь к repository

Add (newDns string)
Delete (dnsToDelete string)
List ()

### `repository.go` - слой данных, отвечает за добавление dns серверов в файл в соответствии с форматом файла

и игнорирование иных строк (комментариев)

### `errors.go` - файл с общими ошибками приложения

## etc

### `resolv.conf` - тестовый файл, был создан для тестирования

## .env-example 
следует переименовать в .env чтобы подгружать зависимости


