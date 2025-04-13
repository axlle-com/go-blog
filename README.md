### Запуск локально
```bash
  go run cmd/main.go
```
### Запуск команды локально 
```bash
  go run cmd/cli/cli.go -command=refill
```
### Запуск тестов
```bash
  go test -count=1 -race ./...
```

### Запуск анализа
```bash
  staticcheck ./...
```