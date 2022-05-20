package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

//*********************************************
// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string, minchar int, maxchar int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) >= minchar && len(scanner.Text()) <= maxchar {
			lines = append(lines, strings.ToUpper(scanner.Text()))
		}
	}
	return lines, scanner.Err()
}

//*********************************************
func parseArgs() (int, int, string, bool) {

	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		fmt.Println("Usage: vanityfarm <minchar> <maxchar> [logfile]")
		return 0, 0, "", false
	}

	if len(os.Args) >= 3 {
		minchar, err1 := strconv.Atoi(os.Args[1])
		maxchar, err2 := strconv.Atoi(os.Args[2])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid arguments")
			fmt.Println("Usage: vanityfarm [minchar] [maxchar]")
			return 0, 0, "", false
		} else {
			if len(os.Args) == 3 {
				return minchar, maxchar, "", true
			} else {
				return minchar, maxchar, os.Args[3], true
			}
		}
	}
	return 0, 0, "", false
}

//*********************************************
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCtrl+C pressed, bye...")
		os.Exit(0)
	}()
}

//*********************************************
func LogToFile(logfile string, word string, address string, mnemonic string) {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Fprintf(f, "%s,%s,%s\n", word, address, mnemonic)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

//*********************************************
func main() {

	minchar, maxchar, logfile, passed := parseArgs()
	if !passed {
		return
	}

	SetupCloseHandler()

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("Matching started using %d CPU\n", runtime.NumCPU())
	fmt.Printf("String length between %d and %d characters\n", minchar, maxchar)

	dictionary, err := readLines("dictionary.txt", minchar, maxchar)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	fmt.Printf("%d words to compare\n", len(dictionary))

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				dictlimit := len(dictionary)
				account := crypto.GenerateAccount()
				for a := 0; a < dictlimit; a++ {
					if strings.HasPrefix(account.Address.String(), dictionary[a]) {
						fmt.Printf("\nFound: %s\n", dictionary[a])
						fmt.Println(account.Address)
						mnemonic, _ := mnemonic.FromPrivateKey(account.PrivateKey)
						fmt.Println(mnemonic)
						if logfile != "" {
							LogToFile(logfile, dictionary[a], account.Address.String(), mnemonic)
						}
					}
				}
				runtime.Gosched()
			}

		}()
	}

	for {

		time.Sleep(10 * time.Second)

	}
}
