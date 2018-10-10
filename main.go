package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

// Declaramos la paleta de Colores que usaremos
// `color` viene de el import `image/color`, y siempre se usa la última
// parte en el código
var palette = []color.Color{
	color.Black,
	color.RGBA{255, 10, 100, 1},
	color.RGBA{0, 200, 50, 1},
	color.RGBA{255, 0, 0, 1},
	color.RGBA{0, 255, 0, 1},
	color.RGBA{0, 0, 255, 1},
}
var cl = len(palette) - 1

func main() {
	// fmt.Println(cl)
	lissajous(os.Stdout)
	// tweetPlease()
}

// Improve this
func tweetPlease() {
	// "178848278-zkemS201RssNCnxRo6fF0dka4jvwJmdEn1KLvt4F",
	//"n64pxfdGmYZZuyOBiETKdNx0sy4BuoC8bAtyYFukPcWn2", "X7k4YMZVwo1EPxA4apRPoEJPS", "6EgrBsM2XhoIGQloajpKjz0lPlDWBv8xGXKRX4wrzJbqldyzHZ"
	api := anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TW_ACCESS_TOKEN"),
		os.Getenv("TW_ACCESS_TOKEN_SECRET"),
		os.Getenv("TW_CLIENT"),
		os.Getenv("TW_CLIENT_SECRET"))
	var buff = new(bytes.Buffer)
	lissajous(buff)
	vals := url.Values{}
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	media, err := api.UploadMedia(encodedString)
	if err != nil {
		fmt.Println(err)
		return
	}
	vals.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
	fmt.Println(media.MediaID)
	api.PostTweet("Lissajous", vals)
}
func lissajous(out io.Writer) {
	args := os.Args[1:] // ignorar el primer argumento que es el nombre del comando
	cycles, _ := strconv.ParseFloat(args[0], 64)
	const ( // las constantes están disponibles en tiempo de compilación, ser números, strings o booleanos
		res     = 0.00001 // 'sharpnesss'
		size    = 250     // la imagen medirá lo doble
		nframes = 64
		delay   = 4
	)
	m, _ := strconv.ParseFloat(args[1], 64) // Mulitplicador de la Frecuencia
	freq := rand.Float64() * m
	anim := gif.GIF{LoopCount: nframes} // Creando un GIF
	phase := 0.0
	for i := 0; i < nframes; i++ { // Creando cada cuadro de la animación
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) // Se usará como un plano cartesiano
		img := image.NewPaletted(rect, palette)
		var index = uint8(rand.Intn(cl) + 1)
		var t2 float64
		for t := 0.0; t < cycles*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			t2 += res

			// Changing color every cycle
			// if math.Pi*-t2 <= 0.1 {
			// 	index = uint8(rand.Intn(cl) + 1)
			// 	t2 = 0
			// }
			// Creating stripes of specific colors across all
			// the frames

			// if y >= 0.0 && y <= 0.3 {
			// 	index = 5
			// } else if y >= 0.3 && y <= 0.6 {
			// 	index = 1
			// } else if y >= 0.6 && y <= 0.9 {
			// 	index = 2
			// } else {
			// 	index = 0
			// }
			img.SetColorIndex(size+int(x*size), size+int(y*size), index)
		}
		phase += math.Pi / 64
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
