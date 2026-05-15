# CHANNELS
## Задание: Анализ и исправление кода с гонками данных

### Описание задачи
1. Внимательно изучить код.
2. Найти все ошибки, описать их в комментариях прямо в коде.
3. Исправить код, обеспечив корректную работу.

```golang
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	alreadyStored := make(map[int]struct{})
	capacity := 1000
	doubles := make([]int, 0, capacity)
	for i := 0; i < capacity; i++ {
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, capacity)
	wg := sync.WaitGroup{}
	for i := 0; i < capacity; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := alreadyStored[doubles[i]]; !ok {
				alreadyStored[doubles[i]] = struct{}{}
				uniqueIDs <- doubles[i]
			}
		}()
	}
	wg.Wait()
	for val := range uniqueIDs {
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs)
}
```

```go
package main

  

import (

    "fmt"

    "math/rand"

    "sync"

)

  

func main() {

    alreadyStored := make(map[int]struct{})

    capacity := 1000

    doubles := make([]int, 0, capacity)

    for i := 0; i < capacity; i++ {

        doubles = append(doubles, rand.Intn(10))

    }

    uniqueIDs := make(chan int, capacity)

    wg := sync.WaitGroup{}

    var mu sync.RWMutex

    for i := 0; i < capacity; i++ {

        wg.Add(1)

        go func() {

            defer wg.Done()

            mu.RLock()

            _, ok := alreadyStored[doubles[i]]

            mu.RUnlock()

            if !ok {

                mu.Lock()

                alreadyStored[doubles[i]] = struct{}{}

                mu.Unlock()

                uniqueIDs <- doubles[i]

            }

        }()

    }

    wg.Wait()

    close(uniqueIDs)

    for val := range uniqueIDs {

        fmt.Println(val)

    }

  

    fmt.Println(uniqueIDs)

}
```
---
## Задание: Анализ времени выполнения параллельных операций в Go

**Цель задания**:  
Понять, как работают горутины и каналы в Go, определить время выполнения программы и исправить код для достижения ожидаемого результата.

---

### Описание задачи

Данная программа создает два параллельных "воркера", каждый из которых через 3 секунды отправляет число в канал. 
Задача:
1. Определить, какое число выведет программа (сколько секунд она будет выполняться).
2. Объяснить, почему результат получился именно таким.
3. Исправить код (если необходимо), чтобы время выполнения составило **ровно 3 секунды**.
```go
package main

import (
	"time"
)

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()
	return ch
}
func main() {
	timeStart := time.Now()
	_, _ = <-worker(), <-worker()
	println(int(time.Since(timeStart).Seconds()))
}
```

```go
package main

import (
    "sync"
    "time"
)

var wg sync.WaitGroup

func worker() chan int {
  
    wg.Add(1)
    ch := make(chan int)
    go func() {
        defer wg.Done()
        time.Sleep(3 * time.Second)
        ch <- 42
    }()

    return ch
}
func main() {
    timeStart := time.Now()
    ch1 := worker()
    ch2 := worker()

    <-ch1
    <-ch2
 
    wg.Wait()
    println(int(time.Since(timeStart).Seconds()))
}
```
---
## Задание: Конвейер чисел на Go

**Цель задания**:  
Реализовать конвейер для обработки чисел с использованием горутин и каналов. Числа из первого канала должны читаться по мере поступления, обрабатываться (например, возводиться в квадрат) и записываться во второй канал.

---

### Описание задачи

Даны два канала:
- `naturals` (для передачи исходных чисел),
- `squares` (для передачи обработанных чисел).

Необходимо:
1. **Генерировать** числа и отправлять их в канал `naturals`.
2. **Читать** числа из `naturals`, обрабатывать их (возводить в квадрат) и отправлять результат в `squares`.
3. **Выводить** результаты из `squares` в консоль.

```go
package main

func main() {
	naturals := make(chan int)
	squares := make(chan int)
}
```

```go
package main

  

import (

    "fmt"

    "math/rand"

)

  

func generate(amountNumbers int) <-chan int {

    out := make(chan int)

  

    go func() {

        defer close(out)

        for i := 0; i < amountNumbers; i++ {

            out <- rand.Intn(100)

        }

    }()

    return out

}

  

func square(in <-chan int) <-chan int {

    out := make(chan int)

  

    go func() {

        defer close(out)

        for n := range in {

            out <- n * n

        }

    }()

    return out

}

  

func print(in <-chan int) {

  

    for n := range in {

        fmt.Println(n)

    }

  

}

  

func main() {

    naturales := generate(100)

    squares := square(naturales)

  

    print(squares)

}
```
---
## Задание: Объединение каналов в Go

**Цель задания**:  
Написать функцию `mergeChannels`, которая объединяет данные из нескольких каналов в один общий канал, используя паттерн `FAN-IN`.

---

### Описание задачи

Дано:
- `n` каналов типа `<-chan int`.
- Функция должна вернуть канал `<-chan int`, в который попадают все значения из исходных каналов.

Требования:
1. Все значения из входных каналов должны быть отправлены в выходной канал.
2. Выходной канал должен быть закрыт после завершения всех входных каналов.
3. Решение должно быть потокобезопасным и эффективным.

```go
package main

func mergeChannels(channels ...<-chan int) <-chan int {

}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
}
```

```go
package main

  

import (

    "sync"

)

  

func mergeChannels(channels ...<-chan int) <-chan int {

    out := make(chan int)

    wg := sync.WaitGroup{}

    wg.Add(len(channels))

    for _, in := range channels {

        go func(ch <-chan int) {

            defer wg.Done()

            for n := range ch {

                out <- n

            }

        }(in)

    }

  

    go func() {

        wg.Wait()

        close(out)

    }()

  

    return out

}

  

func main() {

    a := make(chan int)

    b := make(chan int)

    c := make(chan int)

  

    mergeChannels(a, b, c)

}
```
---
## Задание: Параллельный подсчет слов в файлах с использованием паттерна Fan-Out

### Цель задания
Реализовать параллельную обработку текстовых файлов с использованием паттерна **Fan-Out**, чтобы ускорить подсчет слов в каждом файле.

### Описание задачи
Есть директория с текстовыми файлами. Нужно:
1. Прочитать все файлы.
2. Распределить их обработку между несколькими горутинами.
3. Подсчитать количество слов в каждом файле.
4. Вывести общую статистику.

### Требования
- Использовать паттерн **Fan-Out** для распределения задач.
- Обработка каждого файла должна выполняться в отдельной горутине.
- Результаты должны агрегироваться в основном потоке.

```go
package main

  

import (

    "bytes"

    "fmt"

    "os"

    "runtime"

    "sync"

)

  

const FilePath string = "../files/"

  

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {

    defer wg.Done()

  

    for fileName := range jobs {

  

        fileDirectory := fmt.Sprintf("%s%s", FilePath, fileName)

        text, err := os.ReadFile(fileDirectory)

        if err != nil {

            fmt.Printf("file %s is not valid error: %v", fileDirectory, err)

        }

  

        countWordsInFiles := len(bytes.Fields(text))

  

        results <- fmt.Sprintf("File name: %s Count worlds in file: %d", fileName, countWordsInFiles)

    }

}

  

func main() {

    files, err := os.ReadDir(FilePath)

    if err != nil {

        fmt.Printf("directory is not valid: %v", err)

    }

    workers := runtime.NumCPU()

  

    wg := sync.WaitGroup{}

  

    results := make(chan string)

    jobs := make(chan string)

  

    go func() {

        defer close(jobs)

        for _, f := range files {

  

            jobs <- f.Name()

        }

    }()

  

    for i := 0; i < workers; i++ {

        wg.Add(1)

        go worker(i, jobs, results, &wg)

    }

  

    go func() {

        wg.Wait()

        close(results)

    }()

  

    for n := range results {

        fmt.Println(n)

    }

  

}
```

---
## Задание: Реализация паттерна "Tee" для записи в несколько реплик БД

**Цель задания**:  
Реализовать паттерн "Разветвитель", при котором данные из одного источника параллельно записываются в несколько реплик базы данных (имитированных каналами).



---

### Описание задачи

Есть сервис, который записывает данные в кластер БД, состоящий из нескольких реплик. Требуется:
1. Принимать данные из входного канала.
2. Параллельно отправлять их во все реплики (каналы).
3. Гарантировать, что данные записаны во все реплики.
4. Корректно закрыть реплики после завершения работы.

```go
package main

import (
	"fmt"
	"time"
)

// Реплика БД (имитация)
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}
}
```

```go
package main

  

import (

    "fmt"

    "sync"

    "time"

)

  

var wg sync.WaitGroup

  

// Реплика БД (имитация)

func dbReplica(name string, in <-chan int) {

    defer wg.Done()

    for data := range in {

        fmt.Printf("Запись в %s: %d\n", name, data)

        time.Sleep(100 * time.Millisecond) // Имитация задержки записи

    }

    fmt.Printf("Реплика %s закрыта\n", name)

}

  

func tee(in <-chan int, replicas []chan int) {

  

    for n := range in {

        for _, channel := range replicas {

            channel <- n

        }

    }

  

    for _, channel := range replicas {

        close(channel)

    }

}

  

func main() {

    input := make(chan int) // Канал для входящих данных

    replicas := []chan int{ // Реплики БД (каналы)

        make(chan int),

        make(chan int),

        make(chan int),

    }

    data := [6]int{1, 2, 3, 4, 5, 6}

  

    for i, replica := range replicas {

  

        wg.Add(1)

        go dbReplica(fmt.Sprintf("replica %d", i+1),

            replica,

        )

    }

  

    go tee(input, replicas)

  

    go func() {

        defer close(input)

        for _, i := range data {

            input <- i

        }

    }()

  

    wg.Wait()

}
```
---
## Задание: Реализация декоратора для преобразования метрик в реальном времени

**Цель задания**:  
Создать гибкий декоратор для каналов, который будет автоматически преобразовывать метрики серверов из байтов в мегабайты перед отправкой в API. Используя паттерн `TRANSFORMER`
---

### Описание задачи

В системе мониторинга серверов:
1. **Источник данных**: Канал `metrics <-chan ServerMetric` получает метрики в формате:
   ```go
   type ServerMetric struct {
       Name  string  // Название метрики (например, "memory_usage")
       Value float64 // Значение в байтах
   }
   ```

```go
package main

  

import (

    "fmt"

)

  

type ServerMetric struct {

    Name  string  // Название метрики (например, "memory_usage")

    Value float64 // Значение в байтах

}

  

func generate(nameMetric string, value float64) <-chan ServerMetric {

    serverMetricObject := ServerMetric{Name: nameMetric, Value: value}

  

    out := make(chan ServerMetric)

  

    go func() {

        defer close(out)

        out <- serverMetricObject

    }()

    return out

}

  

func ConvertBytesToMB(in <-chan ServerMetric) <-chan ServerMetric {

    out := make(chan ServerMetric)

  

    go func() {

        defer close(out)

        for n := range in {

            n.Value /= 1024 * 1024

            out <- n

        }

    }()

  

    return out

}

  

func main() {

  

    inputData := generate("memory_usage", 1024)

  

    transformResult := ConvertBytesToMB(inputData)

  

    for result := range transformResult {

        fmt.Printf("Metric name: %s\nMetric value: %f",

            result.Name,

            result.Value,

        )

    }

  

}
```
# Задание: Реализация конвейерной обработки данных (Pipeline паттерн)

**Цель задания**:  
Создать конвейер из трех этапов для обработки строковых данных:
1. **Парсинг** — добавление метки "parsed" к данным.
2. **Разделение** — распределение данных между N каналами (round-robin).
3. **Отправка** — параллельная обработка данных в N горутинах с добавлением метки "sent".

---

## Описание задачи

Ваша задача — реализовать систему, которая:
- Обрабатывает данные в строгом порядке: **Parse → Split → Send**.
- Корректно закрывает все каналы после завершения работы.
- Гарантирует потокобезопасность и отсутствие утечек горутин.

### Этапы конвейера

1. **Parse**:
   - Принимает канал сырых данных (`<-chan string`).
   - Добавляет к каждой строке префикс "parsed - ".
   - Возвращает канал обработанных данных.

2. **Split**:
   - Принимает канал данных и число `N` (количество выходных каналов).
   - Распределяет данные между `N` каналами в порядке round-robin.
   - Возвращает слайс каналов (`[]<-chan string`).

3. **Send**:
   - Принимает слайс каналов и запускает `N` горутин.
   - Каждая горутина добавляет к данным префикс "sent - ".
   - Возвращает объединенный канал результатов.
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот ресурсы и пиши по вопросам. Можем отдельно встречу организовать по вопросам:
https://www.youtube.com/watch?v=luQlkud-jKE&t=5s
https://habr.com/ru/companies/pt/articles/764850/
```go
package main

  

import (

    "fmt"

    "sync"

)

  

func generate(data []string) <-chan string {

    out := make(chan string)

  

    go func() {

        defer close(out)

        for _, d := range data {

            out <- d

        }

    }()

    return out

}

  

// - Принимает канал сырых данных (`<-chan string`).

// - Добавляет к каждой строке префикс "parsed - ".

// - Возвращает канал обработанных данных.

func parse(in <-chan string) <-chan string {

    out := make(chan string)

    go func() {

        defer close(out)

        for n := range in {

            out <- fmt.Sprintln("parsed - ", n)

        }

    }()

    return out

}

  

// - Принимает канал данных и число `N` (количество выходных каналов).

// - Распределяет данные между `N` каналами в порядке round-robin.

// - Возвращает слайс каналов (`[]<-chan string`).

func split(in <-chan string, amountChannels int) []chan string {

    channels := make([]chan string, 0, amountChannels)

    var wg sync.WaitGroup

    wg.Add(1)

  

    for i := 0; i < amountChannels; i++ {

        channels = append(channels, make(chan string))

    }

  

    go func() {

        defer wg.Done()

        i := 0

        for n := range in {

            channels[i] <- n

            i++

            if i >= amountChannels {

                i = 0

            }

        }

    }()

  

    go func() {

        wg.Wait()

        for i := 0; i < amountChannels; i++ {

            close(channels[i])

        }

    }()

    return channels

}

  

// - Принимает слайс каналов и запускает `N` горутин.

// - Каждая горутина добавляет к данным префикс "sent - ".

// - Возвращает объединенный канал результатов.

func send(channelsArray []chan string) <-chan string {

    out := make(chan string)

    var wg sync.WaitGroup

    wg.Add(len(channelsArray))

  

    for _, channel := range channelsArray {

        go func(chan string) {

            defer wg.Done()

            for n := range channel {

                out <- fmt.Sprintln("sent - ", n)

            }

        }(channel)

    }

  

    go func() {

        wg.Wait()

        close(out)

    }()

    return out

}

  

func main() {

    inputData := generate([]string{"qwe", "asd", "zxc", "ewq", "dsa", "cxz"})

  

    result := send(

        split(parse(inputData),

            4,

        ),

    )

  

    for n := range result {

        fmt.Println(n)

    }

}
```
--------------------------------------------
# CONCURRENCY
## Реализация потокобезопасного кеша 

### Описание задачи
Ваша задача — реализовать потокобезопасный кеш для хранения данных в формате ключ-значение. Кеш должен безопасно обрабатывать одновременные операции записи и чтения из множества горутин.

### Требования
1. Реализовать структуру `SafeCache` с методами:
   - `Set(key string, value string)` — добавляет значение в кеш.
   - `Get(key string) (string, bool)` — возвращает значение по ключу.
2. Гарантировать отсутствие data race при параллельном доступе.




```go
package main

  

import (

    "sync"

)

  

type storage struct {

    caches map[string]string

    mu     sync.RWMutex

}

  

type SafeCache struct {

    key   string

    value string

}

  

func NewSafeCache() storage {

    return storage{caches: make(map[string]string), mu: sync.RWMutex{}}

}

  

func (s *storage) Set(keyInput string, valueInput string) {

    s.mu.Lock()

    s.caches[keyInput] = valueInput

    s.mu.Unlock()

}

  

func (s *storage) Get(keyInput string) (string, bool) {

    s.mu.RLock()

    result, ok := s.caches[keyInput]

    s.mu.RUnlock()

    if !ok {

        return "key is not valid", ok

    }

  

    return result, true

}
```
# Параллельная загрузка данных из нескольких источников
---
## Описание задачи
Реализовать систему параллельной загрузки данных из независимых источников:
1. Асинхронная загрузка комментариев из БД
2. Параллельная загрузка данных пользователей на основе полученных комментариев
3. Загрузка данных сессии и условная загрузка вложений

**Цель:**   
Освоить работу с горутинами, `sync.Once` и синхронизацией через `sync.WaitGroup`.

---
## Требования
1. Загрузка комментариев и данных сессии должна выполняться параллельно
2. Загрузка данных пользователей должна стартовать только после получения комментариев
3. Загрузка вложений должна выполняться только при наличии session-id
4. Использовать минимум 3 горутины для разных этапов
5. Синхронизировать все операции перед завершением
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://victoriametrics.com/blog/go-sync-once/

```go
package main

  

import (

    "fmt"

    "math/rand"

    "sync"

    "time"

)

  

var newSessionInfo sessionInfo

var once sync.Once

  

type sessionInfo struct {

    session

    amountSpectators []amountSpectatorsOnTime

    comments         []comment

    users            []User

}

  

type session struct {

    sessionID   int

    sessioName  string

    sessionTime time.Duration

}

  

type amountSpectatorsOnTime struct {

    time   time.Duration

    amount int

}

  

type comment struct {

    value string

    time  time.Duration

    User

}

  

type User struct {

    userID   int

    userName string

}

  

func NewSession(id int, name string, timesession time.Duration) *sessionInfo {

    return &sessionInfo{

        session: session{

            sessionID:   id,

            sessioName:  name,

            sessionTime: timesession * time.Second,

        },

        amountSpectators: make([]amountSpectatorsOnTime, 0),

        comments:         make([]comment, 0),

        users:            make([]User, 0),

    }

}

  

func (s *sessionInfo) SetSpectatorInfo(maxAmountSpectator int) {

  

    for i := 0; i < int(s.session.sessionTime.Seconds()); i++ {

        randAmount := rand.Intn(maxAmountSpectator) + 1

  

        s.amountSpectators = append(s.amountSpectators,

            amountSpectatorsOnTime{

                time:   time.Duration(i) * time.Second,

                amount: randAmount,

            },

        )

    }

}

  

func (s *sessionInfo) SetUsers(usersAmount int) {

    for i := 0; i < usersAmount; i++ {

        s.users = append(s.users, User{userName: "bot", userID: i + 1})

    }

}

  

func (s *sessionInfo) SetComment() {

    exampleComments := [7]string{

        "4:11.2.0-1ubuntu1",

        "sing /usr/bin/g++ to provide ",

        "randAmount",

        "spectators",

        "wsl -u root",

        "Hello world",

        "looozer",

    }

  

    len := int(s.session.sessionTime.Seconds())

    for i := 0; i <= len; {

        s.comments = append(s.comments,

            comment{value: exampleComments[rand.Intn(6)],

                time: time.Duration(i) * time.Second,

                User: User{

                    userID:   rand.Intn(1000),

                    userName: "bot",

                },

            },

        )

        i += rand.Intn(len - len/2)

    }

}

  

func LoadSession(stream sessionInfo) <-chan string {

    out := make(chan string)

  

    stream.SetSpectatorInfo(250)

  

    go func() {

        defer close(out)

        out <- fmt.Sprintf("id session: %d \nName session: %s\n", stream.sessionID,

            stream.sessioName,

        )

  

        for _, spectator := range stream.amountSpectators {

            out <- fmt.Sprintf("time : %s\n amount spectators: %b",

                spectator.time.String(),

                spectator.amount,

            )

        }

    }()

  

    return out

}

  

func LoadComments(stream sessionInfo, wg *sync.WaitGroup) <-chan string {

    out := make(chan string)

  

    stream.SetComment()

  

    if stream.sessionID == 0 {

        return nil

    }

  

    go func() {

        defer wg.Done()

  

        defer close(out)

        for _, com := range stream.comments {

            out <- fmt.Sprintf("time: %s\nuser name: %s\ncomment:\n \"%s\"\n\n",

                com.time.String(),

                com.User.userName,

                com.value,

            )

        }

    }()

    return out

}

  

func LoadUsers(stream sessionInfo) <-chan string {

    out := make(chan string)

  

    stream.SetUsers(310)

  

    go func() {

        defer close(out)

        for _, user := range stream.users {

            out <- fmt.Sprintf(" user ID: %d\nusername: %s\n",

                user.userID,

                user.userName,

            )

        }

    }()

    return out

}

  

func printChannels(channelArray []<-chan string) {

    for _, channel := range channelArray {

  

        go func(ch <-chan string) {

            for n := range ch {

                fmt.Println(n)

            }

        }(channel)

    }

}

  

func printUsers(in <-chan string, wg *sync.WaitGroup) {

    go func() {

        defer wg.Done()

        for n := range in {

            fmt.Println(n)

        }

    }()

}

  

func GetSession(sessionID int, nameSession string, time time.Duration) sessionInfo {

    once.Do(func() {

        newSessionInfo = NewSession(sessionID, nameSession, time)

    })

    return newSessionInfo

}

  

// 1. Загрузка комментариев и данных сессии должна выполняться параллельно

// 2. Загрузка данных пользователей должна стартовать только после получения комментариев

// 3. Загрузка вложений должна выполняться только при наличии session-id

// 4. Использовать минимум 3 горутины для разных этапов

// 5. Синхронизировать все операции перед завершением

func main() {

    wg := sync.WaitGroup{}

  

    GetSession(45, "Minecraft", 250)

  

    wg.Add(1)

    spectatorsInfo := LoadSession(newSessionInfo)

  

    commentsInfo := LoadComments(newSessionInfo, &wg)

  

    printChannels([]<-chan string{spectatorsInfo, commentsInfo})

    wg.Wait()

    users := LoadUsers(newSessionInfo)

    wg.Add(1)

    printUsers(users, &wg)

    wg.Wait()

}
```
-----------------------------------------

# SYNC/COND
## Реализация очереди с ограниченной емкостью на sync.Cond

### Описание задачи
В распределенных системах часто требуется синхронизировать работу продюсеров (добавляющих задачи) и консьюмеров (обрабатывающих задачи). Очередь с фиксированной емкостью (`BoundedQueue`) решает следующие проблемы:
- **Блокировка продюсеров** при заполнении очереди.
- **Блокировка консьюмеров** при опустошении очереди.
- **Потокобезопасность** в многогоруточной среде.
- **Корректное завершение** работы через `Shutdown()`.

**Цель:**  
Реализовать очередь, использующую `sync.Cond` для эффективной синхронизации горутин.

### Требования
1. Реализация методов:
    - `Put(task interface{})` — блокируется, если очередь заполнена.
    - `Get() interface{}` — блокируется, если очередь пуста.
    - `Shutdown()` — завершает работу очереди.
2. Использование `sync.Cond` и `sync.Mutex` для синхронизации.
3. Гарантия отсутствия гонок и утечек.
```go
package main

  

import (

    "fmt"

    "sync"

    "time"

)


  

type User struct {

    userID   int

    userName string

}

  

type BoundedQueue struct {

    queue []interface{}

  

    mutex sync.Mutex

    cond  *sync.Cond

  

    capacity int

  

    closed bool

}

  

func NewBoundedQueue(capacityQueue int) *BoundedQueue {

    q := &BoundedQueue{

        queue:    make([]interface{}, 0),

        capacity: capacityQueue,

    }

  

    q.cond = sync.NewCond(&q.mutex)

  

    return q

}

  

func (b *BoundedQueue) Put(task interface{}) {

    b.mutex.Lock()

    defer b.mutex.Unlock()

  

    for b.capacity == len(b.queue) {

        if b.closed {

            fmt.Println("Shutdown")

            return

        }

        fmt.Println("queue full producer waits")

        b.cond.Wait()

    }

  

    b.queue = append(b.queue, task)

  

    b.cond.Signal()

  

}

  

func (b *BoundedQueue) Get() interface{} {

    b.mutex.Lock()

    defer b.mutex.Unlock()

    for len(b.queue) == 0 {

        if b.closed {

            fmt.Println("Shutdown")

            return nil

        }

        fmt.Println("queue empty consumer waits")

        b.cond.Wait()

    }

  

    result := b.queue[0]

  

    b.queue = b.queue[1:]

  

    b.cond.Signal()

  

    return result

}

  

func (b *BoundedQueue) Shutdown() {

    b.mutex.Lock()

    defer b.mutex.Unlock()

    b.closed = true

  

    b.cond.Broadcast()

}

func main() {

  

    usernames := [5]string{

        "Igor",

        "Ruslan",

        "Jo",

        "lenar",

        "jey",

    }

    boundedQueue := NewBoundedQueue(4)

  

    go func() {

        for i, name := range usernames {

  

            boundedQueue.Put(User{userID: i, userName: name})

        }

    }()

  

    go func() {

        for i := 0; i < 8; i++ {

            fmt.Println(boundedQueue.Get())

        }

    }()

    time.Sleep(1 * time.Second)

    boundedQueue.Shutdown()

    time.Sleep(5 * time.Second)

}
```

---
## Моделирование работы ресторана с использованием `sync.Cond`

### Описание задачи
Реализовать систему управления столиками в ресторане, где:
- Количество столиков фиксировано (например, 5).
- Посетители (горутины) занимают столики, если они свободны.
- Если все столики заняты, посетители ожидают в очереди.
- При освобождении столика его получает первый ожидающий посетитель.

**Цель:**  
Научиться синхронизировать горутины с помощью `sync.Cond`, моделируя реальный сценарий с ограниченными ресурсами.


### Требования
1. Реализовать структуру `Restaurant` с методами:
    - `OccupyTable()` — блокируется, если нет свободных столиков.
    - `ReleaseTable()` — освобождает столик и уведомляет ожидающих.
2. Использовать `sync.Cond` для управления очередью ожидания.

```go
package main

  

import (

    "fmt"

    "sync"

    "time"

)

  

type visitor struct {

    Id    int

    value chan struct{}

}

type Restaurant struct {

    numberOfoccupiedTables int

  

    visitors []visitor

  

    mutex sync.Mutex

    cond  *sync.Cond

  

    numberOfTables int

  

    closed bool

}

  

func NewRestaurant(capacityQueue int) *Restaurant {

    q := &Restaurant{

        visitors:               make([]visitor, 0),

        numberOfoccupiedTables: 0,

        numberOfTables:         capacityQueue,

    }

  

    q.cond = sync.NewCond(&q.mutex)

  

    return q

}

  

func (r *Restaurant) OccupyTable() {

    r.mutex.Lock()

  

    if r.numberOfoccupiedTables >= r.numberOfTables {

        fmt.Println("no free tables")

        v := visitor{value: make(chan struct{})}

        r.visitors = append(r.visitors, v)

        r.mutex.Unlock()

        <-v.value

        return

    }

    r.numberOfoccupiedTables++

    r.mutex.Unlock()

    fmt.Println("1 table is busy")

  

}

  

func (r *Restaurant) ReleaseTable() {

    r.mutex.Lock()

  

    if len(r.visitors) == 0 {

        fmt.Println("no visitors")

        r.mutex.Unlock()

        return

    }

  

    r.numberOfoccupiedTables--

    fmt.Println("1 table is free")

    result := r.visitors[0]

  

    r.visitors = r.visitors[1:]

  

    r.mutex.Unlock()

    close(result.value)

  

}

  

func main() {

  

    restaurant := NewRestaurant(4)

  

    for i := 10; i > 0; i-- {

        go restaurant.OccupyTable()

    }

  

    for i := 0; i < 4; i++ {

        go restaurant.ReleaseTable()

    }

  

    time.Sleep(2 * time.Second)

  

    for i := 0; i < 4; i++ {

        go restaurant.ReleaseTable()

    }

  

    time.Sleep(2 * time.Second)

  

    for i := 0; i < 4; i++ {

        go restaurant.ReleaseTable()

    }

  

    time.Sleep(2 * time.Second)

}
```
---

## Пул подключений к БД с использованием `sync.Cond`

### Описание задачи
Реализовать пул подключений к базе данных с ограничением на максимальное количество активных подключений. Если все подключения заняты, новые запросы должны блокироваться до освобождения ресурсов. Использовать `sync.Cond` для синхронизации.

---

### Требования
1. Реализовать методы:
    - `Get() *Connection` — возвращает свободное подключение или блокирует горутину.
    - `Release(*Connection)` — освобождает подключение и уведомляет ожидающих.
2. Ограничить максимальное количество подключений (например, 3).
3. Гарантировать потокобезопасность.
4. Смоделировать работу с задержками (имитация запросов к БД).

---

```go
func main() {
    pool := NewConnectionPool(3) // Пул на 3 подключения

    for i := 0; i < 10; i++ {
        go func(id int) {
            conn := pool.Get()
            defer pool.Release(conn)

            fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
            time.Sleep(2 * time.Second) // Имитация работы
        }(i)
    }

    time.Sleep(10 * time.Second)
}

```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://ubiklab.net/posts/go-sync-cond/
https://dev.to/func25/go-synccond-the-most-overlooked-sync-mechanism-1fgd
https://wcademy.ru/go-multithreading-sync-cond/

```go
package main

  

import (

    "fmt"

    "math/rand/v2"

    "sync"

    "time"

)

  

type Connection struct {

    ID           int

    path         string

    valuestorage float64

}

  

func CreatePool() Connection {

    c := Connection{

        ID:           rand.IntN(100),

        path:         "./",

        valuestorage: rand.ExpFloat64(),

    }

  

    fmt.Println("create pool")

    time.Sleep(2 * time.Second)

    return c

}

  

type Pool struct {

    connections   []Connection

    maxConnection int

    mutex         sync.Mutex

    cond          *sync.Cond

    closed        bool

}

  

func NewConnectionPool(connectionAmount int) *Pool {

    q := &Pool{

        connections:   make([]Connection, 0),

        maxConnection: connectionAmount,

    }

  

    q.cond = sync.NewCond(&q.mutex)

  

    return q

}

  

func (r *Pool) Get() *Connection {

    r.mutex.Lock()

    defer r.mutex.Unlock()

    defer r.cond.Signal()

    for len(r.connections) >= r.maxConnection {

        fmt.Println("no free connections")

        r.cond.Wait()

    }

  

    c := CreatePool()

  

    r.connections = append(r.connections, c)

  

    fmt.Println("1 connection is busy")

  

    return &c

}

  

func (r *Pool) Release(conn *Connection) {

    r.mutex.Lock()

    defer r.mutex.Unlock()

    defer r.cond.Signal()

    amountConnections := len(r.connections)

  

    for amountConnections == 0 {

        fmt.Println("no current pool")

        r.cond.Wait()

    }

  

    for i := 0; i < amountConnections; i++ {

        if r.connections[i].ID == conn.ID {

  

            fmt.Println("release the next connection: ", r.connections[i])

            r.connections = append(r.connections[:i], r.connections[i+1:]...)

  

            return

        }

    }

    fmt.Println("connection is not found")

}

  

func main() {

    pool := NewConnectionPool(3) // Пул на 3 подключения

  

    for i := 0; i < 10; i++ {

        go func(id int) {

            conn := pool.Get()

            defer pool.Release(conn)

  

            fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)

            time.Sleep(2 * time.Second) // Имитация работы

        }(i)

    }

  

    time.Sleep(20 * time.Second)

}
```
-----------------------------------------

# SYNC/ONCE
## Задание: Ленивая инициализация подключения к базе данных с использованием `sync.Once`

**Цель задания**:  
Реализовать потокобезопасный механизм однократной инициализации подключения к базе данных с помощью `sync.Once`.

---

### Описание задачи

При работе с базами данных в многопоточной среде важно, чтобы подключение инициализировалось только один раз, даже если метод запроса подключения вызывается из нескольких горутин одновременно. В этом задании вам предстоит реализовать такую логику.

### Требования:
1. Создайте структуру `Database`, которая хранит подключение к БД (`conn`) и экземпляр `sync.Once`.
2. Реализуйте метод `GetConnection()`, который:
    - Инициализирует подключение к БД только при первом вызове.
    - Гарантирует, что последующие вызовы возвращают уже созданное подключение.
3. Убедитесь, что код потокобезопасен (нет гонок данных).


```go
package main

  

import (

    "fmt"

    "math/rand/v2"

    "sync"

    "time"

)

  

type Connection struct {

    ID           int

    path         string

    valuestorage float64

}

  

func CreatePool() Connection {

    c := Connection{

        ID:           rand.IntN(100),

        path:         "./",

        valuestorage: rand.ExpFloat64(),

    }

  

    fmt.Println("create pool")

    time.Sleep(2 * time.Second)

    return c

}

  

type Database struct {

    connections Connection

    onse        sync.Once

}

  

func (d *Database) GetConnection() Connection {

    d.onse.Do(func() {

        d.connections = CreatePool()

  

    })

    return d.connections

}

func main() {

    d := Database{}

  

    for i := 0; i < 10; i++ {

        conn := d.GetConnection()

        fmt.Println(conn)

    }

}
```
---
## Конфигуратор приложения с `sync.Once`

**Описание**
Этот проект реализует **потокобезопасный** менеджер конфигурации, который загружает настройки **только один раз** при первом запросе.  
Используется `sync.Once`, чтобы избежать повторной загрузки при одновременном доступе из нескольких горутин.
---

**Возможности**
1. Ленивая инициализация – загрузка конфигурации только при первом вызове.  
2. Потокобезопасность – отсутствие гонок данных при многопоточной работе.  
3. Гибкость – возможность загружать конфигурацию из файла, переменных окружения или базы данных.
---

**Реализованные методы**
- `LoadConfig()` – загружает конфигурацию **один раз** и сохраняет в памяти.
- `Get(key string) string` – возвращает значение конфигурации по ключу.
- `PrintConfig()` – выводит загруженные параметры.
---

```go
// Имитация загрузки конфигурации
cm.config = map[string]string{
"app_name":  "MyApp",
"port":      "8080",
"log_level": "debug",
}

func main() {
    keys := []string{"app_name", "port", "log_level"}
    configManager.PrintConfig()
}
```


```go
package main

  

import (

    "fmt"

    "sync"

)

  

type ConfigManager struct {

    config map[string]string

    id     int

    onse   sync.Once

}

  

func (cm *ConfigManager) GetConnection() map[string]string {

    cm.onse.Do(func() {

        cm.config = map[string]string{

            "app_name":  "MyApp",

            "port":      "8080",

            "log_level": "debug",

        }

  

    })

    return cm.config

}

  

func (cm *ConfigManager) Get(key string) string {

    return cm.config[key]

}

  

func (cm *ConfigManager) PrintConfig() {

    fmt.Println(cm.config)

}

func main() {

    keys := []string{"app_name", "port", "log_level"}

    configManage := ConfigManager{}

    configManage.GetConnection()

    for _, key := range keys {

        fmt.Println(configManage.Get(key))

    }

    configManage.PrintConfig()

}
```

---
## Инициализация плагинов с `sync.Once`

**Цель задания**
Реализовать систему безопасной инициализации плагинов, где:
- Каждый плагин инициализируется **только один раз**
- Инициализация потокобезопасна
- Ошибки при инициализации корректно обрабатываются
- Плагины доступны для использования из разных компонентов


**Требования**
1. **Структура `PluginManager`**:
    - Хранит загруженные плагины
    - Использует `sync.Once` для каждого плагина
    - Поддерживает конкурентный доступ

2. **Методы**:
    - `GetPlugin(name string) (Plugin, error)` – возвращает инициализированный плагин
    - `RegisterPlugin()` – регистрирует плагины (симуляция)
```golang
package main

import (
	"fmt"
	"log"
	"sync"
)

// Интерфейс для всех плагинов
type Plugin interface {
	Execute() string
}

// Управляет инициализацией и доступом к плагинам
type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

type pluginEntry struct {
	//Добавить необходимые поля для однократной инициализации
	initFn func() (Plugin, error)
}

// NewPluginManager создает новый менеджер плагинов
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

// RegisterPlugin регистрирует новый плагин
func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

// GetPlugin возвращает инициализированный плагин
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	// Реализовать:
	// 1. Проверку существования плагина
	// 2. Потокобезопасную однократную инициализацию
	// 3. Обработку и кэширование ошибок
	// 4. Возврат кэшированного результата
	return nil, fmt.Errorf("not implemented")
}

// DemoPlugin реализация плагина
type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	// Имитация длительной инициализации
	// time.Sleep(500 * time.Millisecond)
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)
	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	})

	var wg sync.WaitGroup

	// Тестирование рабочего плагина
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())
		}(i)
	}

	// Тестирование плагина с ошибкой
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}
```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://victoriametrics.com/blog/go-sync-once/
https://dev.to/jones_charles_ad50858dbc0/a-developers-guide-to-synconce-your-go-concurrency-lifesaver-3kf2
https://backendinterview.ru/goLang/concurrency/sync.html

```go
package main

  

import (

    "fmt"

    "log"

    "sync"

    "time"

)

  

// Интерфейс для всех плагинов

type Plugin interface {

    Execute() string

}

  

// Управляет инициализацией и доступом к плагинам

type PluginManager struct {

    plugins map[string]*pluginEntry

    mu      sync.RWMutex

}

  

type pluginEntry struct {

    once   sync.Once

    plugin Plugin

    err    error

    initFn func() (Plugin, error)

}

  

// NewPluginManager создает новый менеджер плагинов

func NewPluginManager() *PluginManager {

    return &PluginManager{

        plugins: make(map[string]*pluginEntry),

    }

}

  

// RegisterPlugin регистрирует новый плагин

func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {

    pm.mu.Lock()

    defer pm.mu.Unlock()

    pm.plugins[name] = &pluginEntry{

        initFn: initFn,

    }

}

  

// GetPlugin возвращает инициализированный плагин

func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {

    // Реализовать:

    var getCurrentPlugin *pluginEntry

    var ok bool

    pm.mu.RLock()

    if getCurrentPlugin, ok = pm.plugins[name]; !ok {

        return nil, fmt.Errorf("not emplimented")

    }

    pm.mu.RUnlock()

    getCurrentPlugin.once.Do(func() {

        getCurrentPlugin.plugin, getCurrentPlugin.err = getCurrentPlugin.initFn()

    })

  

    return getCurrentPlugin.plugin, getCurrentPlugin.err

}

  

// DemoPlugin реализация плагина

type DemoPlugin struct{}

  

func (p *DemoPlugin) Execute() string {

    return "DemoPlugin executed successfully!"

}

  

func initDemo() (Plugin, error) {

    // Имитация длительной инициализации

    time.Sleep(500 * time.Millisecond)

    return &DemoPlugin{}, nil

}

  

func main() {

    pm := NewPluginManager()

  

    pm.RegisterPlugin("demo", initDemo)

    pm.RegisterPlugin("broken", func() (Plugin, error) {

        return nil, fmt.Errorf("simulated error")

    })

  

    var wg sync.WaitGroup

  

    // Тестирование рабочего плагина

    for i := 0; i < 5; i++ {

        wg.Add(1)

        go func(id int) {

            defer wg.Done()

            p, err := pm.GetPlugin("demo")

            if err != nil {

                log.Printf("Goroutine %d error: %v", id, err)

                return

            }

            log.Printf("Goroutine %d: %s", id, p.Execute())

        }(i)

    }

  

    // Тестирование плагина с ошибкой

    for i := 5; i < 7; i++ {

        wg.Add(1)

        go func(id int) {

            defer wg.Done()

            _, err := pm.GetPlugin("broken")

            if err != nil {

                log.Printf("Goroutine %d error: %v", id, err)

            }

        }(i)

    }

  

    wg.Wait()

}
```

--------------------------------
# SYNC/POOL
## Оптимизация обработки строк с sync.Pool

### Описание задачи
В высоконагруженном сервисе частые аллокации буферов для преобразования строк создают нагрузку на GC.
Цель — реализовать оптимизированную функцию `ProcessString` с использованием `sync.Pool`, чтобы переиспользовать буферы `[]byte`.

### Требования
1. Функция `ProcessString(s string) string` преобразует строку в верхний регистр.
2. Использование `sync.Pool` для буферов `[]byte`.
3. Потокобезопасность, отсутствие утечек памяти.
```go
func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}
```

```go
package main

  

import (

    "fmt"

    "sync"

)

  

var bytePool = sync.Pool{

    New: func() any {

        bytes := make([]byte, 0, 1024)

        return &bytes

    },

}

  

func ProcessString(s string) string {

    buffer := bytePool.Get().(*[]byte)

    *buffer = (*buffer)[:0]

  

    *buffer = append(*buffer, s...)

    for i := range *buffer {

        if (*buffer)[i] >= 'a' && (*buffer)[i] <= 'z' {

            (*buffer)[i] -= 32

        }

    }

    result := string(*buffer)

  

    bytePool.Put(buffer)

    return result

}

  

func main() {

    examples := []string{

        "hello, world!",

        "gopher",

        "lorem ipsum dolor sit amet",

    }

  

    for _, s := range examples {

        processed := ProcessString(s)

        fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)

    }

}
```

## Оптимизация HTTP-обработчика с sync.Pool

### Описание задачи
В высоконагруженных сервисах, обрабатывающих тысячи HTTP-запросов в секунду, частая аллокация объектов для декодирования JSON становится узким местом. Каждый вызов `json.NewDecoder` создает новый экземпляр `RequestData`, что приводит к:
- Высокой нагрузке на GC (сборщик мусора).
- Увеличению времени обработки запросов.
- Нестабильной работе при пиковых нагрузках.

**Цель:**  
Использовать `sync.Pool` для переиспользования объектов `RequestData`, сократив аллокации и улучшив производительность.

---

### Требования
1. **Реализация пула объектов**
    - Создать пул для структур `RequestData` с предварительной инициализацией вложенных полей (например, `map` или `slice`).
    - Гарантировать потокобезопасность.

2. **Метод `Reset()`**
    - Очистить все поля объекта перед возвратом в пул.
    - Для слайсов: сохранить базовый массив (`items = items[:0]`).
    - Для мап: явно удалить все ключи.

3. **Отсутствие утечек данных**
    - Убедиться, что объекты из пула не сохраняют данные предыдущих запросов.

---
```go
func main() {
    http.HandleFunc("/", handleRequest)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
```

```go
package main

  

import (

    "encoding/json"

    "fmt"

    "net/http"

    "sync"

)

  

type RequestData struct {

    ID          int               `json:"user_id"`

    Name        string            `json:"user_name"`

    Tags        []string          `json:"tags"`

    Directories map[string]string `json:"directories"`

}

  

func (r *RequestData) Reset() {

  

    r.ID = 0

    r.Name = ""

    r.Tags = r.Tags[:0]

  

    for key := range r.Directories {

        delete(r.Directories, key)

    }

}

  

var requestDataPool = sync.Pool{

    New: func() any {

        return &RequestData{

            Tags:        make([]string, 0, 16),

            Directories: make(map[string]string),

        }

    },

}

  

func handleRequest(w http.ResponseWriter, req *http.Request) {

    data := requestDataPool.Get().(*RequestData)

    defer func() {

        data.Reset()

  

        requestDataPool.Put(data)

    }()

  

    err := json.NewDecoder(req.Body).Decode(data)

    if err != nil {

        http.Error(w, err.Error(), http.StatusBadRequest)

        return

    }

    fmt.Printf("UserID: %d\n", data.ID)

    fmt.Printf("Name: %s\n", data.Name)

    fmt.Printf("Tags: %+v\n", data.Tags)

    fmt.Printf("Meta: %+v\n", data.Directories)

    w.WriteHeader(http.StatusOK)

    w.Write([]byte("ok"))

}

  

func main() {

    http.HandleFunc("/", handleRequest)

    fmt.Println("Server started at :8080")

    http.ListenAndServe(":8080", nil)

}
```

## JSON-кэш с `sync.Pool` и `map + RWMutex`

### **Описание**
Этот проект демонстрирует **потокобезопасный JSON-кэш** с поддержкой TTL и оптимизированной сериализацией.  
Используется `sync.Pool` для **эффективной работы с JSON**, а также `map + RWMutex` для **более быстрого доступа к данным**.

---

### **Основные возможности**
- **Хранение объектов в `map` (с TTL)**
- **Автоматическое удаление устаревших объектов**
- **Быстрая сериализация JSON с `sync.Pool`**
- **Использование `sync.RWMutex` для конкурентного доступа**

---

### **Методы**
#### **Базовые операции**
- `Set(key string, value interface{})` – **добавить объект в кэш**
- `Get(key string) (interface{}, bool)` – **получить объект по ключу**
- `Delete(key string)` – **удалить объект**
- `ToJSON() ([]byte, error)` – **сериализовать кэш в JSON**

---

### **Как это работает?**
- Все объекты хранятся в **`map[string]item`** (ключ → объект с TTL).
- `sync.Pool` позволяет **переиспользовать JSON-буферы**, снижая нагрузку на GC.
- Очистка устаревших данных выполняется **в отдельной горутине**.  

```go
func main() {
    cache := NewObjectCache(5 * time.Second)

    // Добавляем данные в кэш
    cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
    cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

    // Получаем объект
    if user, found := cache.Get("user:1"); found {
        fmt.Println("Найден:", user)
    }

    // Выводим JSON
    jsonData, _ := cache.ToJSON()
    fmt.Println("Кэш в JSON:", string(jsonData))

    // Ждём истечения TTL и проверяем снова
    time.Sleep(6 * time.Second)
    _, found := cache.Get("user:1")
    fmt.Println("После TTL, user:1 найден?", found)
}
```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://ubiklab.net/posts/go-pool-and-mechanics-behind-it/
https://reliasoftware.com/blog/golang-sync-pool
https://dev.to/func25/go-syncpool-and-the-mechanics-behind-it-52c1
https://engineer.yadro.com/article/three-ways-to-optimize-memory-performance-on-go-with-memory-pools/
https://leapcell.io/blog/boost-go-performance-sync-pool
https://www.sobyte.net/post/2022-06/go-sync-pool/
https://goperf.dev/01-common-patterns/object-pooling/

```go
package main

  

import (

    "bytes"

    "encoding/json"

    "fmt"

    "sync"

    "time"

)

  

type item struct {

    value     interface{}

    expiresAt time.Time

}

  

type ObjectCache struct {

    mu      sync.RWMutex

    items   map[string]item

    ttl     time.Duration

    closed  chan struct{}

    bufPool sync.Pool

}

  

func NewObjectCache(ttl time.Duration) *ObjectCache {

    c := &ObjectCache{

        items:  make(map[string]item),

        ttl:    ttl,

        closed: make(chan struct{}),

    }

  

    c.bufPool.New = func() any {

        return new(bytes.Buffer)

    }

  

    go c.startGC()

  

    return c

}

  

func (c *ObjectCache) Set(key string, value interface{}) {

    c.mu.Lock()

    defer c.mu.Unlock()

  

    c.items[key] = item{

        value:     value,

        expiresAt: time.Now().Add(c.ttl),

    }

}

  

func (c *ObjectCache) Get(key string) (interface{}, bool) {

    c.mu.RLock()

    it, ok := c.items[key]

    c.mu.RUnlock()

  

    if !ok {

        return nil, false

    }

  

    if time.Now().After(it.expiresAt) {

        c.mu.Lock()

        delete(c.items, key)

        c.mu.Unlock()

        return nil, false

    }

  

    return it.value, true

}

  

func (c *ObjectCache) Delete(key string) {

    c.mu.Lock()

    delete(c.items, key)

    c.mu.Unlock()

}

  

func (c *ObjectCache) ToJSON() ([]byte, error) {

    buf := c.bufPool.Get().(*bytes.Buffer)

    buf.Reset()

  

    c.mu.RLock()

    defer c.mu.RUnlock()

  

    enc := json.NewEncoder(buf)

    err := enc.Encode(c.items)

  

    data := make([]byte, buf.Len())

    copy(data, buf.Bytes())

  

    c.bufPool.Put(buf)

  

    return data, err

}

  

func (c *ObjectCache) startGC() {

    ticker := time.NewTicker(time.Second)

    defer ticker.Stop()

  

    for {

        select {

        case <-ticker.C:

            c.cleanup()

        case <-c.closed:

            return

        }

    }

}

  

func (c *ObjectCache) cleanup() {

    now := time.Now()

  

    c.mu.Lock()

    defer c.mu.Unlock()

  

    for k, v := range c.items {

        if now.After(v.expiresAt) {

            delete(c.items, k)

        }

    }

}

  

func (c *ObjectCache) Close() {

    close(c.closed)

}

  

func main() {

    cache := NewObjectCache(5 * time.Second)

  

    // Добавляем данные в кэш

    cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})

    cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

  

    // Получаем объект

    if user, found := cache.Get("user:1"); found {

        fmt.Println("Найден:", user)

    }

  

    // Выводим JSON

    jsonData, _ := cache.ToJSON()

    fmt.Println("Кэш в JSON:", string(jsonData))

  

    // Ждём истечения TTL и проверяем снова

    time.Sleep(6 * time.Second)

    _, found := cache.Get("user:1")

    fmt.Println("После TTL, user:1 найден?", found)

}
```
--------------------

# SYNC/WAIT

# Задание: Синхронизация горутин с использованием `sync.WaitGroup`

**Цель задания**:  
Исправить код, чтобы все горутины корректно выводили значения от 0 до 99, и обеспечить завершение всех горутин перед выходом из программы.
```go
package main

import "fmt"

func main() {
	cnt := 100
	for i := 0; i < cnt; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
```


## Задание: Параллельные HTTP-запросы с синхронизацией через `sync.WaitGroup`

**Цель задания**:  
Исправить код так, чтобы основная горутина дожидалась завершения всех HTTP-запросов.
```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchUrl(url string) error {
	_, err := http.Get(url)
	return err
}
func main() {
	urls := []string{
		"https://www.lamoda.ru",
		"https://www.yandex.ru",
		"https://www.mail.ru",
		"https://www.google.ru",
	}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s....\n", url)
			err := fetchUrl(url)
			if err != nil {
				fmt.Printf("Error feaching %s: %v\n", url, err)
				return
			}
			fmt.Printf("Fetched %s\n", url)
		}(url)
	}
	fmt.Println("All request launched!")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Program finished")
}

```


# Задание: Исправление синхронизации горутин 

**Цель задания**:  
Исправить код, чтобы основная горутина дожидалась завершения.
```go
package main

import (
	"context"
	"fmt"
	"time"
)

type logic struct{}

var Logic logic

func (l *logic) UpdateDB(ctx context.Context, item *Item) error {
	return nil // Заглушка
}

func (l *logic) FetchItems(ctx context.Context) ([]*Item, error) {
	return []*Item{
		{Value: 5},
		{Value: 15},
		{Value: 7},
	}, nil // Заглушка
}

type Item struct {
	Value int
}

func processItem(item *Item) {
	time.Sleep(time.Second)
	if item.Value > 10 {
		fmt.Printf("ERROR: item %d can't be more than 10\n", item.Value)
		return
	}

	err := Logic.UpdateDB(context.Background(), item)
	if err != nil {
		fmt.Println("ERROR: can't process item")
	}
}

func DoBusinessLogic() error {
	items, err := Logic.FetchItems(context.Background())
	if err != nil {
		return err
	}

	for _, item := range items {
		go processItem(item)
	}

	return nil
}

func main() {
	err := DoBusinessLogic()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("All items processed")
}
```
https://dev.to/jpoly1219/waitgroups-in-go-3dkj
https://bytegoblin.io/blog/a-beginners-guide-to-waitgroups-in-go.mdx
https://www.golinuxcloud.com/golang-waitgroup/
https://habr.com/ru/articles/850018/