
# Docker Monitor
Система для мониторинга состояния Docker контейнеров.

## Описание
Система состоит из 2х основных сервисов:
 - Pinger
 - Backend

А также фронтенда на React для динамического отображения состояния контейнеров.

### Pinger
Получает список ВСЕХ Docker контейнеров на машине и пытается открыть TCP соединение через открытый порт. 
- Если контейнер поднят в сети хост-машины, то пингуется напрямую
- Если в виртуальной подсети, то через PortBinding (если есть)

Далее отправляет данные на REST API сервиса Backend.

### Backend

Сохраняет текущее состояние контейнеров в базе данных, а также лог изменения их стостояний. 

#### [GET] api/v1/container?offset=x&limit=y
Список контейнеров с их текущим состоянием.

**Тело ответа**
```json
[
  {
    "dockerId": "string",
    "id": 0,
    "ip": "string",
    "lastActivity": "string",
    "lastCheck": "string",
    "name": "string",
    "status": "online"
  }
]
```

#### [GET] api/v1/container/:id?from_date=x&to_date=y
Лог состояния контейнера за определенный период.

**Тело ответа**
```json
{
"container_id": 0,
"from_date": "string",
"to_date": "string",
"health_log": [
    {
    "status": "offline",
    "timestamp": "string"
    }
],
}
```

#### [POST] api/v1/container/health
Загрузить снепшот состояния контейнера.

**Тело запроса**
```json
{
  "docker_id": "string",
  "health_report": {
    "status": "offline",
    "timestamp": "string"
  },
  "ip": "string",
  "name": "string",
  "port": "string"
}
```

**Тело ответа при ошибке**
```json
{
  "details": "string",
  "status": "string",
  "status_code": 0
}
```

Также при сборке проекта через `task build` генерируется `swagger` документация к API.

## Сборка
Создать конфиг `liquibase.properties` в `backend/cmd/liquibase`. Есть шаблон `liquibase.properties.example` в этой же папке, можно просто скопировать его содерджимое. 

Для чувствительных данных могут использоваться конфиги `local.yml` в папке `config` проектов. Они перезаписывают основную конфигурацию и не индексируются в VCS.

Для сборки и запуска в контейнерах:
```bash
docker compose up
```

Также есть возможность сборки Golang проектов по отдельности с помощью `task`. 

```bash
task build # собрать проект + swagger для API
task run # build + запуск
```