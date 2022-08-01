package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
)

var (
	help = flag.Bool("help", false, "show help")
)

type config struct {
	Token string `envconfig:"SLACK_USER_OAUTH_TOKEN" required:"true"`
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if err := runMain(); err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to process envconfig: %w", err)
	}

	c := slack.New(cfg.Token)
	emojis, err := c.GetEmoji()
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	fmt.Printf("Output directory: %s\n", tmpDir)

	for emoji, url := range emojis {
		if strings.HasPrefix(url, "alias:") {
			continue
		}

		if err := download(url, filepath.Join(tmpDir, fmt.Sprintf("%s.png", emoji))); err != nil {
			return err
		}
	}

	return nil
}

func download(url, filepath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, response.Body)

	return nil
}
