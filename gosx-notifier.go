package gosxnotifier

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Sound string

const (
	Default Sound = "'default'"
	Basso   Sound = "Basso"
	Blow    Sound = "Blow"
	Bottle  Sound = "Bottle"
	Frog    Sound = "Frog"
	Funk    Sound = "Funk"
	Glass   Sound = "Glass"
	Hero    Sound = "Hero"
	Morse   Sound = "Morse"
	Ping    Sound = "Ping"
	Pop     Sound = "Pop"
	Purr    Sound = "Purr"
	Sosumi  Sound = "Sosumi"
	Tink    Sound = "Tink"
)

type Notification struct {
	Message      string //required
	Title        string //optional
	Subtitle     string //optional
	Sound        Sound  //optional
	Link         string //optional
	Sender       string //optional
	Group        string //optional
	Remove       string //optional
	AppIcon      string //optional
	ContentImage string //optional
	Timeout      int    //optional
	Debug        bool
}

var (
	rootPath  string
	FinalPath string
)

func init() {
	if supportedOS() {
		FinalPath = "/usr/local/Cellar/terminal-notifier/1.7.1/bin/terminal-notifier"
	}
}

func NewNotification(message string) *Notification {
	n := &Notification{Message: message}
	return n
}

func (n *Notification) Push() error {
	if supportedOS() {
		commandTuples := make([]string, 0)

		//check required commands
		if n.Message == "" {
			return errors.New("Please specifiy a proper message argument.")
		} else {
			commandTuples = append(commandTuples, []string{"-message", n.Message}...)
		}

		//add title if found
		if n.Title != "" {
			commandTuples = append(commandTuples, []string{"-title", n.Title}...)
		}

		//add subtitle if found
		if n.Subtitle != "" {
			commandTuples = append(commandTuples, []string{"-subtitle", n.Subtitle}...)
		}

		//add sound if specified
		if n.Sound != "" {
			commandTuples = append(commandTuples, []string{"-sound", string(n.Sound)}...)
		}

		//add group if specified
		if n.Group != "" {
			commandTuples = append(commandTuples, []string{"-group", n.Group}...)
		}

		//add group if specified
		if n.Remove != "" {
			commandTuples = append(commandTuples, []string{"-remove", n.Remove}...)
		}

		//add appIcon if specified
		if n.AppIcon != "" {
			img, err := normalizeImagePath(n.AppIcon)

			if err != nil {
				return err
			}

			commandTuples = append(commandTuples, []string{"-appIcon", img}...)
		}

		//add contentImage if specified
		if n.ContentImage != "" {
			img, err := normalizeImagePath(n.ContentImage)

			if err != nil {
				return err
			}
			commandTuples = append(commandTuples, []string{"-contentImage", img}...)
		}

		//add url if specified
		url, err := url.Parse(n.Link)
		if err != nil {
			n.Link = ""
		}
		if url != nil && n.Link != "" {
			commandTuples = append(commandTuples, []string{"-open", n.Link}...)
		}

		//add bundle id if specified
		if strings.HasPrefix(strings.ToLower(n.Link), "com.") {
			commandTuples = append(commandTuples, []string{"-activate", n.Link}...)
		}

		//add sender if specified
		if strings.HasPrefix(strings.ToLower(n.Sender), "com.") {
			commandTuples = append(commandTuples, []string{"-sender", n.Sender}...)
		}

		if len(commandTuples) == 0 {
			return errors.New("Please provide a Message and Type at a minimum.")
		}

		//add timeout if specified
		if n.Timeout > 0 {
			commandTuples = append(commandTuples, []string{"-timeout", strconv.Itoa(n.Timeout)}...)
		}

		if n.Debug {
			fmt.Println(FinalPath)
			fmt.Println(commandTuples)
		}

		_, err = exec.Command(FinalPath, commandTuples...).Output()
		if err != nil {
			return err
		}
	}
	return nil
}

func supportedOS() bool {
	if runtime.GOOS == "darwin" {
		return true
	} else {
		log.Print("OS does not support terminal-notifier")
		return false
	}
}

func normalizeImagePath(image string) (string, error) {
	if imagePath, err := filepath.Abs(image); err != nil {
		return "", errors.New("Could not resolve image path of image: " + image)
	} else {
		return imagePath, nil
	}
}
