DEFER
1.# Задание: Анализ кода на Go  
   
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main  
  
import (  
    "fmt"  
)  
  
func main() {  
    fmt.Println("start")  
    for i := 1; i < 4; i++ {  
       defer fmt.Println(i)  
    }  
    fmt.Println("end")  
}
// start end 1 2 3
2.
# Задание: Анализ кода на Go  
  
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.  
Исправьте код так, чтобы defer вывел обновленное значение value.
package main  
  
import "fmt"  
  
func main() {  
    value := 123  
    defer fmt.Println(value)  
    changeValue(&value)  
}  
func changeValue(value *int) {  
    *value = 456  
}
//123
import "fmt"  
  
func main() {  
    value := 123 
    changeValue(&value)   
    defer fmt.Println(value)  
}  
func changeValue(value *int) {  
    *value = 456  
}
//456
3.
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main

import (
	"errors"
	"fmt"
)

func main() {
    println("Case 1")
    case1()
    println()
    println()
    // Case 1 \n without: n\ nil n\ default error \n with: n\ nil n\ default error
    // extra error не выводиться так как метод возвращает ошибку быстрее чем срабатывает defer, поэтому ошибка все еще буде дефолтной
    
    println("Case 2")
    case2()
    println()
    println()
    // case 2 nil default error extra error extra error
    // если возвращаемое значение именновоное то defer может изменять возвращаемое значение

    println("Case 3")
    case3()
    println()
    println()
    // Case 3 nil default error first error first error
    // возвращается первая ошибка так как дефер работает по принципу LIFO

}

func case1() {
    helperWithDefer := func(isError bool) error {
        var retVal error

        defer func() {
            retVal = errors.New("Extra error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal
    }

    helperWithoutDefer := func(isError bool) error {
        var retVal error

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
}

func case2() {
    helperWithDefer := func(isError bool) (retVal error) {
        defer func() {
            retVal = errors.New("Extra error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
}

func case3() {
    helperWithDefer := func(isError bool) (retVal error) {
        defer func() {
            retVal = errors.New("First Error")
        }()

        defer func() {
            retVal = errors.New("Second Error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false))
    fmt.Println(helperWithoutDefer(true))
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false))
    fmt.Println(helperWithDefer(true))
}

-------------------------------------------------------------




ERRORS
1.
# Задание: Возврат ошибки без дополнительных пакетов  
  
Реализуйте функцию `handle() error`, которая возвращает ошибку, **не используя дополнительные пакеты** (кроме уже импортированного `fmt`).  
  
## Требования  
1. Функция `handle()` должна возвращать тип `error`.  
2. Запрещено подключать пакеты, кроме `fmt`.  
3. Ошибка должна содержать текст (не `nil`).
```go
import (

    "fmt"

)

  

func handle() error {

    return fmt.Errorf("ошибка ошибка")

  

}

func main() {

  

    err := handle()

    if err != nil {

        fmt.Println(handle())

    }

}
```
2.
# Кастомные ошибки в Go  
  
Реализуйте три типа кастомных ошибок:  
1. Простую ошибку через `errors.New()`.  
2. Ошибку с форматированием через `fmt.Errorf()`.  
3. Структуру `MyError`, реализующую интерфейс `error`.  
  
## Требования  
1. **Простая ошибка**:  
   - Создайте функцию `SimpleError() error`, возвращающую ошибку с текстом "простая ошибка".  
  
2. **Форматированная ошибка**:  
   - Создайте функцию `FormattedError(age int) error`, которая возвращает ошибку в формате: "ошибка: возраст %d недопустим".  
   - Добавьте оборачивание ошибки с `%w`.  
  
3. **Структура `MyError`**:  
   - Реализуйте метод `Error() string`.  
   - Добавьте поле `Code int` для кода ошибки.  
   - Создайте функцию `StructError() error`, возвращающую `MyError{Code: 404, Msg: "не найдено"}`.
```go
import (

    "errors"

    "fmt"

)

  

type MyError struct {

    Code int

    Msg  string

}

  

func (m MyError) Error() string {

    return fmt.Sprintf("code error: %d message: %s", m.Code, m.Msg)

}

  

func SimpleError() error {

    return errors.New("default error")

}

  

func FormattedError(age int) error {

    return fmt.Errorf("ошибка: возраст %d недопустим", age)

}

  

func StructError() error {

    m := MyError{Code: 404, Msg: "не найдено"}

    return m

}
```


3.
# Анализ цепочек ошибок в Go error.Is

Напишите функцию `ProcessError(err error)`, которая:
1. Проверяет, содержит ли ошибка в своей цепочке ошибку типа `TimeoutError` (кастомный тип).
2. Если содержит, выводит "Требуется повторная попытка".
3. Проверяет, содержит ли ошибка в своей цепочке ошибку `ErrNotFound` (стандартная ошибка).
4. Если содержит, выводит "Ресурс не найден".
5. Для остальных ошибок выводит "Неизвестная ошибка".

## Требования
1. Определите кастомные ошибки:
   ```go
   var (
       ErrNotFound   = errors.New("ресурс не найден")
       TimeoutError = errors.New("таймаут операции")
   )
   ```
   
2. Создайте функцию SimulateRequest() error, которая:
   - В 50% случаев возвращает TimeoutError, обёрнутую в fmt.Errorf("запрос не выполнен: %w", TimeoutError).
   - В 30% случаев возвращает ErrNotFound, обёрнутую в fmt.Errorf("ошибка: %w", ErrNotFound).
   - В 20% случаев возвращает новую ошибку "неизвестная ошибка".
   - Реализуйте логику анализа ошибок в ProcessError.
3. Реализуйте логику анализа ошибок в ProcessError.
```go
import (

    "errors"

    "fmt"

    "math/rand"

)

  

var (

    ErrNotFound  = errors.New("ресурс не найден")

    TimeoutError = errors.New("таймаут операции")

)

  

func ProcessError(err error) {

    switch {

    case errors.Is(err, TimeoutError):

        fmt.Println("Требуется повторная попытка", err)

    case errors.Is(err, TimeoutError):

        fmt.Println("Требуется повторная попытка", err)

    default:

        fmt.Println(err)

    }

}

  

func SimulateRequest() error {

    rnd := rand.Float64()

    if rnd <= 0.5 {

        return fmt.Errorf("запрос не выполнен: %w", TimeoutError)

    }

    if rnd <= 0.8 {

        return fmt.Errorf("ошибка: %w", ErrNotFound)

    }

    return fmt.Errorf("неизвестная ошибка")

}

  

func main() {

    ProcessError(SimulateRequest())

}
```

-------------------------------------------------------



GENERICS
# Универсальный стек (LIFO) на Go с дженериками  
  
Реализация простого обобщенного стека с методами `Push`, `Pop`, `Peek` и `IsEmpty`. Используются дженерики Go (1.18+).  
  
## Требования  
  
1. **Структура `Stack[T]`**:  
   - Поле `elements` для хранения элементов (слайс `[]T`).  
   - Методы:  
        - `NewStack[T]() *Stack[T]`: конструктор.  
        - `Push(value T)`: добавление элемента в стек.  
        - `Pop() (T, bool)`: удаление и возврат верхнего элемента (с проверкой на пустоту).  
        - `Peek() (T, bool)`: возврат верхнего элемента без удаления.  
        - `IsEmpty() bool`: проверка стека на пустоту.  
  
2. **Дополнительно**:  
   - Гарантировать безопасность операций (например, `Pop` на пустом стеке возвращает `false`).  
   - Использовать слайс для эффективного добавления/удаления элементов.

```go
import (

    "fmt"

)

  

type Stack[T any] struct {

    elements []T

}

  

func NewStack[T any]() *Stack[T] {

    return &Stack[T]{

        elements: make([]T, 0),

    }

}

  

func (s *Stack[T]) Push(value T) {

    s.elements = append(s.elements, value)

}

  

func (s *Stack[T]) Pop() (T, bool) {

    var zero T

    if len(s.elements) == 0 {

        return zero, false

    }

  

    element := s.elements[len(s.elements)-1]

    s.elements[len(s.elements)-1] = zero

    s.elements = s.elements[:len(s.elements)-1]

    return element, true

}

  

func (s Stack[T]) Peek() (T, bool) {

    var zero T

    if len(s.elements) == 0 {

        return zero, false

    }

    return s.elements[len(s.elements)-1], true

}

  

func (s Stack[T]) IsEmpty() bool {

    if len(s.elements) == 0 {

        return true

    }

    return false

}
```
----------------------------------------------------
Interface
1.
# Универсальный потокобезопасный кэш с TTL, очисткой и JSON-сериализацией  
  
Реализация in-memory кэша на Go с расширенными возможностями: автоматическое удаление ключей, очистка и сериализация данных.  
  
---  
  
## **Основные возможности**  
  
**TTL (Time-To-Live)**    
- Автоматическое удаление ключей по истечении времени жизни.   
  
**Очистка кэша**    
- Мгновенное удаление всех данных одной командой.    
  
**Сериализация в JSON**    
- Преобразование актуальных данных в JSON-формат.   
  
**Потокобезопасность**    
- Использование `sync.RWMutex` для конкурентного доступа.   
  
**Универсальное хранение**    
- Поддержка любых типов данных через `interface{}`.  
---  
  
## **Методы**  
### **Базовые операции**  
  
 - `Set(key string, value interface{}, ttl time.Duration)` Добавляет значение с указанным TTL   
 - `Get(key string) (interface{}, bool)`  Возвращает значение (с проверкой TTL)   
 - `Delete(key string)`  Удаляет конкретный ключ  
 - `Exists(key string) bool` Проверяет наличие непросроченного ключа   
  
### **Расширенные функции**  
 - `Clear()` Полностью очищает кэш   
 - `ToJSON() ([]byte, error)` Сериализует данные в JSON   
 - `GetAs[T any](key string) (T, error)`  Типизированное получение  
  
---  
  
## **Пример использования**  
```go  
func main() {  
    cache := NewCache()  
    // Добавление данных с TTL    cache.Set("user", User{Name: "Alice"}, time.Hour) // Хранится 1 час    cache.Set("temp_data", 42, time.Minute)           // Хранится 1 минуту  
    // Сериализация в JSON    jsonData, _ := cache.ToJSON()    fmt.Println(string(jsonData))    // {"temp_data":42,"user":{"Name":"Alice"}}  
  
    // Очистка кэша    cache.Clear()    fmt.Println("Exists (user):", cache.Exists("user")) // false}
```

```go
package main

  

import (

    "encoding/json"

    "fmt"

    "sync"

    "time"

)

  

type Storage struct {

    mu     sync.RWMutex

    caches map[string]Cache

}

  

type Cache struct {

    val   interface{}

    tlive time.Time

}

  

type User struct {

    Name string

}

  

type Item struct {

    Key   string      `json:"key"`

    Value interface{} `json:"value"`

}

  

func NewCache() *Storage {

    return &Storage{

        caches: make(map[string]Cache),

    }

}

  

func (c *Storage) Set(key string, value interface{}, ttl time.Duration) {

    c.mu.Lock()

    defer c.mu.Unlock()

  

    cache := Cache{val: value, tlive: time.Now().Add(ttl)}

    c.caches[key] = cache

}

  

func (c *Storage) Get(key string) (interface{}, bool) {

    c.mu.RLock()

    cache, ok := c.caches[key]

    c.mu.RUnlock()

  

    if !ok {

        return nil, false

    }

  

    if time.Now().After(cache.tlive) {

        c.mu.Lock()

        delete(c.caches, key)

        c.mu.Unlock()

        return nil, false

    }

    return cache.val, true

}

  

func (c *Storage) Delete(key string) {

    c.mu.Lock()

    defer c.mu.Unlock()

  

    delete(c.caches, key)

}

  

func (c *Storage) Exists(key string) bool {

    _, ok := c.Get(key)

    return ok

}

  

func (c *Storage) Clear() {

    c.mu.Lock()

    defer c.mu.Unlock()

  

    c.caches = make(map[string]Cache)

}

  

func (c *Storage) ToJSON() ([]byte, error) {

    res := make([]Item, 0, len(c.caches))

  

    c.mu.RLock()

    for key, cache := range c.caches {

        res = append(res, Item{Key: key, Value: cache.val})

    }

    c.mu.RUnlock()

  

    return json.Marshal(res)

}

  

func GetAs[T any](c *Storage, key string) (T, error) {

    var zero T

    value, ok := c.Get(key)

    if ok {

        c.mu.Lock()

        defer c.mu.Unlock()

  

        value, ok := value.(T)

        if !ok {

            return zero, fmt.Errorf("неправильный формат")

        }

  

        return value, nil

    }

    return zero, fmt.Errorf("неправильный формат")

}

  

func main() {

    cache := NewCache()

  

    cache.Set("user", User{Name: "Alice"}, time.Hour) // Хранится 1 час

    cache.Set("temp_data", 42, time.Minute)           // Хранится 1 минуту

  

    jsonData, _ := cache.ToJSON()

    fmt.Println(string(jsonData)) // {"temp_data":42,"user":{"Name":"Alice"}}

  

    cache.Clear()

    fmt.Println("Exists (user):", cache.Exists("user")) //false

}
```

--------------------
# Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (interface), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод каждой из предложенных программ и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
1.
```go
package main  
  
import (  
    "fmt"  
)  
  
type MyError struct {  
    data string  
}  
  
func (m *MyError) Error() string {  
    return m.data  
}  
func foo(i int) error {  
    var err *MyError  
    if i > 5 {  
       err = &MyError{data: "i>5"}  
    }  
    return err  
}  
func main() {  
    err := foo(4)  
    if err != nil {  
       fmt.Println("oops")  
    } else {  
       fmt.Println("ok")  
    }  
}
// fmt.Println("oops")

2.
package main  
  
import (  
    "fmt"  
)  
  
type errorString struct {  
    s string  
}  
  
func (e errorString) Error() string {  
    return e.s  
}  
  
func checkErr(err error) {  
    fmt.Println(err == nil)  
}  
  
func main() {  
    var e1 error  
    checkErr(e1)  //true
  
    var e *errorString  
    checkErr(e)  //false
  
    e = &errorString{}  
    checkErr(e)  //false
  
    e = nil  
    checkErr(e)  //true
}

3.
package main

import "fmt"

type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

func returnError(flag bool) error {
	if flag {
		return &CustomError{"Что-то пошло не так"}
	}
	var err *CustomError
	return err
}

func main() {
	err1 := returnError(true)
	err2 := returnError(false)

	fmt.Println("err1 == nil:", err1 == nil)//false
	fmt.Println("err2 == nil:", err2 == nil)//false

}
```

--------------------------------------------------------
PANICS
1.# Паники и их обработка в Go

В этом задании вам нужно работать с `panic` и `recover`, чтобы понять их механику.

## Требования

1. **Функция с паникой**
    - Реализуйте функцию `CausePanic()`, которая вызывает `panic("что-то пошло не так!")`.

2. **Обработка паники с `recover`**
    - Создайте функцию `HandlePanic()`, которая:
        - вызывает `CausePanic()`,
        - использует `defer` + `recover` для перехвата паники,
        - при панике выводит `"Паника перехвачена: ..."`, но программа продолжает выполняться.

3. **Запуск с демонстрацией поведения**
    - В `main()` вызовите `CausePanic()` напрямую и посмотрите, что произойдет.
    - Затем вызовите `HandlePanic()` и убедитесь, что паника обработана.


```go   
import (

    "fmt"

)

  

func CausePanic() {

    panic("что-то пошло не так!")

}

  

func HandlePanic() {

    defer func() {

        err := recover()

        if err != nil {

            fmt.Println("Паника перехвачена: ", err)

        }

    }()

  

    CausePanic()

}

  

func main() {

  

    HandlePanic()

}
```
```
```

2.# Обработка паники при делении на ноль

В этом задании вы научитесь использовать `panic` и `recover` для безопасного выполнения кода.

## Задание

1. Реализуйте функцию `SafeDivide(a, b int) int`, которая:
    - Вызывает `panic("деление на ноль")`, если `b == 0`.
    - Обрабатывает панику с `recover()` и возвращает `0` вместо аварийного завершения.

2. В `main()` протестируйте вызовы:

   SafeDivide(10, 2) // Ожидаемый результат: 5
   SafeDivide(10, 0) // Ожидаемый результат: 0 (без паники)

```go
package main

  

import (

    "fmt"

)

  

func SafeDivide(a, b int) int {

    defer func() {

        err := recover()

        if err != nil {

            fmt.Println(0)

        }

    }()

  

    return a / b

}

  

func main() {

  

    SafeDivide(10, 0)

}
```


3.# Обработка паники в многоуровневых вызовах функций

В этом задании вы разберетесь, как `panic` поднимается по стеку вызовов и как `recover` может ее обработать.

## Задание

1. Реализуйте три функции:
    - `Level1()`, которая вызывает `Level2()`.
    - `Level2()`, которая вызывает `Level3()`.
    - `Level3()`, которая вызывает `panic("ошибка в Level3")`.

2. В `Level1()` используйте `recover()` для перехвата паники и выведите:
   Паника обработана на уровне 1: ошибка в Level3

3.  В `Level2()` добавьте `defer`, который печатает `"Завершаем Level2"`, чтобы убедиться, что `defer` выполняется даже при панике.

4. В `main()` вызовите `Level1()` и убедитесь, что программа не завершилась аварийно.


```go
package main

  

import (

    "fmt"

)

  

func Level1() {

    defer func() {

        err := recover()

        if err != nil {

            fmt.Println("Паника обработана на уровне 1: ошибка в Level3")

        }

    }()

    Level2()

}

  

func Level2() {

    Level3()

    defer fmt.Println("Завершаем Level2")

}

  

func Level3() {

    panic("ошибка в Level3")

}

  

func main() {

  

    Level1()

}
```