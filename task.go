package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type clock struct {
	open  time.Time // время начала работы
	close time.Time // время окончания работы
}

type event struct {
	time   time.Time //время события
	id     int       //id события
	client string    //имя клиента
}

type table struct { // стол
	id        int
	sum       int
	timecheck time.Time
	duration  time.Duration
	isbusy    bool
	client    string
}

var (
	n       int      // количество столов в компьютерном клубе
	cost    int      // стоимость часа в компьютерном клубе
	num     int      // номер стола за который сядет клиент/уйдет клиент
	work    clock    // время работы
	ev      []event  // события
	clients []string // клиенты в ожидании
	tables  []table  // состояние столов
)

func ParseHHMM(form string, t time.Time) time.Time {
	tt := strings.Split(form, ":")
	h1, _ := strconv.Atoi(tt[0])
	m1, _ := strconv.Atoi(tt[1])
	return time.Date(t.Year(), t.Month(), t.Day(), h1, m1, 0, t.Nanosecond(), t.Location())
}

func findIndex(sl []string, name string) int {
	// iterate over the array and compare given string to each element
	for index, value := range sl {
		if value == name {
			return index
		}
	}
	return -1
}

func findIndexByTable(sl []table, name string) int {
	// iterate over the array and compare given string to each element
	for index, value := range sl {
		if value.client == name {
			return index
		}
	}
	return -1
}

func isfree() bool {
	for _, value := range tables {
		if value.isbusy == false {
			return true
		}
	}
	return false
}

func check(e event) (event, bool) {
	var ev event
	ev.time = e.time
	if e.id == 1 {
		if e.time.Before(work.open) || e.time.After(work.close) {
			ev.id = 13
			ev.client = "NotOpenYet"
			goto exit
		} else if findIndex(clients, e.client) >= 0 {
			ev.id = 13
			ev.client = "YouShallNotPass"
			goto exit
		}
		clients = append(clients, e.client)
	} else if e.id == 2 {
		if tables[e.id].isbusy {
			ev.id = 13
			ev.client = "PlaceIsBusy"
			goto exit
		} else if findIndex(clients, e.client) == -1 && findIndexByTable(tables, e.client) == -1 {
			ev.id = 13
			ev.client = "ClientUnknown"
			goto exit
		}
		index := findIndexByTable(tables, e.client)
		if index >= 0 {
			table := &tables[num]
			table.client = tables[index].client
			table.isbusy = true
			table.timecheck = e.time

			table2 := &tables[index]
			table2.client = ""
			table2.isbusy = false
			hour := e.time.Sub(table2.timecheck)
			table2.duration += hour
			table2.sum += (int(hour.Hours()) + 1) * cost
		}
		index = findIndex(clients, e.client)
		if index >= 0 {
			table := &tables[num]
			table.client = clients[index]
			table.isbusy = true
			table.timecheck = e.time
			clients = append(clients[:index], clients[index+1:]...)
		}

	} else if e.id == 3 {
		if isfree() {
			ev.id = 13
			ev.client = "ICanWaitNoLonger"
			goto exit
		} else if len(clients) > n {
			ev.id = 11
			ev.client = e.client
			index := findIndex(clients, e.client)
			clients = append(clients[:index], clients[index+1:]...)
			goto exit
		}
	} else if e.id == 4 {
		if findIndex(clients, e.client) == -1 && findIndexByTable(tables, e.client) == -1 {
			ev.id = 13
			ev.client = "ClientUnknown"
			goto exit
		}
		index := findIndexByTable(tables, e.client)
		if index >= 0 {
			// ev.id = 12
			table := &tables[index]
			table.isbusy = false
			hour := e.time.Sub(table.timecheck)
			table.duration += hour
			table.sum += (int(hour.Hours()) + 1) * cost
			table.client = ""
		}

		if len(clients) > 0 {
			ev.id = 12
			index = 0
			table := &tables[index]
			table.client = clients[index]
			ev.client = clients[index]
			table.isbusy = true
			table.timecheck = e.time
			clients = append(clients[:index], clients[index+1:]...)
			goto exit
		}
	}
	return ev, false
exit:
	return ev, true

}

func main() {
	argsWithoutProg := os.Args
	f, err := os.Open(argsWithoutProg[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// scanner.Split(bufio.ScanWords)

	if scanner.Scan() {
		s := scanner.Text()
		n, err = strconv.Atoi(s) // количество столов в компьютерном клубе
	}
	tables = make([]table, n)

	t := time.Now()
	if scanner.Scan() {
		s := scanner.Text()
		tt := strings.Split(s, " ")
		work.open = ParseHHMM(tt[0], t)  // время начала работы
		work.close = ParseHHMM(tt[1], t) // время окончания работы
		fmt.Println(work.open.Format("15:04"))
	}
	if scanner.Scan() {
		s := scanner.Text() // считываем количество столов в компьютерном клубе
		cost, _ = strconv.Atoi(s)
	}

	for scanner.Scan() {
		// do something with a word
		s := scanner.Text()
		matched, _ := regexp.MatchString(`^([0-2]\d:\d{2}) [1-4] ([a-zA-Z]+[0-9]+_??-??)([\s][\d+])?$`, s)
		if matched {
			fmt.Println(s)
		} else {
			break
		}
		var e event
		matched, _ = regexp.MatchString(`^([0-2]\d:\d{2}) [134] ([a-zA-Z]+[0-9]+_??-??)$`, s)
		if matched {
			tt := strings.Split(s, " ")
			e.time = ParseHHMM(tt[0], t)
			e.id, _ = strconv.Atoi(tt[1])
			e.client = tt[2]
			ev = append(ev, e)
		}
		matched, _ = regexp.MatchString(`^([0-2]\d:\d{2}) 2 ([a-zA-Z]+[0-9]+_??-??) (\d+)$`, s)
		if matched {
			// var e event
			tt := strings.Split(s, " ")
			e.time = ParseHHMM(tt[0], t)
			e.id, _ = strconv.Atoi(tt[1])
			e.client = tt[2]
			ev = append(ev, e)
			num, _ = strconv.Atoi(tt[3])
			num--
		}

		res, is := check(e)
		if is {
			if res.id != 12 {
				fmt.Println(res.time.Format("15:04"), res.id, res.client)
			} else {
				fmt.Println(res.time.Format("15:04"), res.id, res.client, num)
			}
		}

	}
	// конец рабочего дня
	for index, value := range tables {
		if value.client != "" {
			table := &tables[index]
			clients = append(clients, table.client)
			table.client = ""
			table.isbusy = false
			hour := work.close.Sub(table.timecheck)
			table.duration += hour
			table.sum += (int(hour.Hours()) + 1) * cost
		}
	}
	sort.Sort(sort.StringSlice(clients))
	for _, value := range clients {
		fmt.Println(work.close.Format("15:04"), 11, value)
	}
	fmt.Println(work.close.Format("15:04"))

	// вывод прибыли
	for index, value := range tables {
		// fmt.Println()
		duration := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
		duration = duration.Add(value.duration)
		fmt.Println(index+1, value.sum, duration.Format("15:04"))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
