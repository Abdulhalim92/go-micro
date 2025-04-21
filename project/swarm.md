### Инициализация docker swarm
```bash
docker swarm init --advertise-addr <ip>
```

### Получение токена для добавления ноды в кластер
```bash
docker swarm join-token manager
````

### Расширение количество реплик сервиса
```bash
docker service scale myapp_listener-service=3
```

### Обновление образа сервиса
```bash
docker service update --image <image> <service>
```

### Поднятие сервиса в кластере
```bash
docker stack deploy -c <docker-compose-file> <stack-name>
```

## Поднятие сервиса в кластере с помощью docker-compose
```bash
docker stack deploy -c swarm.yml myapp
```

### Добавление ноды в кластер
```bash
docker swarm join --token <token> <ip>:<port>
```

### Получение токена для добавления ноды
```bash
docker swarm join-token worker
```

### Получение информации о кластере
```bash
docker info
```
### Получение информации о нодах
```bash
docker node ls
```
### Получение информации о сервисах
```bash
docker service ls
```