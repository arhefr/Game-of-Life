package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"Game-of-Life-main/pkg/life"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	fileSaveState string = "stateGame.txt"                // Сохранение состояния игры в stateGame.txt
	fileIconGame  string = "/internal/gameFiles/Icon.png" // Иконка игры
)

var (
	world *life.World // Игровой мир

	worldWidth           int  // Ширина игрового мира
	worldHeight          int  // Высота игрового мира
	chanceAlive          uint // Шанс появления живой клетки
	backgroundColor      color.RGBA
	backgroundColorGame  color.RGBA = color.RGBA{40, 48, 68, 0xff}
	backgroundColorPause color.RGBA = color.RGBA{32, 15, 26, 0xff}

	scaleWindow int  = 1 // Масштаб игры в соответсвие с экраном пользователя
	pauseGame   bool     // Пауза игры
	tick, timer uint     // Обновление кадра
)

type Game struct{}

func main() {

	// Параметры запуска
	flag.IntVar(&worldWidth, "width", 50, "Ширина мира")
	flag.IntVar(&worldHeight, "height", 50, "Высота мира")
	flag.UintVar(&chanceAlive, "chance", 5, "Процент живых клеток")
	flag.UintVar(&tick, "tick", 5, "Обновление кадра раз в tick/60 секунд")
	flag.Parse()

	// Создание игрового мира
	game := &Game{}
	world = life.NewWorld(worldHeight, worldWidth)
	if err := world.Seed(chanceAlive); err != nil {
		log.Fatal(err)
	}

	// Настройка игры
	scaleWindow, err := resizeWindow(worldWidth, worldHeight)
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(worldWidth*scaleWindow, worldHeight*scaleWindow)
	ebiten.SetWindowTitle("Game of Live")
	ebiten.SetMaxTPS(30)
	if icon, err := openImage(fileIconGame); err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	} else {
		log.Fatal(err)
	}

	// Запуск игры
	log.Printf("running game with parameters:\nsize world:%d, %d;\nchance of alive cell:%d;\nscale window: %d.", worldWidth, worldHeight, chanceAlive, scaleWindow)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("error with starting game")
	}
}

func (g *Game) Update() error {

	// SPACE -> пауза
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		pauseGame = !pauseGame
	}

	// S -> сохранение состояния мира в файл
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if err := world.SaveState(fileSaveState); err != nil {
			log.Fatal(err)
		}

		log.Printf("succecs save worldstate in file, height:%d, width:%d", world.Width, world.Height)
	}

	// L -> загрузка состояния мира из файла
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		log.Printf("load worldstate from file, height:%d, width:%d", world.Width, world.Height)

		if err := world.LoadState(fileSaveState); err != nil {
			log.Fatal(err)
		}
	}

	// ENTER -> создание нового мира
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		log.Println("refreshing the world")

		world = life.NewWorld(worldHeight, worldWidth)
		world.Seed(chanceAlive)
	}

	// N -> очистить текущий мир
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world = life.NewWorld(world.Height, worldWidth)
	}

	// Обновление состояния мира
	if !pauseGame && timer >= tick {
		world = life.NextState(world)

		timer = 0
	}

	// Режим римования
	if pauseGame {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			world.Cells[y][x] = true
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			world.Cells[y][x] = false
		}
	}

	timer++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Настройка отрисовки мира
	options := &ebiten.DrawImageOptions{}

	// Отрисовка мира
	switch pauseGame {
	case false:
		backgroundColor = backgroundColorGame
	default:
		backgroundColor = backgroundColorPause
	}
	screen.Fill(backgroundColor)
	background := ebiten.NewImage(worldWidth*scaleWindow, worldHeight*scaleWindow)
	world.Print(background, backgroundColor)
	screen.DrawImage(background, options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return worldWidth, worldHeight
}

// Возвращает масштаб игрового поля в соответсвии с расширением экрана пользователя
func resizeWindow(widthWindow, heightWindow int) (int, error) {
	widthMonitor, heightMonitor := ebiten.Monitor().Size()

	for widthWindow*scaleWindow <= widthMonitor && heightWindow*scaleWindow <= heightMonitor {
		scaleWindow++
	}
	scaleWindow--

	// Обработка ошибок
	if scaleWindow == 0 {
		return 0, fmt.Errorf("error, size world: %d, %d, but size monitor: %d, %d", widthWindow, heightWindow, widthMonitor, heightMonitor)
	}
	if widthWindow < 50 || heightWindow < 50 || widthWindow > 600 || heightWindow > 600 {
		return 0, fmt.Errorf("error, size world: %d, %d, but expected width and height must be > 50 and < 600", widthWindow, heightWindow)
	}

	return scaleWindow, nil
}

// Возвращает картинку по имени файла
func openImage(filename string) (image.Image, error) {

	// Получение директории проекта
	directory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error directory project: %s", err)
	}

	// Чтение файла
	imageFile, err := os.Open(directory + filename)
	if err != nil {
		return nil, fmt.Errorf("error file: %s not found", filename)
	}
	defer imageFile.Close()

	// Декодирование картинки
	image, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("error image decode: %s", err)
	}

	return image, nil
}
