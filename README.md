# Старцев Иван | Avito Backend Bootcamp

## Установка и конфигурация
+ Склонировать репозиторий:
  ```
  git clone https://github.com/IvanTvardovsky/avito_backend_bootcamp
  ```
+ Настроить конфигурацию в файле `configs/config.yaml`
+ Настроить переменные окружения в файле `.env` для контейнера с сервисом и БД
+ Запустить *docker compose* из корневой директории
  ```make
  docker-compose -f deployments/docker-compose.yml up
  ```
## Использование
### Сервис поддерживает 7 эндпоинтов:
+ `GET    /dummyLogin -- (no Auth)`
+ `POST   /login -- (no Auth)`
+ `POST   /register -- (no Auth)`
+ `POST   /flat/create -- (client/moderator Auth)`
+ `POST   /flat/update -- (moderator Auth)`
+ `POST   /house/create -- (moderator Auth)`
+ `GET    /house/:id -- (client/moderator Auth)`

## Запуск тестов
+ Написаны интеграционные и модульные тесты для сценариев получения списка квартир и процесса публикации новой квартиры
+ Для тестов в папке `test` существуют свои `Dockerfile` и `docker-compose.yml` с минорными изменениями
+ Токены авторизации для запросов в тестах получаются через запросы к `/dummyLogin`, а не через `/login` и `/register`, хотя они и реализованы
+ Запустить тесты из корневой директории
  ```make
  go test -v ./...
  ```
+ В интеграционных поднимаются Docker контейнеры, поэтому это может занять какое-то время 
  

## Архитектура
+ Сервис написан с использованем чистой архитектуры. 
+ Бизнес-логика расположена в папках `internal/entity` и `internal/usecase`
+ Настроено подробное логгирование ошибок
+ Используется многоступенчатая сборка Dockerfile
+ Если есть какие-то замечания или любой фидбек по архитектуре или в целом по проекту, то с радостью бы получил его в каком-нибудь формате.


### Возникшие вопросы
+ Были вопросы, связанные с моделями квартиры и дома. На мой взгляд, представленные в описании задания структуры не подходили, поэтому для ясности я добавил свои поля, сохранив связи между сущностями.
+ Получились следующие модели:
```sql
CREATE TABLE houses (
    id INT PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    developer VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    number INT NOT NULL,
    house_id INT NOT NULL,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    FOREIGN KEY (house_id) REFERENCES houses(id)
);
```
