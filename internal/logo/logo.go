package logo

import (
	"image/color"
	"math/cmplx"
	"math/rand"

	"github.com/fogleman/gg"
)

const (
	pngPath       = "./static/images/amazing_logo.png"
	width, height = 300, 300

	// комплексное число которое определяет форму фрактала
	c = complex(0.35, 0.35)

	//граница выхода за пределы множества
	// определяет детализацию фрактала
	maxIter = 600

	// множитель масштаба
	//чем ближе к 0, тем больше масштаб
	scale = 0.6

	// множитель контраста
	// чем ближе к 0, тем меньше контрастность
	contrast = 5

	//Диапозон интересной части фрактала, нужен для отображения и сопоставления с позицией пикселя
	xMin, xMax = -scale, scale
	yMin, yMax = -scale, scale
)

var (
	baseColour = color.RGBA{
		R: uint8(rand.Float32() * 255),
		G: uint8(rand.Float32() * 255),
		B: uint8(rand.Float32() * 255),
		A: 255,
	}
)

func Create() {
	drawCtx := gg.NewContext(width, height)

	for px := 0; px < width; px++ {
		for py := 0; py < height; py++ {
			// 	конвертация координат текущего пикселя в координаты
			//  комплексной координатной плоскости
			x := xMin + float64(px)*(xMax-xMin)/float64(width)
			y := yMin + float64(py)*(yMax-yMin)/float64(height)

			// комплексное число соответствующее текущим координатам
			z := complex(x, y)

			i := 0
			for i = 0; i < maxIter; i++ { // если модуль cmplx Z превысит 2 за n итераций,
				z = z*z + c           // то в дальнейшем Z гарантированно убегает к бесконечности
				if cmplx.Abs(z) > 2 { // количество итераций определяет удаленность от множества
					break // чем меньше итераций для выхода за пределы, тем дальше Z
				}
			}

			drawCtx.SetColor(shadedColour(baseColour, i, maxIter))
			drawCtx.SetPixel(px, py)

		}
	}
	if err := drawCtx.SavePNG(pngPath); err != nil {
		panic(err)
	}
}

func shadedColour(baseColour color.RGBA, i, maxIter int) color.RGBA {

	shader := (float64(i) / float64(maxIter) * contrast)

	return color.RGBA{
		R: uint8(float64(baseColour.R) * shader),
		G: uint8(float64(baseColour.G) * shader),
		B: uint8(float64(baseColour.B) * shader),
		A: baseColour.A,
	}
}
