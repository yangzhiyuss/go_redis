package main
// interface 的实现

import (
	"fmt"
	"math/rand"
	"time"
)

type IUser interface {
	SetName(string)
	GetName() string
}

type User struct {
	name string
}


func (self User) GetName() string {
	return self.name
}

func (self *User) SetName(name string) {
	self.name = name
}

func main() {
	var i IUser
	i = &User{}
	i.SetName("wxnacy")
	fmt.Println(i.GetName())
	// wxnacy

	for true {
		rand.Seed(time.Now().UnixNano())
		fmt.Println(rand.Int())
		time.Sleep(10 *time.Second)
	}

}