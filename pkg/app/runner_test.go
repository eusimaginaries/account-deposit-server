package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp_shouldReturnAppWithoutAnyCustomers(t *testing.T) {
	testNewApp := func() {
		res := NewApp()
		assert.Equal(t, App{customers: []*Customer{}}, res)
	}

	testNewApp()
}

func TestProcessInput_shouldProcessInput(t *testing.T) {
	app := NewApp()
	testProcessInput := func(command string, expectedErr error) {
		_, err := app.processInput(command)
		if expectedErr != nil {
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		} else {
			assert.NoError(t, err)
		}
	}

	testProcessInput("newcustomer", app.createNewCustomer([]string{}))
	testProcessInput("addportfolio", app.addPortfolio([]string{}))
	testProcessInput("startDeposit", app.startDeposit())
	testProcessInput("addOneTimePlan", app.addPlan("one-time", []string{}))
	testProcessInput("addMonthlyPlan", app.addPlan("monthly", []string{}))
	testProcessInput("deposit", app.deposit([]string{}))
	testProcessInput("endDeposit", app.endDeposit())
}

func TestProcessInput_shouldReturnError_givenInvalidCommand(t *testing.T) {
	testProcessInput := func(command string) {
		app := NewApp()
		_, err := app.processInput(command)
		assert.Error(t, err)
		assert.Equal(t, "invalid command", err.Error())
	}

	testProcessInput("error-command")
	testProcessInput("\"error command\"")
}

func TestProcessInput_shouldReturnNil_givenNoCommand(t *testing.T) {
	testProcessInput := func(command string) {
		app := NewApp()
		_, err := app.processInput(command)
		assert.NoError(t, err)
	}

	testProcessInput("")
}

func TestCliCreateNewCustomer_shouldCreateNewCustomer(t *testing.T) {
	testCreateNewCustomer := func(args []string) {
		app := NewApp()
		err := app.createNewCustomer(args)
		customer, err := NewCustomer(args[0])
		assert.NoError(t, err)
		assert.Equal(t, []*Customer{&customer}, app.customers)
		assert.Equal(t, &customer, app.currentCustomer)
	}

	testCreateNewCustomer([]string{"test"})
}

func TestCliCreateNewCustomer_shouldReturnError_givenNoArgs(t *testing.T) {
	testCreateNewCustomer := func(args []string) {
		app := NewApp()
		err := app.createNewCustomer(args)
		assert.Error(t, err)
		assert.Equal(t, "invalid number of args", err.Error())
		assert.Empty(t, app.customers)
		assert.Nil(t, app.currentCustomer)
	}

	testCreateNewCustomer([]string{})
}

func TestCliAddPortfolio_shouldAddPortfolioToCustomer(t *testing.T) {
	testAddPortfolio := func(args []string) {
		app := NewApp()
		app.createNewCustomer([]string{"test"})
		err := app.addPortfolio(args)
		assert.NoError(t, err)
		assert.Equal(t, []*Portfolio{{Name: args[0]}}, app.currentCustomer.portfolios)
	}

	testAddPortfolio([]string{"Retirement"})
}

func TestCliAddPortfolio_shouldThrowError_givenNoArgs(t *testing.T) {
	testAddPortfolio := func(args []string) {
		app := NewApp()
		app.createNewCustomer([]string{"test"})
		err := app.addPortfolio(args)
		assert.Error(t, err)
		assert.Equal(t, "invalid number of args", err.Error())
		assert.Empty(t, app.currentCustomer.portfolios)
	}

	testAddPortfolio([]string{})
}

func TestCliAddPortfolio_shouldThrowError_givenNoActiveCustomer(t *testing.T) {
	testAddPortfolio := func(args []string) {
		app := NewApp()
		err := app.addPortfolio(args)
		assert.Error(t, err)
		assert.Equal(t, "no active customer", err.Error())
	}

	testAddPortfolio([]string{"Retirement"})
}

func TestStartDeposit_shouldStartASession(t *testing.T) {
	testStartDeposit := func() {
		app := NewApp()
		app.createNewCustomer([]string{"test"})
		err := app.startDeposit()
		assert.NoError(t, err)
		assert.NotNil(t, app.currentCustomer.DepositSession)
		assert.Equal(t, &DepositSession{[]DepositPlan{}, []float32{}}, app.currentCustomer.DepositSession)
	}

	testStartDeposit()
}

func TestStartDeposit_shouldThrowError_givenNoActiveCustomer(t *testing.T) {
	testStartDeposit := func() {
		app := NewApp()
		err := app.startDeposit()
		assert.Error(t, err)
		assert.Equal(t, "no active customer", err.Error())
	}

	testStartDeposit()
}
