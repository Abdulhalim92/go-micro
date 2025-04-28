### Создание кластера
```bash
minikube start --nodes=2 --cpus=4 --memory=4096
```

### Проверка статуса кластера
```bash
minikube status
```

### Остановка kubectl
```bash
minikube stop
```

### Просмотр 
```bash
minikube dashboard
```

### Вывод списка Pod - ов кластера Kubernetes для текущего контекста
```bash
kubectl get pods
```

### Вывод списка Pod - ов кластера Kubernetes для текущего контекста с дополнительной информацией
```bash
kubectl get pods -o wide
```
### Вывод списка Pod - ов кластера Kubernetes для текущего контекста с дополнительной информацией и фильтрацией по имени
```bash
kubectl get pods -o wide | grep <имя>
```
### Вывод списка Pod - ов кластера Kubernetes для всех контекстов
```bash
kubectl get pods --all-namespaces
```
### Вывод списка Pod - ов кластера Kubernetes для всех контекстов с дополнительной информацией
```bash
kubectl get pods --all-namespaces -o wide
```

### Вывод сервисов
```bash
kubectl get services
```

### Вывод сервисов с дополнительной информацией
```bash
kubectl get services -o wide
```

### Запуск сервиса с помощью манифеста
```bash
kubectl expose deployment broker-service --type=LoadBalancer --port=8080 --target-port=8080
```

### Добавление ingress
```bash
minikube addons enable ingress
```

```bash
kubectl get ingress
```