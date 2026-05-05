INTERFACE
1.
# Базовые интерфейсы в Go

Это задание направлено на освоение работы с интерфейсами в Go.

---
## Задание 1: Реализация интерфейса `Shape`

### Описание
1. Создать интерфейс `Shape` с методами:
    - `Area() float64` — возвращает площадь фигуры.
    - `Perimeter() float64` — возвращает периметр (длину окружности для круга).
2. Реализовать интерфейс для двух фигур:
    - **Круг** (`Circle`), задаваемый радиусом.
    - **Прямоугольник** (`Rectangle`), задаваемый шириной и высотой.

### Требования
- Для `Circle`:
    - Площадь: `π * r²`
    - Периметр: `2 * π * r`
- Для `Rectangle`:
    - Площадь: `width * height`
    - Периметр: `2 * (width + height)`
- Используйте константу `math.Pi`.
package main

type Shape interface {
	Area() float64
	Perimeter() float64
}
____________
```golang
package main

import (

    "fmt"

    "math"

)

type Shape interface {

    Area() float64

    Perimeter() float64

}

type Circle struct {

    Radius float64

}

type Rectangle struct {

    Width float64

    Hight float64

}

func (c Circle) Area() float64 {

    return math.Pi * math.Pow(c.Radius, 2)

}

func (c Circle) Perimeter() float64 {

    return math.Pi * 2 * c.Radius

}

func (r Rectangle) Area() float64 {

    return r.Width * r.Hight

}

func (r Rectangle) Perimeter() float64 {

    return 2 * (r.Width + r.Hight)

}

func Calculate(s Shape) {

    fmt.Println(s.Area())

    fmt.Println(s.Perimeter())

}

func main() {

    circle := Circle{Radius: 5}

    rectangle := Rectangle{Width: 14, Hight: 5}

  

    Calculate(circle)

    Calculate(rectangle)

}
```
------------
# 2. Калькулятор платежей

Реализация интерфейса для платежей через разные Банки.

## Задание

### Цель
1. Создать интерфейс `PaymentProcessor` с методом `ProcessPayment(amount float64) error`.
2. Реализовать интерфейс для трех Банков:
    - **Sberbank**
    - **Tbank**
    - **Alfabank**

### Требования
- Каждый провайдер должен иметь уникальный идентификатор (например, `APIKey`).
- Метод `ProcessPayment` должен:
    - Возвращать `nil`, если сумма платежа положительная.
    - Возвращать ошибку `ErrInvalidAmount`, если сумма ≤ 0.
    - Возвращать ошибку `ErrProviderUnavailable`, если провайдер недоступен (заглушка). Сделать рандомный шанс, что банк недоступен.
package main

import "errors"

// Общие ошибки
var (
	ErrInvalidAmount       = errors.New("некорректная сумма платежа")
	ErrProviderUnavailable = errors.New("провайдер недоступен")
)

// PaymentProcessor - интерфейс обработки платежей
type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

```golang
package main

  

import (

    "errors"

    "fmt"

    "math/rand/v2"

)

  

type PaymentProcessor interface {

    ProcessPayment(amount float64) error

}

  

type Sber struct {

    APIKey string

}

  

type Tbank struct {

    APIKey string

}

  

type Alfabank struct {

    APIKey string

}

  

var (

    ErrInvalidAmount = errors.New("некорректная сумма платежа")

  

    ErrProviderUnavailable = errors.New("провайдер недоступен")

  

    ErrInvalidAPIkey = errors.New("неизвестный провайдер")

)

  

func (b Sber) ProcessPayment(a float64) error {

    return process(a, b.APIKey)

}

  

func (b Alfabank) ProcessPayment(a float64) error {

    return process(a, b.APIKey)

}

  

func (b Tbank) ProcessPayment(a float64) error {

    return process(a, b.APIKey)

}

  

func process(a float64, api string) error {

    if api == "" {

        return ErrInvalidAPIkey

    }

  

    if a < 0 {

        return ErrInvalidAmount

    }

  

    err := checkProviderUnavailable()

    if err != nil {

        return err

    }

  

    return nil

}

  

func checkProviderUnavailable() error {

    r := rand.Float64()

    if r < 0.2 {

        return ErrProviderUnavailable

    }

    return nil

}

func main() {

  

    processor := []PaymentProcessor{

        Tbank{APIKey: "ASDASDqw123sad1"},

        Alfabank{APIKey: "123213asdasd1213"},

        Sber{APIKey: "qweqwesadasd12321"},

        Sber{APIKey: ""},

    }

  

    for _, p := range processor {

        err := p.ProcessPayment(111)

        if err != nil {

            switch t := p.(type) {

            case Sber:

                fmt.Println("Sber error:", t.APIKey, err)

            case Alfabank:

                fmt.Println("Alfabank error:", t.APIKey, err)

            case Tbank:

                fmt.Println("Tbank error:", t.APIKey, err)

            }

        }

    }

  

}
```

--------------

# 3. Управление устройствами с интерфейсами в Go

Реализуйте систему управления различными устройствами, используя интерфейсы и методы с указателями.

## Задание

### Цель
1. Создать интерфейс `Device` с методами:
    - `UpdateOS(version string) error` — обновляет ОС устройства.
    - `GetInfo() string` — возвращает информацию об устройстве.
2. Реализовать интерфейс для трех устройств:
    - **Смартфон** (`Smartphone`)
    - **Ноутбук** (`Laptop`)
    - **Умные часы** (`Smartwatch`)

### Требования
- Каждое устройство должно иметь:
    - Поле `OSVersion string` (текущая версия ОС).
    - Поле `Model string` (модель устройства).
- Методы:
    - `UpdateOS`:
        - Обновляет `OSVersion`.
        - Возвращает ошибку `ErrUnsupported`, если обновление невозможно.
    - `GetInfo`:
        - Возвращает строку в формате: `"Модель: [модель], ОС: [версия]"`.
- **Специфичные правила**:
    - Смартфон нельзя обновить, если текущая версия ОС ≥ `"12.0"`.
    - Ноутбук поддерживает только версии ОС с префиксом `"Windows"`.
    - Умные часы нельзя обновить, если новая версия короче 5 символов.
package main

import "errors"

var (
	ErrUnsupported = errors.New("обновление недоступно")
)

type Device interface {
	UpdateOS(string) error
	GetInfo() string
}
_________

```go
package main

  

import (

    "errors"

    "fmt"

    "regexp"

    "strconv"

    "strings"

)

  

type Device interface {

    UpdateOS(version string) error

    GetInfo() string

}

  

type Smartfone struct {

    OSVersion string

    Model     string

}

  

type Laptop struct {

    OSVersion string

    Model     string

}

  

type SmartWatch struct {

    OSVersion string

    Model     string

}

  

var (

    ErrUnsupported = errors.New("обновление недоступно")

)

  

func (d *Smartfone) UpdateOS(version string) error {

    re := regexp.MustCompile(`\d+(\.\d+)?`)

    match := re.FindString(d.OSVersion)

  

    curentVersion, _ := strconv.ParseFloat(match, 64)

  

    if curentVersion >= 12.0 {

        return ErrUnsupported

    }

  

    d.OSVersion = version

    return nil

}

  

func (d *Laptop) UpdateOS(version string) error {

    IsCorectVersion := strings.Contains(version, "Windows")

  

    if !IsCorectVersion {

        return ErrUnsupported

    }

  

    d.OSVersion = version

    return nil

}

  

func (d *SmartWatch) UpdateOS(version string) error {

    count := len([]rune(version))

  

    if count < 5 {

        return ErrUnsupported

    }

  

    d.OSVersion = version

    return nil

}

  

func (d Smartfone) GetInfo() string {

  

    return getinfo(d.OSVersion, d.Model)

  

}

  

func (d Laptop) GetInfo() string {

  

    return getinfo(d.OSVersion, d.Model)

  

}

  

func (d SmartWatch) GetInfo() string {

  

    return getinfo(d.OSVersion, d.Model)

  

}

  

func getinfo(os string, model string) string {

    t := fmt.Sprintf("OS version: %s Model name: %s", os, model)

    return t

}

func main() {

  

    devices := []Device{

        &Smartfone{OSVersion: "Android 12.0", Model: "Samsung"},

        &Laptop{OSVersion: "Linux Ubuntu", Model: "Lenovo 123123"},

        &SmartWatch{OSVersion: "Fenix 5 pro", Model: "Fuu"},

    }

  

    newVers := []string{

        "Android 12.3",

        "Window 11",

        "hhhh",

    }

  

    for i, p := range devices {

        fmt.Println(p.GetInfo())

  

        err := p.UpdateOS(newVers[i])

        if err != nil {

            fmt.Println(err)

        }

        fmt.Println(p.GetInfo())

    }

  

}
```


------------

# 4. Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (interface), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод программы и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main

```
import "fmt"

// создается структура с переменной целочисленного типа
type MyStruct struct {
    MyInt int
}

// создается новое значение структуры и возращается при вызове функции
func func1() MyStruct {
    return MyStruct{MyInt: 1}
}

// возвращается указатель на структуру
func func2() *MyStruct {
    return &MyStruct{}
}

// принимается указатель и изменяет значение оригинальной структуры
func func3(s *MyStruct) {
    s.MyInt = 333
}

// принимается копия структуры и изменяется внутри метода метод ничего не возращает поэтому он бесмысленный(ничего не изменяет в основной структуре и не возвращает никаких значений)
func func4(s MyStruct) {
    s.MyInt = 923
}

// возвращает пустой указатель на структуру
func func5() *MyStruct {
    return nil
}

func main() {
    ms1 := func1()
    fmt.Println(ms1.MyInt)
    //MyInt = 1

    ms2 := func2()
    fmt.Println(ms2.MyInt)
    //MyInt = 0

    func3(ms2)
    fmt.Println(ms2.MyInt)
    //MyInt = 333

    func4(ms1)
    fmt.Println(ms1.MyInt)
    //MyInt = 1

    ms5 := func5()
    fmt.Println(ms5.MyInt)
    //panic ошибка с пустым указателем (nil pointer)
}

```
--------------------------------------------------------------------
# Maps
# 1. Задание: Работа с map в Go
## Описание
В этом задании вам нужно реализовать функции для работы с map в Go.
Вам предстоит создать, заполнить и обработать map, а затем выполнить некоторые операции с ним.
## Задачи
1. Создайте map, где ключ - это строка (имя человека), а значение - его возраст.
2. Добавьте в map несколько записей.
3. Реализуйте функцию `GetAge(name string) int`, которая возвращает возраст человека по его имени.
4. Реализуйте функцию `DeletePerson(name string)`, которая удаляет запись из map.
5. Реализуйте функцию `PrintAll()`, которая выводит все записи в map.
package main
```
// Объявите переменную для хранения map

func init() {
	// Инициализируйте map
}

// Добавление записей
func AddPerson(name string, age int) {
	// Реализуйте добавление записи
}

// Получение возраста
func GetAge(name string) int {
	// Реализуйте получение возраста
	return 0
}

// Удаление записи
func DeletePerson(name string) {
	// Реализуйте удаление
}

// Вывод всех записей
func PrintAll() {
	// Реализуйте вывод всех записей
}

func main() {
	// Тестирование функций

}
```


```go
package main

import "fmt"

var register map[string]int

func init() {

    register = make(map[string]int)
}

// Добавление записей
func AddPerson(name string, age int) {
    // Реализуйте добавление записи
    register[name] = age
}

// Получение возраста
func GetAge(name string) int {
    // Реализуйте получение возраста
    return register[name]
}

// Удаление записи
func DeletePerson(name string) {
    // Реализуйте удаление
    delete(register, name)
}

// Вывод всех записей
func PrintAll() {
    // Реализуйте вывод всех записей
    for name, age := range register {
        fmt.Println(name, age)
    }
}

func main() {
    // Тестирование функций
    AddPerson("Staniclav", 18)
    AddPerson("Georgiy", 25)
    AddPerson("Maria", 24)

    name := "Staniclav"
    fmt.Printf("%s %d years old\n", name, GetAge(name))

    DeletePerson("Georgiy")

    PrintAll()
}
```
-----------------

# 2. Частотный анализ слов

## Описание задания

В этом задании вам необходимо реализовать программу на Go, которая проводит частотный анализ слов в заданном тексте. Ваша программа должна использовать `map` для подсчета количества повторений каждого слова и выводить результат, отсортированный по убыванию частоты.

## Задачи

1. **Реализуйте функцию `WordFrequency(text string) map[string]int`:**
    - Принимает строку `text`.
    - Разбивает строку на слова (например, с помощью `strings.Fields`).
    - Подсчитывает количество повторений каждого слова.
    - Возвращает `map[string]int`, где ключ – слово, а значение – количество его вхождений.

2. **Реализуйте функцию `PrintWordFrequency(freqMap map[string]int)`:**
    - Принимает `map[string]int` с данными о частоте слов.
    - Выводит слова и их количество, отсортированные по убыванию частоты.

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

type pair struct {
    word   string
    amount int
}

var pairs []pair

// WordFrequency принимает строку текста и возвращает map с частотой слов.
func WordFrequency(text string) map[string]int {
    value := strings.Fields(text)
    counter := make(map[string]int)
    for _, v := range value {
        counter[v]++
    }
    return counter
}

// PrintWordFrequency выводит частотный анализ слов, отсортированный по убыванию частоты.
func PrintWordFrequency(freqMap map[string]int) {
    for key, count := range freqMap {
        pairs = append(pairs, pair{word: key, amount: count})
    }
  
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].amount > pairs[j].amount
    })

    for _, p := range pairs {
        fmt.Println(p.word, p.amount)
    }
}

func main() {

    text := "golang is great and golang is fast"
    PrintWordFrequency(WordFrequency(text))
}
```

# 3. Работа с map в Go: Фильтрация и инвертирование

Это задание направлено на освоение продвинутых операций с map в Go, включая фильтрацию по значениям и инвертирование ключей/значений.

## Задания

### Фильтрация по значению: `FilterByValue`

**Задача**:  
Реализуйте функцию `FilterByValue`, которая фильтрует элементы map, оставляя только те, чьи значения присутствуют в разрешённом списке.

**Требования**:
- Функция должна принимать:
    - Исходную map типа `map[int]string`.
    - Список разрешённых значений типа `[]string`.
- Возвращает новую map, содержащую только элементы с значениями из списка.
- Исходная map не должна изменяться.
- Эффективная проверка значений (используйте set для оптимизации)`make(map[string]struct{}`.

```go
package main

import (
    "fmt"
)

type pair struct {
    word   string
    amount int
}

var pairs []pair

// FilterByValue возвращает новую map, содержащую только элементы,
// значения которых присутствуют в allowedValues.
func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
    // Преобразовать allowedValues в set для быстрой проверки
    whiteList := make(map[string]struct{})
    for _, a := range allowedValues {
        whiteList[a] = struct{}{}
    }
    // Создать новую map и заполнить её подходящими элементами
    n := make(map[int]string)

    for key, value := range m {
        if _, ok := whiteList[value]; ok {
            n[key] = value
        }
    }
    return n
}

// InvertMap меняет ключи и значения местами.
// Если значения исходной map не уникальны, возвращает ошибку.
func InvertMap(m map[string]int) (map[int]string, error) {
    // Проверять уникальность значений
    // При обнаружении дубликата вернуть ошибку с описанием конфликта
    invetredMap := make(map[int]string)
    for key, value := range m {
        if _, ok := invetredMap[value]; ok {
            return nil, fmt.Errorf("duplicate value: %d", value)
        }
        
        invetredMap[value] = key
    }

    return invetredMap, nil
}

```
-------------------------------------------------------------------------
# OOP
# 1. Система управления транспортом (ООП в Go)

Реализуйте иерархию классов транспорта, используя принципы ООП: наследование (композицию), инкапсуляцию и полиморфизм.

## Задание

### Цель
1. Создать базовый интерфейс `Vehicle` с методами:
    - `StartEngine() error` — запускает двигатель.
    - `StopEngine() error` — останавливает двигатель.
    - `GetInfo() string` — возвращает информацию о транспорте.
2. Реализовать три типа транспорта:
    - **Автомобиль** (`Car`):
        - Имеет поле `Brand` (марка) и `EngineOn` (состояние двигателя).
        - Метод `Honk() string` — возвращает "Beep beep!".
    - **Грузовик** (`Truck`):
        - Наследует функциональность `Car`.
        - Добавляет поле `CargoCapacity` (грузоподъемность в тоннах).
        - Переопределяет `Honk()` — возвращает "Honk Honk!".
    - **Электрокар** (`ElectricCar`):
        - Наследует функциональность `Car`.
        - Добавляет поле `BatteryLevel` (уровень заряда в %).
        - Переопределяет `StartEngine()`: запускается только если `BatteryLevel > 5%`.

### Требования
- Используйте **композицию** для наследования (встраивание структур).
- Поля `EngineOn`, `BatteryLevel` и `CargoCapacity` должны быть **инкапсулированы** (не экспортируемы).
- Для работы с полями добавьте методы:
    - `GetEngineStatus() bool` — возвращает состояние двигателя.
    - `GetBatteryLevel() int` — возвращает уровень заряда.
    - `GetCargoCapacity() float64` — возвращает грузоподъемность.
- Напишите unit-тесты, проверяющие:
    - Корректность запуска/остановки двигателя.
    - Полиморфизм через интерфейс `Vehicle`.
    - Уникальное поведение методов (например, `Honk()`).

```
package main

import "errors"

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}
```

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
    ErrEngineOff            = errors.New("двигатель не запущен")
    ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
    StartEngine() error
    StopEngine() error
    GetInfo() string
}

type Car struct {
    Brand    string
    EngineOn bool
}

func (c Car) Honk() string {
    return "Beep beep!"
}

func (c Car) GetEngineStatus() bool {
    return c.EngineOn
}

func (c *Car) StartEngine() error {
    if c.EngineOn == true {
        return ErrEngineAlreadyRunning
    }
    c.EngineOn = true
    return nil
}

func (c *Car) StopEngine() error {
    if c.EngineOn == false {
        return ErrEngineOff
    }
    c.EngineOn = false
    return nil
}

func (c Car) GetInfo() string {

    return fmt.Sprintf("Марка машины: %s Состояние двигателя: %t", c.Brand, c.EngineOn)
}

type Truck struct {
    Car
    CargoCapacity float64
}

func (t Truck) Honk() string {
    return "Honk Honk!"
}

func (t Truck) GetCargoCapacity() float64 {
    return t.CargoCapacity
}

type ElectricCar struct {
    Car
    BatteryLevel int
}

func (ecar *ElectricCar) StartEngine() error {
    if ecar.BatteryLevel < 5 {
        return ErrLowBattery
    }
    ecar.EngineOn = true
    return nil
}

func (ecar ElectricCar) GetBatteryLevel() int {
    return ecar.BatteryLevel
}
```
-------------

# 2. Система управления пользователями и ролями (ООП в Go)

Реализуйте систему управления пользователями с различными ролями и правами доступа, используя принципы ООП: инкапсуляцию, композицию и полиморфизм.

## Задание

### Цель
1. Создать базовый интерфейс `User` с методами:
    - `GetUsername() string` — возвращает имя пользователя.
    - `HasPermission(permission string) bool` — проверяет наличие права доступа.
    - `GetRole() string` — возвращает роль пользователя.
2. Реализовать три типа пользователей:
    - **Обычный пользователь** (`BasicUser`):
        - Может читать данные (`read`), но не может их изменять.
    - **Модератор** (`Moderator`):
	        - Наследует права `BasicUser`.
        - Добавляет право `edit` (редактирование данных).
        - Может банить пользователей (`ban_user`).
    - **Администратор** (`Admin`):
        - Наследует права `Moderator`.
        - Добавляет право `delete` (удаление данных).
        - Может управлять ролями (`manage_roles`).

### Требования
- Поля, хранящие права доступа, должны быть **инкапсулированы**.
- Используйте **композицию** для наследования прав.
- Для каждого типа пользователя реализуйте:
	    - Конструктор `NewAdmin(username string)`,`NewModerator(username string)`,`NewBasicUser(username string)`.
    - Уникальные права доступа.

```
package main

// Базовый интерфейс
type User interface {
	GetUsername() string
	HasPermission(permission string) bool
	GetRole() string
}
```

```go
type User interface {
    GetUsername() string
    HasPermission(permission string) bool
    GetRole() string
}

type BasicUser struct {
    username   string
    role       string
    permission map[string]struct{}
}

type Moderator struct {
    BasicUser
}

type Admin struct {
    Moderator
}

func NewBasicUser(username string) BasicUser {
    return BasicUser{
        username: username,
        role:     "User",
        permission: map[string]struct{}{
            "read": {},
        },
    }
}

func (b BasicUser) GetUsername() string {
    return b.username
}
  
func (b BasicUser) HasPermission(permission string) bool {
    _, ok := b.permission[permission]
    return ok
}

func (b BasicUser) GetRole() string {
    return b.role
}

func NewModerator(username string) Moderator {
    moderator := NewBasicUser(username)
    
    moderator.role = "Moderator"
    
    moderator.permission["edit"] = struct{}{}
    moderator.permission["ban_user"] = struct{}{}

    return Moderator{
        BasicUser: moderator,
    }
}

func NewAdmin(username string) Admin {
    admin := NewModerator(username)

    admin.BasicUser.role = "admin"
    
    admin.BasicUser.permission["delete"] = struct{}{}
    admin.BasicUser.permission["manage_roles"] = struct{}{}

    return Admin{
        Moderator: admin,
    }
}
```
---------------------------------------------------------------
# Slice TASK1
# Задание: Анализ кода на Go

Это задание направлено на глубокое понимание работы срезов (slices), их модификации и передачи в функциях Go.  
**Ваша задача:** Определить вывод каждой из предложенных программ и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
### 1.
```
package main

import "fmt"

type account struct {
	value int
}

func main() {
// инициилизуруем срез с длиной 0 и емкостью 2
	s1 := make([]account, 0, 2)
// добавлеяем в s1 срез структуру
	s1 = append(s1, account{})
// создаем новую переменную копирующую срез s1 и добавлем в скопираванный срез s2 еще одну структуру
	s2 := append(s1, account{})
// создаем переменную которая ссылается на первый объект в срезе s1 (account{})
	acc := &s2[0]
// задаем значение 
	acc.value = 100
	fmt.Println(s1, s2) //[{100}] [{100} {0}]
	s1 = append(s1, account{})
	acc.value += 100
	fmt.Println(s1, s2) ////[{200} {0}] [{200} {0}]
}
```
-----
2.
```
package main

import "fmt"

func main() {
	slice := make([]string, 0, 5)
	slice = append(slice, "0", "1", "2", "3")
	fmt.Println(slice, len(slice), cap(slice)) //[0 1 2 3] 4 5
	addToSlice1(slice)
	fmt.Println(slice, len(slice), cap(slice)) //[0 1 2 one] 4 5
	addToSlice2(slice)
	fmt.Println(slice, len(slice), cap(slice)) //[0 1 2 one] 4 5
}
// создает срез слайса [1 2] добавляет one но не возвращает его и поэтому 3 заменяет его на one, длина и капасити не меняется
func addToSlice1(slice []string) {
	slice = append(slice[1:3], "one")
}
//  ничего не делает
func addToSlice2(slice []string) {
	slice = append(slice, "two")
}
```
---
3.
```
package main

import "fmt"

func main() {
	a1 := make([]int, 0, 10)
	a1 = append(a1, []int{1, 2, 3, 4, 5}...)
	a2 := append(a1, 6)
	a3 := append(a1, 7)
	fmt.Println(a1, a2, a3) //[0 1 2 3 4 5] [0 1 2 3 4 5 7] 
	[0 1 2 3 4 5 7]
}
```
---
4.
```
package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := a[:2]
	b = append(b, 4)
	fmt.Println(b) //[1 2 4]
	fmt.Println(a) //[1 2 4]
}
```
-----
5.
```
package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	src := arr[:1]
	foo(src)
	fmt.Println(src) //[1]
	fmt.Println(arr) //[1 5 3]
}

func foo(src []int) {
	src = append(src, 5)
}
```
-----
6.
```
package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	bar := arr[1:3]
	bar = append(bar, 10, 11, 12, 13)
	fmt.Println(arr, bar) //[1 2 3 4 5] [2 3 10 11 12 13]
}
```
-----
7.
```
package main

import "fmt"

func main() {
	a := []string{"a", "b", "c"}
	b := a[1:2]
	fmt.Println(b, cap(b), len(b)) //[b] 2 1
	b[0] = "q"
	fmt.Println(a) //[a q c]
}
```
---
8.
```
package main

import (
	"fmt"
)

func main() {
	nums := make([]int, 1, 3)
	fmt.Println(nums) //[0]
	appendSlice(nums, 1)
	fmt.Println(nums) //[0]
	copySlice(nums, []int{2, 3})
	fmt.Println(nums) //[2]
	mutateSlice(nums, 1, 4)
	fmt.Println(nums) //out of range
}

func appendSlice(sl []int, val int) {
	sl = append(sl, val)
}

func copySlice(sl, cp []int) {
	copy(sl, cp)
}

func mutateSlice(sl []int, idx, val int) {
	sl[idx] = val
}
```
---
9.
```
package main

import (
	"fmt"
)

func main() {
	slice := make([]int, 3, 4)
	appendingSlice(slice[:1])
	fmt.Println(slice) //[0 1 0]
}

func appendingSlice(slice []int) {
	slice = append(slice, 1)
}
```
--------------------------------
SLICE TASK 2
# Задание: Анализ и исправление кода на Go

Это задание направлено на понимание работы срезов, функций и передачи данных в Go.  
**Ваша задача:**
1. **Проанализировать вывод программ** и объяснить поведение кода.
2. **Исправить код** так, чтобы достигался корректный результат (в некоторых случаях требуется несколько решений)


### 1.// Версия 1.21
```
package main

import (
	"fmt"
)

func main() {
	var numbers []*int
	for _, value := range []int{10, 20, 30, 40} {
		numbers = append(numbers, &value)
	}
	for _, number := range numbers {
		fmt.Println("d", *number)
	}
}
```
не понял что тут нужно исправлять
сначала мы создаем слайс с указателями типа инт 
после чего заполняем данными и выводим

----
### 2.
```
package main

import (
	"fmt"
	"strings"
)

func chengeSlice(arr []string) {
	arr[0] = "Goodbye"
}

func appendSomeData(arr []string) {
	arr = append(arr, "!")
}

func main() {
	someSlice := []string{"Hello", "World"}
	chengeSlice(someSlice)
	appendSomeData(someSlice)
	fmt.Println(strings.Join(someSlice, ""))
}
```

```go
func chengeSlice(arr []string) {

    arr[0] = "Goodbye"

}

  

func appendSomeData(arr []string) []string {

    return append(arr, "!")

}

  

func main() {

    someSlice := []string{"Hello", "World"}

    chengeSlice(someSlice)

    someSlice = appendSomeData(someSlice)

    fmt.Println(strings.Join(someSlice, " "))

}
```
----
### 3.
```
package main

import "fmt"

func test(testSlice []string) {
	testSlice = append(testSlice, "Пока")
}
func main() {
	testSlice := make([]string, 0, 3)
	testSlice = append(testSlice, "Привет")
	testSlice = append(testSlice, "Привет")
	test(testSlice)
	fmt.Println(testSlice)
}
```

```go
func test(testSlice []string) []string {

    return append(testSlice, "Пока")

}

func main() {

    testSlice := make([]string, 0, 3)

    testSlice = append(testSlice, "Привет")

    testSlice = append(testSlice, "Привет")

    testSlice = test(testSlice)

    fmt.Println(strings.Join(testSlice, " "))

}
```
----
### 4.
```
package main

import "fmt"

func main() {
	first := []int{10, 20, 30, 40}
	second := make([]*int, len(first))
	for i, v := range first {
		second[i] = &v
	}
	fmt.Println(*second[0], *second[1])
}
```
```go
func main() {

    first := []int{10, 20, 30, 40}

    second := make([]*int, len(first))

    for i, v := range first {

        second[i] = &v

    }

    fmt.Println(DerefSlice(second))

}

  

func DerefSlice(sl []*int) []int {

    res := make([]int, 0, len(sl))

    for _, v := range sl {

        if v != nil {

            res = append(res, *v)

        }

    }

    return res

}
```
----
### 5.
```
package main

import (
	"fmt"
)

func main() {
	slice := make([]string, 3, 4)
	fmt.Println(slice)

	appendSlice(slice)
	fmt.Println(slice)

	mutareSlice(slice)
	fmt.Println(slice)
}

func appendSlice(slice []string) {
	slice = append(slice, "privet")
}
func mutareSlice(slice []string) {
	slice[0] = "vasya"
}
```

```go
func main() {

    slice := make([]string, 3, 4)

    fmt.Println(slice) //["" "" ""_]

    slice2 := slice[:1]

    appendSlice(slice2)

    fmt.Println(slice) //["" "privet" ""_]

  

    mutareSlice(slice)

    fmt.Println(slice) //["vasya" "" ""_]

  

}

  

func appendSlice(slice []string) {

    slice = append(slice, "privet")

}

func mutareSlice(slice []string) {

    slice[0] = "vasya"

}
```

-----------------------------
### SLICE TASK 3
# 1. Курс Go: Удаление элементов из слайса

Это задание поможет освоить работу со слайсами в Go, фокусируясь на операциях удаления элементов с учетом эффективности, порядка и управления памятью.

## Цели
- Научиться удалять элементы из слайса с сохранением и без сохранения порядка.
- Понять, как избежать утечек памяти при работе с указателями.
- Оптимизировать использование памяти слайса.
- Реализовать продвинутые операции (удаление дубликатов, фильтрация).

## Задание

### Часть 1: Базовое удаление
Реализуйте функции:
- `RemoveUnordered(s []T, i int) []T` — удаление без сохранения порядка.
- `RemoveOrdered(s []T, i int) []T` — удаление с сохранением порядка.

### Часть 2: Удаление по значению
Реализуйте функцию:
- `RemoveAllByValue(s []T, value T) []T` — удаление всех вхождений `value`.

### Часть 3: Работа с памятью
1. Обнуляйте удаленные элементы-указатели.
2. Сокращайте вместимость (`capacity`) слайса при сильном уменьшении размера.

### Часть 4: Дополнительные задачи
1. `RemoveDuplicates(s []T) []T` — удаление дубликатов.
2. `RemoveIf(s []T, predicate func(T) bool) []T` — удаление по условию.
package main

```
// RemoveUnordered удаляет элемент по индексу без сохранения порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveUnordered[T any](s []T, i int) []T {
	// реализовать
	return s
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveOrdered[T any](s []T, i int) []T {
	// реализовать
	return s
}

// RemoveAllByValue удаляет все вхождения указанного значения.
func RemoveAllByValue[T comparable](s []T, value T) []T {
	// реализовать
	return s
}

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).
func RemoveDuplicates[T comparable](s []T) []T {
	// реализовать
	return s
}

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	// реализовать
	return s
}

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),
// обнуляя удаляемый элемент для предотвращения утечек памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	//реализовать
	return s
}

// ShrinkCapacity сокращает вместимость слайса, если она превышает
// удвоенную длину после удаления элементов.
func ShrinkCapacity[T any](s []T) []T {
	//реализовать
	return s
}

func main() {
	//реализовать
}
```
```go
package main

  

import (

    "fmt"

)

  

// RemoveUnordered удаляет элемент по индексу без сохранения порядка.

// Если индекс выходит за границы слайса, возвращает исходный слайс.

func RemoveUnordered[T any](s []T, i int) []T {

    //[1 2 3 4 5 6] -> [1 2 6 4 5]

    if i < 0 || i >= len(s) {

        return s

    }

    lstElmnt := len(s) - 1

    s[i] = s[lstElmnt]

    return s[:lstElmnt]

}

  

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.

// Если индекс выходит за границы слайса, возвращает исходный слайс.

func RemoveOrdered[T any](s []T, i int) []T {

    // [1 2 3 4 5 6] -> [1 2 4 5 6]

    return append(s[:i], s[i+1:]...)

}

  

// RemoveAllByValue удаляет все вхождения указанного значения.

func RemoveAllByValue[T comparable](s []T, value T) []T {

    res := make([]T, 0, len(s))

  

    for i, j := range s {

        if s[i] == value {

            continue

        }

        res = append(res, j)

    }

  

    return res

}

  

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).

func RemoveDuplicates[T comparable](s []T) []T {

    m := make(map[T]struct{})

    res := make([]T, 0, len(s))

  

    for _, sl := range s {

        if _, ok := m[sl]; !ok {

            res = append(res, sl)

        }

        m[sl] = struct{}{}

    }

    return res

}

  

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.

func RemoveIf[T any](s []T, predicate func(T) bool) []T {

    res := make([]T, 0, len(s))

    for _, sl := range s {

        if !predicate(sl) {

            res = append(res, sl)

        }

    }

    return res

}

  

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),

// обнуляя удаляемый элемент для предотвращения утечек памяти.

func RemoveOrderedWithNil[T any](s []*T, i int) []*T {

    copy(s[i:], s[i+1:])

    s[len(s)-1] = nil

    return s[:len(s)-1]

}

  

// ShrinkCapacity сокращает вместимость слайса, если она превышает

// удвоенную длину после удаления элементов.

func ShrinkCapacity[T any](s []T) []T {

    if cap(s) < 2*len(s) {

        return s

    }

    res := make([]T, len(s))

    copy(res, s)

    return res

}

  

func main() {

    slice := []int{1, 2, 3, 4, 5, 6}

    res := RemoveUnordered(slice, 2)

    fmt.Println(res)

  

    res = RemoveOrdered(slice, 2)

    fmt.Println(res)

  

    slice = []int{1, 2, 2, 4, 2, 6}

    res = RemoveAllByValue(slice, 2)

    fmt.Println(res)

  

    res = RemoveDuplicates(slice)

    fmt.Println(res)

  

    res = RemoveIf(slice, func(i int) bool { return i/2 != 0 })

    fmt.Println(res)

  

    var slice2 []*string

    for _, n := range []string{"en", "ru", "uk", "bg", "tt"} {

        slice2 = append(slice2, &n)

    }

    slice2 = RemoveOrderedWithNil(slice2, 2)

    fmt.Println(DerefSlice(slice2))

  

    res = ShrinkCapacity(slice)

    fmt.Println(res)

}

  

func DerefSlice[T any](sl []*T) []T {

  

    res := make([]T, 0, len(sl))

  

    for _, v := range sl {

  

        if v != nil {

  

            res = append(res, *v)

  

        }

  

    }

  

    return res

  

}
```
------------------------------
# SLICE TASK 4
### Работа со слайсами в Go

Этот проект демонстрирует различные способы работы со слайсами в Go, включая очистку, обнуление и особенности внутренней структуры.

**Ваша задача:** Определить вывод каждого случая и зафиксировать ответы **в сообщениях коммитов** с кратким объяснением логики.
package main

```
import (
	"fmt"
	"unsafe"
)

func main() {
	//1 first: [] : 0 : 5
	first := []int{1, 2, 3, 4, 5}
	first = nil
	fmt.Println("first:", first, ":", len(first), ":", cap(first))

	//2 second: [] : 0 : 5
	second := []int{1, 2, 3, 4, 5}
	second = second[:0]
	fmt.Println("second:", second, ":", len(second), ":", cap(second))

	//3 third: [0 0 0 0 0] : 5 : 5
	third := []int{1, 2, 3, 4, 5}
	clear(third)
	fmt.Println("third:", third, ":", len(third), ":", cap(third))

	//4 fourth: [1 0 0 4 5] : 5 : 5
	fourth := []int{1, 2, 3, 4, 5}
	clear(fourth[1:3])
	fmt.Println("fourth:", fourth, ":", len(fourth), ":", cap(fourth))

	//5 slice = [10 0 0] 3 6
	//array =[0 0 0] 3 3
	slice := make([]int, 3, 6)
	array := [3]int(slice[:3])
	slice[0] = 10

	fmt.Println("slice = ", slice, len(slice), cap(slice))
	fmt.Println("array =", array, len(array), cap(array))

	//6 В каких случаях Slice пустой или нулевой
	//1 var data []string: nil
	//tempty= true nil= false size=0 data=nil pointer\n", data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)
	var data []string
	fmt.Println("var data []string:")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//2 true true 0 nil
	data = []string(nil)
	fmt.Println("data = []string(nil):")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//3 true false 0 addres
	data = []string{}
	fmt.Println("data = []string{}")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//4 true false 0 addres
	data = make([]string, 0)
	fmt.Println("data =make([]string,0)")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data))
	//addres
	empty := struct{}{}
	fmt.Println("empty struct address ", unsafe.Pointer(&empty))
}
```


---------------------
# SLICE TASK 5
### Задача

Реализуйте структуру стека с использованием слайсов, удовлетворяющую следующему интерфейсу:

```go
type Stacker interface {
    Push(v int)
    Pop() int
}
```

### Требования к реализации

1. Операция Push(v int)
    Должна добавлять целочисленное значение v в стек.

2. Операция Pop() int Должна возвращать последний добавленный элемент, реализуя поведение LIFO (последним пришёл — первым ушёл).
    Если стек пуст, вызов метода Pop() должен приводить к панике.

3. Конструктор
    Реализуйте функцию New() *stack, возвращающую новый экземпляр стека.

4. Реализация должна находится в main.go
5. Реализация должна успешно проходить тесты. Для их запуска введите команду `go test ./...` в этой директории


```
package main

import (
	"testing"
)

func TestStack_PushPop(t *testing.T) {
	s := New()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	tests := []struct {
		expected int
	}{
		{3},
		{2},
		{1},
	}

	for _, tc := range tests {
		got := s.Pop()
		if got != tc.expected {
			t.Errorf("Pop() = %d; ожидалось %d", got, tc.expected)
		}
	}
}

func TestStack_PopEmpty(t *testing.T) {
	s := New()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Ожидалась паника при попытке извлечь элемент из пустого стека")
		}
	}()

	s.Pop()
}
```

```
package main

type Stacker interface {
	Push(v int)
	Pop() int
}

type stack struct {
	//...
}

func (s *stack) Push(v int) {
	panic("unimplemented")
}

func (s *stack) Pop() int {
	panic("unimplemented")
}

func New() *stack {
	return &stack{}
}
```



```go
package main

  

type Stacker interface {

    Push(v int)

    Pop() int

}

  

type stack struct {

    slice []int

}

  

func (s *stack) Push(v int) {

    s.slice = append(s.slice, v)

}

  

func (s *stack) Pop() int {

    if len(s.slice) == 0 || s.slice == nil {

        panic("стек пуст")

    }

    res := s.slice[len(s.slice)-1]

    s.slice = s.slice[:len(s.slice)-1]

    return res

}

  

func New() *stack {

    return &stack{slice: make([]int, 0)}

}
```
