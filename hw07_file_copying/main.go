package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

var ErrUserInterruption = errors.New("execution was interrupted by user")

var (
	from, to      string
	limit, offset int64
	force         *bool
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	force = flag.Bool("force", false, "override it if already exist")
}

func main() {
	flag.Parse()

	err := overridingAsk(to, *force)
	if err != nil {
		log.Fatal(err)
	}

	err = Copy(from, to, offset, limit)
	if err != nil {
		log.Fatal(fmt.Errorf("copying: %w", err))
	}
}

func overridingAsk(destinationFile string, rewrite bool) error {
	// Проверяем есть ли уже указываемый нами файл или директория в месте назначения
	_, err := os.Stat(destinationFile)
	if err == nil {
		if !rewrite {
			var in string
			fmt.Printf(
				"File with name \"%s\" already exist in destination. ",
				destinationFile,
			)
			fmt.Print("Do you want to override it? (y/n): ")

		AskLoop:
			for {
				_, err := fmt.Scanln(&in)
				if err != nil {
					return fmt.Errorf("read user input error: %w", err)
				}

				switch in {
				case "y":
					break AskLoop
				case "n":
					return ErrUserInterruption
				default:
					fmt.Printf("Unknown answer - \"%s\", please type y or n: ", in)
				}
			}
		}

		return nil
	}

	return nil
}
