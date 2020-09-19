package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPortfolio_shouldCreateANewPortfolioWithBalance0(t *testing.T) {
	pName := "testname"
	res, err := NewPortfolio(pName)
	assert.NoError(t, err)
	assert.Equal(t, pName, res.Name)
	assert.Equal(t, float32(0), res.Balance)
}

func TestNewPortfolio_shouldReturnError_givenInvalidName(t *testing.T) {
	pName := ""
	res, err := NewPortfolio(pName)
	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Equal(t, Portfolio{}, res)
}

func TestDeposit_shouldUpdateBySetAmount_givenValidAmount(t *testing.T) {
	testDeposit := func(amount float32, initialBalance float32, expectedBalance float32) {
		p := Portfolio{Name: "test", Balance: initialBalance}
		err := p.Deposit(amount)
		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, p.Balance)
	}

	testDeposit(10.10, 0, 10.10)
	testDeposit(0, 0, 0)
	testDeposit(10.10, 10.20, 20.30)
	testDeposit(10, 0, 10)
}

func TestDeposit_shouldReturnError_givenInvalidAmount(t *testing.T) {
	testDeposit := func(amount float32, initialBalance float32) {
		p := Portfolio{Name: "test", Balance: initialBalance}
		err := p.Deposit(amount)
		assert.Error(t, err)
		assert.Equal(t, "amount is negative", err.Error())
		assert.Equal(t, initialBalance, p.Balance)
	}

	testDeposit(-10.10, 0)
	testDeposit(-10.10, 10.20)
}

func TestWithdraw_shouldUpdateBySetAmount_givenValidAmount(t *testing.T) {
	testWithdraw := func(amount float32, initialBalance float32, expectedBalance float32) {
		p := Portfolio{Name: "test", Balance: initialBalance}
		err := p.Withdraw(amount)
		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, p.Balance)
	}

	testWithdraw(10.10, 20, 9.9)
	testWithdraw(0, 20, 20)
	testWithdraw(0, 0, 0)
	testWithdraw(10.10, 10.10, 0)
}

func TestWithdraw_shouldReturnError_givenInvalidAmount(t *testing.T) {
	testWithdraw := func(amount float32, initialBalance float32) {
		p := Portfolio{Name: "test", Balance: initialBalance}
		err := p.Withdraw(amount)
		assert.Error(t, err)
		assert.Equal(t, "amount is negative", err.Error())
		assert.Equal(t, initialBalance, p.Balance)
	}

	testWithdraw(-10.10, 0)
	testWithdraw(-10.10, 10.20)
}

func TestWithdraw_shouldReturnError_givenAmountMoreThanBalance(t *testing.T) {
	testWithdraw := func(amount float32, initialBalance float32) {
		p := Portfolio{Name: "test", Balance: initialBalance}
		err := p.Withdraw(amount)
		assert.Error(t, err)
		assert.Equal(t, "withdrawal amount more than balance", err.Error())
		assert.Equal(t, initialBalance, p.Balance)
	}

	testWithdraw(10.10, 0)
	testWithdraw(10.11, 10.10)
}
