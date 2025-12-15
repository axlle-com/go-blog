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

### deploy
```bash

chmod +x deploy.sh
```


### 
```bash
sudo fuser -k 3000/tcp

```
```bash
sudo ss -ltnp | grep ':3000'
sudo kill -9 12345
```