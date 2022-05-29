# ISO
ISO - это инструмент, предназначенный для изолирования микросервисов от их внешних зависимостей для упрощения функционального и нагрузочного тестирования. 
От пользователя требуется перечислить в файле спецификации `spec.yaml` список внешних сервисов и пути к файлам-контрактам API, которые они реализуют (proto, OpenAPI и т.д.).
Поведение каждой внешней зависимости определяется с помощью простого набора правил и может быть изменено в любой момент.
---
## Как это работает
У ISO есть два основных компонента:
1. isoserver - сервис, который имитирует поведение всех внешних зависимостей изолируемого сервиса в соответсвии с заданными правилами.
2. isoctl - консольная утилита, которая на основе файла спецификации генерирует плагин для имитирующего сервиса.

## Установка
- Установить [Docker](https://docs.docker.com/install/)
- `go install github.com/Speakerkfm/iso/cmd/isoctl@latest`

## Использование
Для тестового примера взят сервис [order-api](https://github.com/Speakerkfm/iso_example/tree/master/order-api).
1. Создать проект `isoctl init example`.
2. Описать внешние зависимости в файле `spec.yaml`.
    ```
   # Данный файл содержит описание внешних зависимостей изолируемого сервиса
   external_dependencies:
     - name: user_service        # Имя внешней зависимости
       host: user-api.localhost  # Хост 
       proto:
         - raw.githubusercontent.com/Speakerkfm/iso_example/master/user-api/api/user_api.proto
   
     - name: address_service
       host: address-api.localhost
       proto:
         - raw.githubusercontent.com/Speakerkfm/iso_example/master/address-api/api/address_api.proto
   
     - name: shipment_service
       host: shipment-api.localhost
       proto:
         - raw.githubusercontent.com/Speakerkfm/iso_example/master/shipment-api/api/shipment_api.proto
   ```
3. Сгенерировать плагин для имитирующего сервиса `isoctl generate example --docker=true`.
4. Запустить имитирующий сервис `isoctl server start example --docker=true`.
5. Скачать автоматически сгенерированные правила `isoctl rules sync example`.
6. Добавить правила
    ```
   service_name: UserService
   method_name: GetUser
   rules:
     - conditions:
         - key: body.id
           value: "10"
       response:
         delay: 5ms
         data: |-
           {
           	"user": {
           		"id": 62,
           		"name": "Aleksandr",
           		"surname": "Usanin"
           	}
           }
         error: ""
     - conditions:
         - key: body.id
           value: "15"
       response:
         delay: 5ms
         error: "Not found"
     - conditions:
         - key: header.x-request-id
           value: '*'
       response:
         delay: 5ms
         data: |-
           {
           	"user": {
           		"id": 61,
           		"name": "aFChDoAtEQRJLckOHgifVhyeD",
           		"surname": "JXPAvjpxdMRsWJwNyclTHPsMO"
           	}
           }
         error: ""
   ```
7. Применить правила `isoctl rules apply example`.
8. Запустить и протестировать изолируемый сервис.
9. Получить отчет о пройденном тестировании `isoctl report load`.
    ```
   ISO Report
   +-----------------+----------------+-----------------------+---------------+
   |  SERVICE NAME   |  METHOD NAME   |       RULE NAME       | REQUEST COUNT |
   +-----------------+----------------+-----------------------+---------------+
   | AddressService  | GetUserAddress | header.x-request-id:* |             8 |
   +-----------------+----------------+-----------------------+---------------+
   | ShipmentService | CreateShipment | header.x-request-id:* |             8 |
   +-----------------+----------------+-----------------------+---------------+
   | UserService     | GetUser        | body.id:10            |             2 |
   +-----------------+----------------+-----------------------+---------------+
   | UserService     | GetUser        | body.id:15            |             1 |
   +-----------------+----------------+-----------------------+---------------+
   | UserService     | GetUser        | header.x-request-id:* |             6 |
   +-----------------+----------------+-----------------------+---------------+
   ```