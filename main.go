package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	fieldsNum = 10 // количество ячеек
	minesNum  = 10 // количество мин
)

var field [fieldsNum][fieldsNum]string // массив отображаемого поля
var opened [fieldsNum][fieldsNum]bool  // массив открытых ячеек
var mines [fieldsNum][fieldsNum]bool   // массив мин

func main() {
	rand.Seed(time.Now().Unix()) // засеиваем генератор случайных чисел (не знаю делается ли это в C++)
	setMines()                   // расставляем мины
	start()                      // готовим пустое поле
	show()                       // отображаем поле
	for {

		// проверяем победу (все поля, кроме мин открыты)
		if check() {
			println("\nYou win!")
			os.Exit(0)
		}

		var i, j int
		// читаем координаты выстрела (если не валидные - повторяем ввод)
		for {
			getInput(&i, &j)
			if i < fieldsNum && j < fieldsNum && i >= 0 && j >= 0 { // проверка на выход за пределы массива
				break
			}
		}

		hit := shut(i, j) // выстрел
		if hit {
			boom() // открываем все мины
			show() // отображаем поле
			println("\nYou lose... Try again!")
			os.Exit(0) // если попал на мину - выходим
		}
		show() // отображаем поле
	}
}

// getInput вводим координаты выстрела
func getInput(i, j *int) {
	print("\nточка (например 1,2) > ")
	fmt.Scanf("%d,%d", i, j)
}

// start забиваем поле точками
func start() {
	for i := 0; i < fieldsNum; i++ {
		for j := 0; j < fieldsNum; j++ {
			field[i][j] = "."
		}
	}
}

// setMines расставляет мины
func setMines() {
	var n int // счетчик мин
	for n < minesNum {
		x := rand.Intn(fieldsNum) // получаем x координату мины
		y := rand.Intn(fieldsNum) // получаем y координату мины
		if mines[x][y] {          // если мина уже есть в точке - повторяем итерацию цикла
			continue
		}
		mines[x][y] = true // выставляем мину
		n++                // увеличиваем счетчик мин
	}
}

// show просто отображаем текущее состояние поля
func show() {
	print("    ")
	for j := 0; j < fieldsNum; j++ {
		fmt.Printf("%d ", j) // нумерация Y
	}
	print("\n\n")
	for i := 0; i < fieldsNum; i++ {

		fmt.Printf("%d   ", i) // нумерация X

		for j := 0; j < fieldsNum; j++ {
			if mines[i][j] {
				print("o")
			} else {
				print(field[i][j])
			}
			print(" ")
		}
		print("\n")
	}
}

// shut выстрел
func shut(i, j int) bool {
	if i > fieldsNum || j > fieldsNum || i < 0 || j < 0 { // проверка на выход за пределы массива
		return false
	}

	if mines[i][j] {
		field[i][j] = "X" // если попали в мину - помечаем точку как 'X' и завершаем програму
		return true
	}

	status := status(i, j)                  // получаем количество мин в точке
	field[i][j] = fmt.Sprintf("%d", status) // отображаем количество мин
	clean(i, j)                             // открываем пустые клетки

	return false
}

// clean рекурсивно открывает пустые клетки. Выглядет это так:
//
//  точка (например 1,2) > 8,2
//   0 1 2 3 4 5 6 7 8 9
// 0 . . . . . 1 0 0 0 0
// 1 . . . . . 1 0 0 0 0
// 2 . . . . . . 1 1 1 0
// 3 . . . . . . . . 1 0
// 4 . . . . . 2 1 1 1 0
// 5 . . . . . 1 0 0 0 0
// 6 . . . 2 1 1 0 1 1 1
// 7 . 3 1 1 0 0 0 1 . .
// 8 . 1 0 0 0 1 2 3 . .
// 9 . 1 0 0 0 1 . . . .

func clean(i, j int) {
	if i >= fieldsNum || j >= fieldsNum || i < 0 || j < 0 { // проверка на выход за пределы массива
		return
	}
	status := status(i, j)                  // получаем количество мин в точке
	field[i][j] = fmt.Sprintf("%d", status) // отображаем количество мин
	if status > 0 {                         // если количество мин больше нуля - выходим
		return
	}
	if opened[i][j] { // если поле уже открыто - выходим
		return
	}
	opened[i][j] = true // помечаем поле как открытое

	// продолжаем углубляться в рекурсию
	clean(i, j+1)
	clean(i, j-1)
	clean(i+1, j)
	clean(i-1, j)
	clean(i-1, j-1)
	clean(i-1, j+1)
	clean(i+1, j+1)
	clean(i+1, j-1)
}

// status возвращает количество мин вокруг точки
func status(i, j int) int {
	var n int
	if i >= 0 && i < fieldsNum && j-1 >= 0 && j-1 < fieldsNum && mines[i][j-1] {
		n++
	}
	if i >= 0 && i < fieldsNum && j+1 >= 0 && j+1 < fieldsNum && mines[i][j+1] {
		n++
	}
	if i+1 >= 0 && i+1 < fieldsNum && j >= 0 && j < fieldsNum && mines[i+1][j] {
		n++
	}
	if i-1 >= 0 && i-1 < fieldsNum && j >= 0 && j < fieldsNum && mines[i-1][j] {
		n++
	}
	if i+1 >= 0 && i+1 < fieldsNum && j+1 >= 0 && j+1 < fieldsNum && mines[i+1][j+1] {
		n++
	}
	if i-1 >= 0 && i-1 < fieldsNum && j-1 >= 0 && j-1 < fieldsNum && mines[i-1][j-1] {
		n++
	}
	if i-1 >= 0 && i-1 < fieldsNum && j+1 >= 0 && j+1 < fieldsNum && mines[i-1][j+1] {
		n++
	}
	if i+1 >= 0 && i+1 < fieldsNum && j-1 >= 0 && j-1 < fieldsNum && mines[i+1][j-1] {
		n++
	}

	return n
}

// check проверка на открытие всех полей, кроме мин (проверка на победу)
func check() bool {
	var n int
	for x := range field {
		for y := range field[x] {
			if field[x][y] != "." {
				n++
			}
		}
	}

	return (fieldsNum*fieldsNum)-minesNum-n == 0
}

// boom отображаем все мины
func boom() {
	for x := range field {
		for y := range field[x] {
			if mines[x][y] {
				field[x][y] = "X"
			}
		}
	}
}
