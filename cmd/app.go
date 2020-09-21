package cmd

import (
	"bufio"
	"fmt"
	"os"

	appMod "bitbucket.org/leeyousheng/account-deposit-server/pkg/app"
)

// Run starts the main loop of the app.
func Run() {
	scanner := bufio.NewScanner(os.Stdin)
	if serr := scanner.Err(); serr != nil {
		fmt.Println("Scanner error: ", serr)
	}

	app := appMod.NewApp()
	app.Run(scanner)
}
