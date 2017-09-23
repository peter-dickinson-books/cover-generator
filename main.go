package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"os"

	"github.com/gosimple/slug"
	"github.com/nfnt/resize"
	yaml "gopkg.in/yaml.v2"

	"image/jpeg"
	_ "image/png"
)

var configPath = flag.String("config", "covers.yaml", "the path of the generator config")

func init() {
	flag.Parse()
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// Format describes an end result for an image
type Format struct {
	Name   string `yaml:"name"`
	MaxX   int    `yaml:"max_x"`
	MinX   int    `yaml:"min_x"`
	Source bool   `yaml:"source"`
}

// Config describes the specifications for the project
type Config struct {
	Formats     []Format          `yaml:"formats"`
	SlugFormat  string            `yaml:"slug_format"`
	Questions   map[string]string `yaml:"questions"`
	Destination string            `yaml:"destination"`
}

func initConfig() Config {
	fmt.Println("loading config")
	reader, err := os.Open(*configPath)
	if err != nil {
		fatal(err)
	}
	defer reader.Close()

	fmt.Println("reading config")
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fatal(err)
	}

	var config Config
	fmt.Println("parsing config")
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		fatal(err)
	}

	return config
}

func makeSlug(cfg Config, format Format, x, y int, answers map[string]string) (string, error) {

	payload := map[string]interface{}{
		"FormatName": slug.Make(format.Name),
		"X":          x,
		"Y":          y,
	}

	for k, v := range answers {
		payload[k] = slug.Make(v)
	}

	tmpl, err := template.New(format.Name).Parse(cfg.SlugFormat)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, payload)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func doFormat(orig image.Image, cfg Config, format Format, answers map[string]string) error {
	bounds := orig.Bounds()

	var x, y int
	x = bounds.Max.X - bounds.Min.X

	if x < format.MinX {
		x = format.MinX
	}

	if x > format.MaxX {
		x = format.MaxX
	}

	var m image.Image
	switch format.Source {
	case false:
		m = resize.Resize(uint(x), 0, orig, resize.Lanczos3)
	default:
		m = orig
	}

	bounds = m.Bounds()
	x = bounds.Max.X - bounds.Min.X
	y = bounds.Max.Y - bounds.Min.Y

	fslug, err := makeSlug(cfg, format, x, y, answers)
	if err != nil {
		return err
	}

	fmt.Println(fslug)

	out, err := os.Create(fmt.Sprintf("%s/%s.jpg", cfg.Destination, fslug))
	if err != nil {
		return err
	}
	defer out.Close()

	err = jpeg.Encode(out, m, nil)

	return err
}

func main() {
	config := initConfig()

	answers := make(map[string]string)

	for k, v := range config.Questions {
		fmt.Printf(v + ": ")
		reader := bufio.NewReader(os.Stdin)
		bytes, _, _ := reader.ReadLine()
		answers[k] = string(bytes)
	}

	paths := os.Args[1:]

	for _, path := range paths {
		fmt.Println(path)
		reader, err := os.Open(path)
		if err != nil {
			fatal(err)
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			fatal(err)
		}

		for _, format := range config.Formats {
			err := doFormat(img, config, format, answers)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
