package main

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var (
	String1 string
	HomeDir string
)

// SetupConfig reads from the teller's config file and authenticates with the platform
func SetupConfig() (string, error) {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return "", errors.Wrap(err, "error while reading email values from config file")
	}

	String1 = viper.Get("String1").(string)
	return "", nil
}

func main() {

	_, err := SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	HomeDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	promptColor := color.New(color.FgHiYellow).SprintFunc()
	whiteColor := color.New(color.FgHiWhite).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      promptColor("emulator") + whiteColor("# "),
		HistoryFile: HomeDir + "/history_emulator.txt",
		// AutoComplete: autoComplete(),
	})

	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		// setup reader with max 4K input chars
		msg, err := rl.Readline()
		if err != nil {
			log.Println("could not read user input / SIGINT received")
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		rl.SaveHistory(msg)

		cmdslice := strings.Fields(msg)
		ColorOutput("entered command: "+msg, YellowColor)

		err = ParseInput(cmdslice)
		if err != nil {
			ColorOutput(err.Error(), RedColor)
			log.Println(err)
			continue
		}
	}
}
