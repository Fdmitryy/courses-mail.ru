package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var kitchen room
var corridor room
var myRoom room
var street room
var Player state

type room struct {
	Name      string
	Items     map[string][]string
	CanGo     []string
	ToSayLook string
	ToSayGo   string
	ToDo      string
	ToApply   []string
}

type state struct {
	WhatRoom  room
	Inventory map[string][]string
	Door      bool
	Commands  map[string]func(args []string) string
}

func main() {
	initGame()
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		text := in.Text()
		fmt.Println(handleCommand(text))
	}
}

func initGame() {
	kitchen = room{
		Name:      "кухня",
		Items:     map[string][]string{"на столе: ": {"чай"}},
		CanGo:     []string{"коридор"},
		ToSayLook: "ты находишься на кухне",
		ToSayGo:   "кухня, ничего интересного",
		ToDo:      ", надо собрать рюкзак и идти в универ",
		ToApply:   []string{},
	}
	corridor = room{
		Name:      "коридор",
		Items:     map[string][]string{},
		CanGo:     []string{"кухня", "комната", "улица"},
		ToSayLook: "",
		ToSayGo:   "ничего интересного",
		ToDo:      "",
		ToApply:   []string{"дверь"},
	}
	myRoom = room{
		Name:      "своя комната",
		Items:     map[string][]string{"на столе: ": {"ключи", "конспекты"}, "на стуле: ": {"рюкзак"}},
		CanGo:     []string{"коридор"},
		ToSayLook: "",
		ToSayGo:   "ты в своей комнате",
		ToDo:      "",
		ToApply:   []string{},
	}
	street = room{
		Name:      "улица",
		Items:     map[string][]string{},
		CanGo:     []string{"домой"},
		ToSayLook: "",
		ToSayGo:   "на улице весна",
		ToDo:      "",
		ToApply:   []string{},
	}
	Player = state{
		WhatRoom:  kitchen,
		Inventory: map[string][]string{},
		Door:      false,
		Commands: map[string]func(args []string) string{
			"осмотреться": lookAround,
			"идти":        goTo,
			"надеть":      putOn,
			"взять":       take,
			"применить":   apply,
		},
	}
}

func handleCommand(command string) string {
	str := strings.Split(command, " ")
	todo := str[0]
	if _, exist := Player.Commands[todo]; !exist {
		return "неизвестная команда"
	}
	var args []string
	for i := 1; i < len(str); i++ {
		args = append(args, str[i])
	}
	result := Player.Commands[todo](args)
	return result
}

func lookAround(args []string) string {
	result := Player.WhatRoom.ToSayLook
	args = append(args, result)
	if len(Player.WhatRoom.Items) != 0 {
		names := make([]string, 0, len(Player.WhatRoom.Items))
		for name := range Player.WhatRoom.Items {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			if len(result) != 0 {
				result += ", "
			}
			result += name
			for _, v := range Player.WhatRoom.Items[name] {
				result = result + v + ", "
			}
			result = result[:len(result)-2]
		}
	} else {
		result += "пустая комната"
	}
	if _, exist := Player.Inventory["рюкзак"]; exist && Player.WhatRoom.Name == "кухня" {
		Player.WhatRoom.ToDo = ", надо идти в универ"
	}
	result += Player.WhatRoom.ToDo
	result += ". можно пройти - "
	for _, v := range Player.WhatRoom.CanGo {
		result += v
	}
	return result
}

func goTo(args []string) string {
	var rooms = map[string]room{
		"кухня":   kitchen,
		"коридор": corridor,
		"комната": myRoom,
		"улица":   street,
	}
	var flag bool
	for _, v := range Player.WhatRoom.CanGo {
		if args[0] == v {
			flag = true
		}
	}
	if !flag {
		return "нет пути в " + args[0]
	}
	if !Player.Door && args[0] == "улица" {
		return "дверь закрыта"
	}
	Player.WhatRoom = rooms[args[0]]
	result := Player.WhatRoom.ToSayGo
	result += ". можно пройти - "
	for _, v := range Player.WhatRoom.CanGo {
		result = result + v + ", "
	}
	result = result[:len(result)-2]
	return result
}

func putOn(args []string) string {
	if findItem(args[0], Player.WhatRoom.Items) {
		Player.Inventory[args[0]] = append(Player.Inventory[args[0]], "")
		delItem(args[0])
		return "вы надели: " + args[0]
	}
	return "вам нечего надеть))"
}

func take(args []string) string {
	if len(Player.Inventory) == 0 {
		return "некуда класть"
	}
	if findItem(args[0], Player.WhatRoom.Items) {
		Player.Inventory["рюкзак"] = append(Player.Inventory["рюкзак"], args[0])
		delItem(args[0])
		return "предмет добавлен в инвентарь: " + args[0]
	}
	return "нет такого"
}

func apply(args []string) string {
	if findItem(args[0], Player.Inventory) {
		for _, v := range Player.WhatRoom.ToApply {
			if v == args[1] {
				Player.Door = !Player.Door
				return args[1] + " открыта"
			}
		}
		return "не к чему применить"
	}
	return "нет предмета в инвентаре - " + args[0]
}

func delItem(arg string) {
	for key, value := range Player.WhatRoom.Items {
		for i, v := range value {
			if arg == v {
				if len(value) == 1 {
					delete(Player.WhatRoom.Items, key)
					return

				}
				value = append(value[:i], value[i+1:]...)
				Player.WhatRoom.Items[key] = value
				return
			}
		}
	}
}

func findItem(arg string, hash map[string][]string) bool {
	for _, value := range hash {
		for _, val := range value {
			if val == arg {
				return true
			}
		}
	}
	return false
}
