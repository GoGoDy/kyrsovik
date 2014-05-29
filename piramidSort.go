package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"bufio"
	"os"
	"strconv"
	"strings"
)
// тут переменная для имени файла (для такой структуры загрузка доступня через файл)
var (
	file string
)
//определяем стрктуру данных
type customer struct {
	Fio     string  `json:"fio"`
	Account string  `json:"acc"`
	Money   float64 `json:"money"`
	Time    int64   `json:"time"` // к примеру в днях
}
//массив структур
type customers struct {
	Csmrs []customer `json:"data"`
}
// дальше методы которые характерны для типа customers
// все методы однобразны используют пирамидальную сортировку и поиск половинным делением
//поэтому опису одну функцию остальные работают также
//поиск по фамилии
func (c *customers) FindFio(fio string) (bool, customer) {
	var (
		exist bool
		cstm  customer
	)
	//сортируем массив
	c.piramidSort("fio")
	//создаем временную переменнюу
	temp := c.Csmrs
//бесконечный цикл для поиска
	for {
//		смотрим длину массива
		lenght := len(temp)
//		сравниваем искомое значение со средним элементом
//		и в зависимости от результата создаем новый массив либо из левой, либо из правой части
		if temp[int(lenght/2)].Fio > fio {
			temp = temp[:int(lenght/2)]
		} else {
			temp = temp[int(lenght/2):]
		}
//		когда остался один элемент присваеваем его временной переменной и указываем искомый он или нет
		if lenght == 1 {
			exist = temp[0].Fio == fio
			cstm = temp[0]
			break
		}
	}
//	возвращаем значение
	return exist, cstm
}

func (c *customers) FindAccount(acc string) (bool, customer) {
	var (
		exist bool
		cstm  customer
	)
	c.piramidSort("acc")
	temp := c.Csmrs

	for {
		lenght := len(temp)
		if temp[int(lenght/2)].Account > acc {
			temp = temp[:int(lenght/2)]
		} else {
			temp = temp[int(lenght/2):]
		}
		if lenght == 1 {
			exist = temp[0].Account == acc
			cstm = temp[0]
			break
		}
	}
	return exist, cstm
}

func (c *customers) FindMoney(mn float64) (bool, customer) {
	var (
		exist bool
		cstm  customer
	)
	*c = c.piramidSort("money")
	temp := c.Csmrs

	for {
		lenght := len(temp)
		if temp[int(lenght/2)].Money > mn {
			temp = temp[:int(lenght/2)]
		} else {
			temp = temp[int(lenght/2):]
		}
		if lenght == 1 {
			exist = temp[0].Money == mn
			cstm = temp[0]
			break
		}
	}
	return exist, cstm
}

func (c *customers) FindTime(tm int64) (bool, customer) {
	var (
		exist bool
		cstm  customer
	)
	c.piramidSort("time")
	temp := c.Csmrs

	for {
		lenght := len(temp)
		if temp[int(lenght/2)].Time > tm {
			temp = temp[:int(lenght/2)]
		} else {
			temp = temp[int(lenght/2):]
		}
		if lenght == 1 {
			exist = temp[0].Time == tm
			cstm = temp[0]
			break
		}
	}
	return exist, cstm
}
// пирамидальная сортировка
func (c customers) piramidSort(field string) customers {
	var (
		element, child, length int
		sravnenie              bool
	)
//	смотрим длину массива
	length = len(c.Csmrs)
//	ставим ключ на первый элемент
	element = 0
	for {
		flagToNextPart := true

		child = 2*element + 1
//		проверяем существует ли элемент 2*k+1
		if child <= (length - 1) {
//			если да
//			проверяем больше он или нет от элемента к (для каждого поля структуры свое сравнение)
			switch field {
			case "fio":
				sravnenie = c.Csmrs[element].Fio < c.Csmrs[child].Fio
			case "acc":
				sravnenie = c.Csmrs[element].Account < c.Csmrs[child].Account
			case "money":
				sravnenie = c.Csmrs[element].Money < c.Csmrs[child].Money
			case "time":
				sravnenie = c.Csmrs[element].Time < c.Csmrs[child].Time
			}
			if sravnenie {
//				если он больше то меняем местами элементы
				c.Csmrs[element], c.Csmrs[child] = c.Csmrs[child], c.Csmrs[element]
//				ставим ключ на 0 чтобы начать сначала массива
				element = 0
//				устанавливаем флаг для того чтобы значение elemet не наращивалось
				flagToNextPart = false
				continue;
			}
		}
//		аналогисно и для элемента 2*k+2
		child = 2*element + 2
		if child <= (length - 1) {
			switch field {
			case "fio":
				sravnenie = c.Csmrs[element].Fio < c.Csmrs[child].Fio
			case "acc":
				sravnenie = c.Csmrs[element].Account < c.Csmrs[child].Account
			case "money":
				sravnenie = c.Csmrs[element].Money < c.Csmrs[child].Money
			case "time":
				sravnenie = c.Csmrs[element].Time < c.Csmrs[child].Time
			}
			if sravnenie {
				c.Csmrs[element], c.Csmrs[child] = c.Csmrs[child], c.Csmrs[element]
				element = 0
				flagToNextPart = false
				continue;
			}
		}
//		если изменений не было
		if flagToNextPart {
//			переходим к след элементы массива
			element++
//			если индекс текущего элемента больше среднего индекса
			if element > (length / 2) {
//				перенести первое значение (оно самое большое) в конец списка
				c.Csmrs[0], c.Csmrs[length-1] = c.Csmrs[length-1], c.Csmrs[0]
//				уменьшить массив на один элемент
				length = length - 1
				element = 0
			}
//			повторить все заново
		}
//		если длина расна 0 выйти
		if length == 0 {
			break
		}
	}
//	вернуть отсортированный массив
	return c
}

func init() {
//	считать флаги запуска
	flag.StringVar(&file, "f", "", "File")
	flag.Parse()
}

func main() {
	var (
		cst  customers
		c    customer
		ex   bool
		err  error
		text []byte
		command string
		flagEndofLoop bool     = true
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
		fio, acc string
		money float64
		time int64
		fioArr []string
	)
//	читаем файл
	text, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
//преобразовываем файл в массив структур
	err = json.Unmarshal(text, &cst)
	if err != nil {
		fmt.Println(err)
	}

	for flagEndofLoop{
//		пишем команду
		fmt.Println("Введите команду и параметр через пробел (для более подробной информации наберите help)")
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		command = scanner.Text()

		switch command{
		case "fio":
//			ищем по фамилии
			scanner.Split(bufio.ScanWords)

			scanner.Scan()
			v := scanner.Text()
			fmt.Printf("%v\n", v)
			fioArr = append(fioArr, string(v))

			scanner.Scan()
			v = scanner.Text()
			fmt.Printf("%v\n", v)
			fioArr = append(fioArr, string(v))


			fio = strings.Join(fioArr, " ")
			fio = strings.Trim(fio, " ")
			ex, c = cst.FindFio(fio);
			if(ex){
				fmt.Printf("%v\n", c )
			}else{
				fmt.Println("нет такого")
			}

		case "acc":
//			ищем по счету
			scanner.Split(bufio.ScanWords)
			scanner.Scan()
			acc = scanner.Text();
			ex, c = cst.FindAccount(acc);
			if(ex){
				fmt.Printf("%v\n", c )
			}else{
				fmt.Println("нет такого")
			}

		case "money":
//			ищем по деньгам
			scanner.Split(bufio.ScanWords)
			scanner.Scan()
			money, _ = strconv.ParseFloat(scanner.Text(), 32);
			ex, c = cst.FindMoney(money);
			if(ex){
				fmt.Printf("%v\n", c )
			}else{
				fmt.Println("нет такого")
			}

		case "time":
//			ищем по сроку
			scanner.Split(bufio.ScanWords)
			scanner.Scan()
			time, _ = strconv.ParseInt(scanner.Text(), 10, 0);
			ex, c = cst.FindTime(time);
			if(ex){
				fmt.Printf("%v\n", c )
			}else{
				fmt.Println("нет такого")
			}

		case "help":
//			помощь по командам
			fmt.Println("Возможные команды:")
			fmt.Println("fio строка : поиск по фамилии и имени")
			fmt.Println("acc строка: поиск по счету")
			fmt.Println("money число: поиск по кол-ву денег")
			fmt.Println("time число:  поиск по сроку  ")
			fmt.Println("end: Выход из программы")

		case "end":
//			если надоело то выходим
			fmt.Println("пока")
			flagEndofLoop = false;
		}



	}
}
