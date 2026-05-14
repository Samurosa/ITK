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
	comments         []comment
	users            []User
}

type session struct {
	sessionID   int
	sessioName  string
	sessionTime time.Duration
}

type amountSpectatorsOnTime struct {
	time   time.Duration
	amount int
}

type comment struct {
	value string
	time  time.Duration
	User
}

type User struct {
	userID   int
	userName string
}

func NewSession(id int, name string, timesession time.Duration) sessionInfo {
	return sessionInfo{
		session: session{
			sessionID:   id,
			sessioName:  name,
			sessionTime: timesession * time.Second,
		},
		amountSpectators: make([]amountSpectatorsOnTime, 0),
		comments:         make([]comment, 0),
		users:            make([]User, 0),
	}
}

func (s *sessionInfo) SetSpectatorInfo(maxAmountSpectator int) {

	for i := 0; i < int(s.session.sessionTime.Seconds()); i++ {
		randAmount := rand.Intn(maxAmountSpectator) + 1

		s.amountSpectators = append(s.amountSpectators,
			amountSpectatorsOnTime{
				time:   time.Duration(i) * time.Second,
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
					userID:   rand.Intn(1000),
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
