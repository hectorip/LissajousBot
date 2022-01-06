package main

import (
	"io/ioutil"
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
	"time"
	// "bufio"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/gomail.v2"
	"github.com/joho/godotenv"
	// "github.com/lucasb-eyer/go-colorful"

)

// Declaramos la paleta de Colores que usaremos
// `color` viene de el import `image/color`, y siempre se usa la última
// parte en el código
var palette []color.Color
var cl int

func main() {
	// fmt.Println(cl)
	godotenv.Load()

	rand.Seed(time.Now().UTC().UnixNano())
	dark := rand.Intn(1) == 1
	palette = selectPalette(dark)
	cl = len(palette) - 1
	storeDir := os.Args[1]
	name, writer := createFile(storeDir)
	cycles, freq, delay, decreasing := lissajous(writer)
	fmt.Println(name)
	body := fmt.Sprintf("Cycle count: %d\nFrequency: %2.f\nDelay: %d\nDecrease cycles: %t", int(cycles), freq, delay, (decreasing == 1))
	// sendMail("trigger@applet.ifttt.com", body, name)
	if len(os.Args) > 2 && os.Args[2] == "tweet" {
		tweetPlease(body, name)
	}
}


func selectPalette(dark bool) []color.Color {
	// load palettes
	// Select based on dark or light
	return []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{ 0, 0, 25, 1},
		color.RGBA{ 0, 0, 50, 1},
		color.RGBA{ 0, 0, 75, 1},
		color.RGBA{ 0, 0, 100, 1},
		color.RGBA{ 0, 0, 125, 1},
		color.RGBA{ 0, 0, 150, 1},
		color.RGBA{ 0, 0, 175, 1},
		color.RGBA{ 0, 0, 200, 1},
		color.RGBA{ 0, 0, 225, 1},
		color.RGBA{ 0, 0, 255, 1},
		color.RGBA{ 0, 0, 225, 1},
		color.RGBA{ 0, 0, 200, 1},
		color.RGBA{ 0, 0, 175, 1},
		color.RGBA{ 0, 0, 150, 1},
		color.RGBA{ 0, 0, 125, 1},
		color.RGBA{ 0, 0, 100, 1},
		color.RGBA{ 0, 0, 75, 1},
		color.RGBA{ 0, 0, 50, 1},
		color.RGBA{ 0, 0, 25, 1},
	}
	return []color.Color{
		selectBGColor(),
		color.RGBA{ 25, 0, 0, 1},
		color.RGBA{50, 0, 0, 1},
		color.RGBA{75, 0, 0, 1},
		color.RGBA{100, 0, 0, 1},
		color.RGBA{125, 0, 0, 1},
		color.RGBA{150, 0, 0, 1},
		color.RGBA{175, 0, 0, 1},
		color.RGBA{200, 0, 0, 1},
		color.RGBA{225, 0, 0, 1},
		color.RGBA{255, 0, 0, 1},
		color.RGBA{225, 0, 0, 1},
		color.RGBA{200, 0, 0, 1},
		color.RGBA{175, 0, 0, 1},
		color.RGBA{150, 0, 0, 1},
		color.RGBA{125, 0, 0, 1},
		color.RGBA{100, 0, 0, 1},
		color.RGBA{75, 0, 0, 1},
		color.RGBA{50, 0, 0, 1},
		color.RGBA{25, 0, 0, 1},
		// color.RGBA{0, 0, 0, 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
	}
	return []color.Color{
		selectBGColor(),
		color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},

		color.RGBA{255, 10, 100, 1},
		color.RGBA{0, 200, 50, 1},
		color.RGBA{255, 0, 0, 1},
		color.RGBA{0, 255, 0, 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
		// color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 1},
	}
}

func selectBGColor() color.Color {
	colors := []color.Color{color.Black, color.White}
	index := int(math.Round(rand.Float64()))
	return colors[index]
}

// createFile names and creates a file in the storeDir using Unixtime as the name
func createFile(dir string) (fileName string, writer io.Writer) {
	fileName = fmt.Sprintf("%s/%v.gif", dir, time.Now().Unix())
	writer,  _ = os.Create(fileName)
	// defer f.Close()
	// writer = bufio.NewWriter(f)
	return
}


// Improve this how?
func tweetPlease(body string, fileName string ) {
	api := anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TW_ACCESS_TOKEN"),
		os.Getenv("TW_ACCESS_TOKEN_SECRET"),
		os.Getenv("TW_CLIENT"),
		os.Getenv("TW_CLIENT_SECRET"))
	// var buff = new(bytes.Buffer)
	// lissajous(buff)
	vals := url.Values{}
	// file.Read(buff.Bytes())
	bytefile, err := ioutil.ReadFile(fileName)

	encodedString := base64.StdEncoding.EncodeToString(bytefile)

	media, err := api.UploadMedia(encodedString)
	if err != nil {
		fmt.Println(err)
		return
	}
	vals.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
	fmt.Println(media.MediaID)
	api.PostTweet(body, vals)
}

func sendMail(mail string, body string, fileName string) {

	from := os.Getenv("FROM_MAIL")
	pass := os.Getenv("MAIL_PASS")
	to := mail

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Lissajous!")
	m.SetBody("text/html", body)
	m.Attach(fileName)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

func lissajous(out io.Writer) (oCycles, freq float64, delay, decreasing int) {
	// args := os.Args[1:] // ignorar el primer argumento que es el nombre del comando
	// cycles, _ := strconv.ParseFloat(args[0], 64)
	cycles := rand.Float64() * 20 // Some number between 0 and 50
	oCycles = cycles
	delay = rand.Intn(10) + 1
	const ( // las constantes están disponibles en tiempo de compilación, ser números, strings o booleanos
		res     = 0.00001 // 'sharpnesss'
		size    = 250     // la imagen medirá lo doble
		nframes = 128
		imgSize = 400
	)
	// m, _ := strconv.ParseFloat(args[1], 64) // Mulitplicador de la Frecuencia
	freq = rand.Float64() * 10
	anim := gif.GIF{LoopCount: nframes} // Creando un GIF
	phase := 0.0
	space := (imgSize - size)
	decreasing = rand.Intn(2)
	r := (cycles / nframes) * float64(decreasing)
	frames_color := int(nframes/cl)

	for i := 0; i < nframes; i++ { // Creando cada cuadro de la animación
		rect := image.Rect(0, 0, 2*imgSize+1, 2*imgSize+1) // Se usará como un plano cartesiano
		img := image.NewPaletted(rect, palette)
		var index = uint8(rand.Intn(cl) + 1) // avoid black
		index = uint8(i / frames_color) + 1
		var t2 float64
		cycles = cycles - r
		for t := 0.0; t < cycles*2*math.Pi; t += res {
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
			img.SetColorIndex(size+space+int(x*size), size+space+int(y*size), index)
		}
		phase += 2 * math.Pi / 64
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
	return
}
