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
go test ./pkg/menu/repository -count=1 -v
go test ./pkg/blog/repository -count=1 -v
go test ./pkg/info_block/repository -count=1 -v
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