package life

import (
	"bufio"
	"os"
)

const (
	path string = "/pkg/life/"
)

// Сохранение состояния игры в файл
func (w *World) SaveState(filename string) error {

	// Форматирование состояния игры
	var cellsString string = w.String("1", "0")

	// Директория проекта
	directory, err := os.Getwd()
	if err != nil {
		return err
	}

	// Запись файла
	err = os.WriteFile(directory+path+filename, []byte(cellsString), 0666)
	return err
}

// Загрузка состояния игры из файла
func (w *World) LoadState(filename string) error {
	var read []string

	// Директория проекта
	directory, err := os.Getwd()
	if err != nil {
		return err
	}

	// Читаем файл
	file, err := os.Open(directory + path + filename)
	if err != nil {
		return ErrInvalidFile
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		read = append(read, scanner.Text())
	}

	if len(read) == 0 {
		return nil
	}

	height, width := len(read), len(read[0])
	newWorld := NewWorld(w.Height, w.Width)

	// Загрузка состояия из файла в мир
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			if i >= w.Height || j >= w.Width {
				break
			}

			switch string(read[i][j]) {
			case "1":
				newWorld.Cells[i][j] = true
			case "0":
				newWorld.Cells[i][j] = false
			}
		}
	}

	w.Cells = newWorld.Cells
	return nil
}

// Метод вывода форматированного состояния игры
func (w *World) String(aliveSym, deadSym string) string {
	var rowString, cellsString string

	// Перебираем все клетки
	for i, row := range w.Cells {
		for _, cell := range row {
			if !cell { // Если клетка мёртвая - выводим deadSym
				rowString += deadSym
			} else { // Если клетка живая - выводим aliveSym
				rowString += aliveSym
			}
		}

		cellsString += rowString
		rowString = ""
		if i != len(w.Cells)-1 {
			cellsString += "\n" // Перенос строки
		}
	}

	return cellsString
}