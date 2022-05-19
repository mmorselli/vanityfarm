package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

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

func parseArgs() (int, int, bool) {

	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		fmt.Println("Usage: vanityfarm [minchar] [maxchar]")
		return 0, 0, false
	}

	if len(os.Args) == 3 {
		minchar, err1 := strconv.Atoi(os.Args[1])
		maxchar, err2 := strconv.Atoi(os.Args[2])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid arguments")
			fmt.Println("Usage: vanityfarm [minchar] [maxchar]")
			return 0, 0, false
		} else {
			return minchar, maxchar, true
		}
	}
	return 0, 0, false
}

func main() {

	minchar, maxchar, passed := parseArgs()
	if !passed {
		return
	}

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
					}
				}
				runtime.Gosched()
			}

		}()
	}

	for {
		time.Sleep(10 * time.Minute)
	}

}
