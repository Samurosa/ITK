## Ответственность сервисов

### UserService

Отвечает за:

- регистрация
- профиль пользователя
- баланс
- адреса
- роль пользователя

Message:


---

### SportInstrumentService

Отвечает за:

- каталог инвентаря
- цены
- остатки
- availability
- категории
  
Message:


---

### OrderService


Знает:
- кто делает заказ
- какие товары покупаются
- статус заказа

Использует:

- UserService
- SportInstrumentService

Message:
