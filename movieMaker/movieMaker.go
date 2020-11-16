package movieMaker

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
	config "youtube-slideshow/configuration"

	"github.com/icza/mjpeg"
)

type MovieMaker struct {
	width       int32
	height      int32
	fps         int32
	saveNoMusic string
	savePath    string
	mp3Path     string
	FramesCh    chan image.Image
	Stop        chan bool
}

func (m *MovieMaker) InitConfig() {
	m.width = int32(config.Width)
	m.height = int32(config.Height)
	m.fps = int32(config.Fps)
}

func (m *MovieMaker) MakeMovie(movieName string) (string, error) {
	m.saveNoMusic = config.Root + "movies/" + movieName + "/" + movieName + "_nosound.avi"
	m.savePath = config.Root + "movies/" + movieName + "/" + movieName + ".avi"

	aw, err := mjpeg.New(m.saveNoMusic, int32(config.Width), int32(config.Height), int32(config.Fps))
	if err != nil {
		return "", err
	} else {
		log.Println("Started generating a movie: ", movieName, ".avi")
	}

	for {
		select {
		case fr, _ := <-m.FramesCh:
			buf := &bytes.Buffer{}
			err = jpeg.Encode(buf, fr, nil)
			if err != nil {
				return "", err
			}

			err = aw.AddFrame(buf.Bytes())
			if err != nil {
				return "", err
			}
		case <-m.Stop:
			err = aw.Close()
			if err != nil {
				return "", err
			}

			err = m.applyMusic()
			if err != nil {
				return "", err
			}

			os.Remove(m.saveNoMusic)
			os.Remove(m.saveNoMusic + ".idx_")

			log.Println("Movie generation finished")

			if config.Parsed {
				log.Println("Last movie created end program ")
				os.Exit(0)
			}

		}
	}
}
func (m *MovieMaker) applyMusic() error {
	err := m.chooseMusic()
	if err != nil {
		return err
	}
	cmd := exec.Command("ffmpeg", "-i", m.saveNoMusic, "-i", m.mp3Path, "-codec", "copy", "-shortest", m.savePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return errors.New("Error occured on aplying sound to video: " + err.Error())
	}

	return nil
}
func (m *MovieMaker) chooseMusic() error {
	files, err := ioutil.ReadDir(config.MusicDir)
	if err != nil {
		return err
	}

	//Generate random song number
	rand.Seed(time.Now().Unix())
	m.mp3Path = config.MusicDir + strconv.Itoa(rand.Intn(len(files))+1) + ".mp3"

	return nil
}
