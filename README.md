# Тестовое задание на позицию стажёра-бэкендера

## Микросервис для работы с балансом пользователей

**Основное задание (минимум):**

Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость.
Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму.
Метод получения баланса пользователя. Принимает id пользователя.

*Доп. задание 1:*

Бухгалтерия раз в месяц просит предоставить сводный отчет по всем пользователем, с указанием сумм выручки по каждой из предоставленной услуги для расчета и уплаты налогов.

Задача: реализовать метод для получения месячного отчета. На вход: год-месяц. На выходе ссылка на CSV файл.

Пример отчета:

название услуги 1;общая сумма выручки за отчетный период

название услуги 2;общая сумма выручки за отчетный период

*Доп. задание 2:*

Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства. 

Задача: необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате. 
