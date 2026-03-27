# Port Scanner

Консольная утилита для сканирования TCP/UDP портов на удалённом хосте. Написана на Go,
сканирует порты параллельно через пул воркеров, поэтому даже полный диапазон 1-65535
проверяется за удивительно небольшое время.

## Установка и запуск

Нужен установленный Go (проверялось на 1.25).

```bash
git clone https://github.com/RishatShay/Port-Scanner.git
cd Port-Scanner
go build -o bin/portscanner ./cmd/portscanner
./bin/portscanner -host scanme.nmap.org -ports 1-1024
```

Либо без сборки бинарника:

```bash
go run ./cmd/portscanner -host scanme.nmap.org -ports 1-1024
```

Если указать в качестве хоста localhost, рядом с портом покажется процесс, который
его занимает:

```bash
go run ./cmd/portscanner -host localhost -ports 8080,8081
```

## Флаги

| Флаг        | По умолчанию         | Описание                                          |
|-------------|----------------------|----------------------------------------------------|
| `-host`     | `scanme.nmap.org`    | хост, который сканируем                            |
| `-protocol` | `tcp`                | протокол: `tcp` или `udp`                          |
| `-ports`    | `1-65535`            | список/диапазоны портов, например `22,80,8000-8100`|
| `-workers`  | `500`                | количество горутин-воркеров                        |
| `-timeout`  | `1s`                 | таймаут подключения на один порт                   |

## Пример

```bash
$ go run ./cmd/portscanner -host scanme.nmap.org -ports 20-25,80,443
22/tcp is open
80/tcp is open
scanned 8 ports in 431.494137ms, found 2 open
```

Если порт закрыт или недоступен, то не попадёт в вывод. При сканировании localhost
рядом с портом также покажется, кто его занимает:

```bash
$ go run ./cmd/portscanner -host localhost -ports 8080,8081
8080/tcp is open (pid 89034, node)
8081/tcp is open (pid 88842, reco-bills)
scanned 2 ports in 12.3ms, found 2 open
```

Сканер понимает Ctrl+C: по SIGINT/SIGTERM он завершает уже запущенные проверки и
выходит, не зависая на недосканированных портах.

## Структура проекта

Проект собран по стандартному для Go расположению пакетов:

```
cmd/portscanner/   - точка входа, разбор флагов, вывод результата
internal/scanner/  - вся логика сканирования, наружу не импортируется
internal/owner/    - определение процесса-владельца порта на локальной машине
```

Внутри `internal/scanner`:

- `scanner.go` - проверка одного порта (`ScanPort`)
- `pool.go` - пул воркеров, который разгоняет проверки портов по горутинам (`Run`)
- `ports.go` - разбор строки с портами вида `22,80,8000-8100` (`ParsePorts`)

Внутри `internal/owner`:

- `owner.go` - общая часть (`Lookup`, `IsLocalHost`)
- `lookup_linux.go` / `lookup_windows.go` / `lookup_fallback.go` - платформенная реализация `Lookup`

## Тесты

```bash
go test ./...
```

Тесты на `ScanPort` поднимают локальный TCP-listener на `127.0.0.1`, чтобы не
зависеть от сети и внешних хостов. Тесты на `ParsePorts` - обычные табличные тесты
на корректные и некорректные строки с портами.
