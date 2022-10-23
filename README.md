# Тестовое задание на позицию стажёра-бэкендера Avito.Tech: микросервис для работы с балансом пользователей

## Структура проекта
 
```
.
├── assets               файлы документации README.md
├── cmd                  основное приложение
├── collections          интеграционные тесты
├── config               конфигурационные файлы
├── docs                 документация Swagger 
├── internal             приватный код библиотек и приложения
│   ├── entity           уровень сущностей и структур запросов
│   ├── handler          уровень обработчиков запросов
│   ├── repository       бизнес-логика на уровне базы данных
│   ├── service          бизнес-логика в форме интерфейсов
└── schema               метаданные о БД
```

## Руководство:

Для развёртывания приложения используется docker-compose и GNU Make:

```
$ make help
Usage:
  make <target>
  help             Output help
  compose-up       Run billing app along with PostgreSQL server using docker-compose
  compose-down     Stop billing app and PostgreSQL server launched using docker-compose
  compose-test     Run integration testing in docker environment
  run              Run billing app locally
  build            Build billing app locally
  test             Run tests on billing app locally
```

Для запуска контейнеров сервиса и базы данных необходимо выполнить команду:

```
make compose-up
```

Остановка сервисов:

```
make compose-down
```

Проект включает **интеграционный тест, созданный в Postman (collections/integration.json)**. Для выполнения тестирования используется вспомонательный контейнер "newman". Запуск тестирования и вывод результатов производится командой:

```
make compose-test
```
Примера вывода результатов интеграционного теста:
#  ![Работа интеграционных тестов:](/assets/testing-works.png) </br>

В проект интегрирована автоматическая документация **Swagger**. Она доступна по адресу: [https://localhost:8080/swagger/index.html](https://localhost:8080/swagger/index.html)
Работа Swagger:
#   ![Работа Swagger:](/assets/swagger-works.png) </br>

Для включения постоянного хранилища данных в PostgreSQL необходимо подмонтировать директорию к контейнеру БД. Для этого нужно раскомментировать в файле docker-compose.yml следующую строку:

```
#      - postgres_data:/var/lib/postgresql/data
```
## Схема базы данных:
  ![Схема базы данных:](/assets/avitoDb_schema.png) </br>


# Реализованные методы:

##### **Метод получения баланса пользователя:** *принимает id пользователя.* <br/>
  Вызов: 
  ``` 
  GET /account/getBalance/:id
  ```

  Возвращаемый JSON: 
  ``` 
  { 
    "user-balance": int,
    "user-pending-amount": int
  } 
  ```
  
  
##### **Метод начисления средств на баланс:** *принимает id пользователя и сколько средств зачислить. Возвращает баланс с двумя параметрами: настоящим балансом и кредитным. При вызове метода на несуществующий аккаунт создается запись с балансом суммы начисления и нулевым кредитным балансом (нулевым долгом).* <br/>
  Вызов: 
  ``` 
  POST /account/depositMoney 
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "update-amount": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "account-id": int,
    "sum-deposited": int, 
    "operation-status": string,
    "operation-event": string,
    "created-at": time.Time()
  } 
  ```
  **Примечания: при начислении денег на несуществующий аккаунт создается новый счет с балансом, равным начисляемой сумме. Реализована проверка на отрицательную сумму пополнения.** </br>
  
##### **Метод снятия средств с баланса:** *принимает id пользователя и сколько средств снять. Проверяет существование аккаунта и уход в отрицательный баланс.* <br/>
  Вызов: 
  ``` 
  POST /account/withdrawMoney 
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "update-amount": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "account-id": int,
    "sum-withdrawn": int, 
    "operation-status": string,
    "operation-event": string,
    "created-at": time.Time()
  } 
  ```
 **Примечания: реализованы проверки на уход баланса в отрицательное число после снятия, отрицательную сумму снятия, существование пользователя с указанным идентификатором.** </br>
 
#####  **Метод перевод средств со счета на другой:** *принимает id пользователя, перечисляющего средства, id пользователя, получающего средства, и сумму перевода. Проверяет существование аккаунта перечисления (второй аккаунт не проверяется в соответствие с методом первого начисления) и уход в отрицательный баланс.* <br/>
  Вызов: 
  ``` 
  POST /account/transfer
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "sender-id": int,
    "receiver-id": int,
    "transfer-amount": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "receive-account": int,
    "transfer-account" int,
    "amount": int,
    "status": string,
    "event-type": string,
    "created-at": time.Time()
  } 
  ```
 **Примечания: реализованы проверки на уход баланса в отрицательное число после снятия, отрицательную сумму перевода, существование пользователя-отправителя с указанным идентификатором.** </br>
 
 ##### **Метод резервирования средств с основного баланса на отдельном счете:** *принимает id пользователя, ИД услуги, ИД заказа, стоимость. Записывает в отдельную таблицу полученные параметры вместе со статусом "Pending" и таймкодом получения. Увеличивает кредитный баланс пользователя на стоимость услуги. * <br/>
 Вызов: 
  ``` 
  POST /account/reserveServiceFee
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "service-id": int,
    "order-id": int,
    "fee": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "account-id": int
    "service-id": int
    "order-id":   int
    "invoice":    int
    "status":     string
    "created-at": time.Time(),
    "updated-at": time.Time()
  } 
  ```
 **Примечания: реализованы проверки на отрицательную стоимость, существование пользователя с указанным идентификатором.** </br>
 
##### **Метод признания выручки:** *списывает из резерва деньги, добавляет данные в отчет для бухгалтерии.  Проверяет последний статус по услуге с принятыми параметрами на предмет конфликта (ранее установленный статус подтверждения "Approved" или статус отмены "Cancelled"), устанавливает статус записи с принятыми параметрами в "Approved", уменьшает кредитный баланс пользователя и пытается списать сумму со счета. В случае нехватки средств откатывает транзакцию и возвращает ошибку. Принимает id пользователя, ИД услуги, ИД заказа, сумму.* <br/>
 Вызов: 
  ``` 
  POST /account/approveServiceFee
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "service-id": int,
    "order-id": int,
    "fee": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "account-id": int
    "service-id": int
    "order-id":   int
    "invoice":    int
    "status":     string
    "created-at": time.Time(),
    "updated-at": time.Time()
  } 
  ```
 **Примечания: реализованы проверки на уход баланса в отрицательное число после списания суммы, отрицательную сумму снятия, существование пользователя с указанным идентификатором, последний статус записи с теми же параметрами.** </br>

##### **Метод разрезервирования средств:** *разрезервирует услугу для пользователя. Проверяет последний статус по услуге с принятыми параметрами на предмет конфликта (ранее установленный статус подтверждения "Approved" или статус отмены "Cancelled") Устанавливает статус по услуге в "Cancelled", уменьшает кредитный баланс на стоимость услуги (уменьшает долг). Принимает id пользователя, ИД услуги, ИД заказа, сумму.* <br/>
 Вызов: 
  ``` 
  POST /account/approveServiceFee
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "service-id": int,
    "order-id": int,
    "fee": int
  }
  ```
  Возвращаемый JSON: 
  ``` 
  { 
    "account-id": int
    "service-id": int
    "order-id":   int
    "invoice":    int
    "status":     string
    "created-at": time.Time(),
    "updated-at": time.Time()
  } 
  ```
 **Примечания: реализованы проверки на существование пользователя с указанным идентификатором, последний статус записи с теми же параметрами.** </br>

## Дополнительные задания:
##### **Метод получения месячного отчета по услугам в формате CSV-файла: получает на вход диапазон дат отчета. Возвращает CSV-файл.**
**Постановка задания:**
 ```
 Бухгалтерия раз в месяц просит предоставить сводный отчет по всем пользователем, с указанием сумм 
 выручки по каждой из предоставленной услуги для расчета и уплаты налогов.
 Задача: реализовать метод для получения месячного отчета. На вход: год-месяц. На выходе ссылка на CSV файл.
 ```
 Вызов: 
  ``` 
  POST /reports/servicesMonthly
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "date-from": "2022-10-22",
    "date-to": "2022-10-23"
  }
  ```
  Возвращаемый CSV-файл: 
  ``` 
  {
     "csv-file-link": "../reports/servicesReport10-23-2022.csv"
  }
  ```
 ![service-reports-csv-image](assets/csv-report-service.png) </br>
##### **Метод получения отчета по операциям пользователя: отчет создается в формате CSV-файла. На вход получает идентификатор пользователя, диапазон дат отчета, поле сортировки (по сумме или дате), метод сортировки (от большего к меньшего, от меньшего к большему) номер выдаваемой страницы (10 записей на странице).** 

**Постановка задания:**

 ```
  Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства.
  Задача: необходимо предоставить метод получения списка транзакций с комментариями 
  откуда и зачем были начислены/списаны средства с баланса. 
  Необходимо предусмотреть пагинацию и сортировку по сумме и дате. 
 ```
 
  Вызов: 
  ``` 
  POST /reports/transactions
  ```
  
  Принимаемый JSON: 
  ``` 
  {
    "user-id": int,
    "date-from": "2022-10-23",
    "date-to": "2022-10-23",
    "sort-by": "date", // допустимые варианты "date" | "amount"
    "sort-order": "descending", // допустимые варианты "descending" | "ascending"
    "page": 1 // страница выбора -- на страницу выводится первые десять записей
  }
  ```
  Возвращаемый CSV-файл: 
  ``` 
  {
      /reports/userTransactionsReport10-23-2022.csv
  }
  ```
  ![transactions-csv-image](assets/csv-report-transactions.png) </br>
  
 ### **Вспомогательный метод получения файла по ссылке: получает адрес на файл и возвращает его.**

  Вызов: 
  ``` 
  GET /reports/:filepath
  ```
