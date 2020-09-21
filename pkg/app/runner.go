package app

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"

	"bitbucket.org/leeyousheng/account-deposit-server/pkg/cli"
)

// App contains the variables needed for the running of the application
type App struct {
	customers       []*Customer
	currentCustomer *Customer
}

// NewApp instantiate a new app without any customers
func NewApp() App {
	return App{customers: []*Customer{}, currentCustomer: nil}
}

// Run performs the app loop
func (a *App) Run(scanner *bufio.Scanner) {
	var end bool
	var err error

	fmt.Println("Welcome to the banking session.")
	fmt.Println("Type \"help\" for help")
	fmt.Print("> ")
	for !end && scanner.Scan() {
		input := scanner.Text()
		end, err = a.processInput(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("\n> ")
	}
}

func (a *App) processInput(input string) (bool, error) {
	command, err := cli.ParseCmdInput(input)
	if err != nil {
		if err.Error() == "no command given" {
			return false, nil
		}
		return false, err
	}
	return a.performCommand(command)
}

func (a *App) performCommand(command cli.Command) (bool, error) {
	var err error
	switch command.Command {
	case "newcustomer":
		err = a.createNewCustomer(command.Args)
		if err == nil {
			fmt.Println("New customer created:", command.Args[0])
		}
	case "addportfolio":
		err = a.addPortfolio(command.Args)
		if err == nil {
			fmt.Println("Portfolio added:", command.Args[0])
		}
	case "startDeposit":
		err = a.startDeposit()
		if err == nil {
			fmt.Println("Deposit session started")
		}
	case "addOneTimePlan":
		err = a.addPlan("one-time", command.Args)
		if err == nil {
			fmt.Println("One time deposit plan selected for deposit:", command.Args[0])
		}
	case "addMonthlyPlan":
		err = a.addPlan("monthly", command.Args)
		if err == nil {
			fmt.Println("Monthly deposit plan selected for deposit:", command.Args[0])
		}
	case "deposit":
		err = a.deposit(command.Args)
		if err == nil {
			fmt.Println("Deposit amount:", command.Args[0])
		}
	case "endDeposit":
		err = a.endDeposit()
		if err == nil {
			fmt.Println("Session completed")
			a.printPortfolios()
		}
	case "printPortfolios":
		err = a.printPortfolios()
	case "exit":
		fmt.Println("See ya!!!")
		return true, nil
	case "help":
		printHelp()
	default:
		err = errors.New("invalid command")
	}
	return false, err
}

func (a *App) createNewCustomer(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid number of args")
	}
	c, err := NewCustomer(args[0])
	if err != nil {
		return err
	}

	a.customers = append(a.customers, &c)
	a.currentCustomer = &c
	return nil
}

func (a *App) addPortfolio(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid number of args")
	}
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}
	return a.currentCustomer.AddPortfolio(args[0])
}

func (a *App) startDeposit() error {
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}
	a.currentCustomer.StartSession()
	return nil
}

func (a *App) addPlan(planType string, args []string) error {
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}

	if a.currentCustomer.DepositSession == nil {
		return errors.New("no active session")
	}

	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return errors.New("invalid number of args")
	}

	portfolioRatio := map[string]float32{}
	for i := 1; i < len(args); i += 2 {
		amt, err := strconv.ParseFloat(args[i+1], 32)
		if err != nil {
			return err
		}
		portfolioRatio[args[i]] = float32(amt)
	}

	var dp DepositPlan
	var err error

	switch planType {
	case "one-time":
		dp, err = NewOneTimeDepositPlan(args[0], portfolioRatio)
	case "monthly":
		dp, err = NewMonthlyDepositPlan(args[0], portfolioRatio)
	default:
		return errors.New("invalid plan type")
	}

	if err != nil {
		return err
	}

	return a.currentCustomer.PayDepositPlan(dp)
}

func (a *App) deposit(args []string) error {
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}

	if len(args) < 1 {
		return errors.New("amount not specified")
	}

	amount, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		return err
	}

	return a.currentCustomer.Deposit(float32(amount))
}

func (a *App) endDeposit() error {
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}

	if a.currentCustomer.DepositSession == nil {
		return errors.New("no active session")
	}

	err := a.currentCustomer.PerformDeposit(a.currentCustomer.DepositSession.depositPlans, a.currentCustomer.DepositSession.deposits)
	if err != nil {
		return err
	}

	a.currentCustomer.DepositSession = nil
	return nil
}

func (a *App) printPortfolios() error {
	if a.currentCustomer == nil {
		return errors.New("no active customer")
	}

	a.currentCustomer.PrintPortfolio()
	return nil
}

func printHelp() {
	fmt.Println("Sample flow:")
	fmt.Println("newcustomer test1")
	fmt.Println("addportfolio Retirement")
	fmt.Println("addportfolio \"High Risk\"")
	fmt.Println("startDeposit")
	fmt.Println("addOneTimePlan \"One Time Plan 1\" \"High Risk\" 10000 Retirement 500")
	fmt.Println("addMonthlyPlan \"Monthly Plan 1\" Retirement 100")
	fmt.Println("deposit 10500")
	fmt.Println("deposit 100")
	fmt.Println("endDeposit")
	fmt.Println("printPortfolios")
}
