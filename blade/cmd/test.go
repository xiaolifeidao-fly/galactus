package main

import (
	"fmt"
	"galactus/blade/internal/user/dao/po"
	"galactus/blade/internal/user/dao/repository"
	"galactus/common/middleware/db"
)

func demo() {
	Init()
	userRepository := db.NewRepository[repository.UserRepository]()
	userPo := &po.User{
		Username: "22dsda",
		Password: "test2",
		Nickname: "test3",
		Sex:      "test1",
		Channel:  "test5222",
		Mobile:   "test1",
	}
	userPo.Id = 2
	r, err := userRepository.SaveOrUpdate(userPo)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
	user, err := userRepository.FindById(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}
}

func main() {
	demo()
}
