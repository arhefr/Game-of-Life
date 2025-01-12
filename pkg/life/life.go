package life

import (
	"fmt"
	"math/rand"
)

// Структура для игрового мира
type World struct {
	Height int      // Высота сетки
	Width  int      // Ширина сетки
	Cells  [][]bool // Состояния клеток
}

// Возвращает новый пустой мир
func NewWorld(height, width int) *World {

	// Создаём и заполняем двумерный булевой массив клеток
	cells := make([][]bool, height)
	for row := range cells {
		cells[row] = make([]bool, width)
	}

	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

// Заполняет игровое поле живыми клетками
func (w *World) Seed(chance uint) error {

	if chance < 1 || chance > 100 {
		return fmt.Errorf("error chance range: 1 <= %d <= 100", chance)
	}

	// Перебираем все клетки
	for _, row := range w.Cells {
		for cell := range row {
			if rand.Intn(100) <= int(chance) { // Шанс оживления клетки
				row[cell] = true
			}
		}
	}

	return nil
}

// Возвращает количество живых соседей у клетки
func (w *World) Neighbors(x, y int) int {
	var neighbors int

	// Перебираем соседние клетки
	for i := x - 1; i < x+2; i++ {
		for j := y - 1; j < y+2; j++ {
			if (i >= 0 && j >= 0) && (i < w.Height && j < w.Width) && // Проверяем принадлежность координаты к сетке мира
				!(i == x && j == y) && w.Cells[i][j] { // Если соседняя клетка живая
				neighbors++ // Прибавляем кол-во живых соседей
			}
		}
	}

	return neighbors
}

// Возвращает следующее состояние клетки
func (w *World) Next(x, y int) bool {
	n, alive := w.Neighbors(y, x), w.Cells[y][x]

	if n < 4 && n > 1 && alive { // Если клетка живая и кол-во соседей 1 < n < 4:
		return true // Клетка остаётся живой
	}
	if n == 3 && !alive { // Если клетка мёртвая и кол-во соседей 3:
		return true // Клетка становится живой
	}

	return false // В остальных случаях клетка считается мёртвой
}

// Возвращает обновлённое состояние игры
func NextState(oldWorld *World) *World {
	var newWorld *World = NewWorld(oldWorld.Height, oldWorld.Width)

	// Перебираем все клетки
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}

	return newWorld
}
