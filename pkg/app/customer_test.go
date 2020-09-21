package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomer_shouldReturnCustomer(t *testing.T) {
	c := NewCustomer("test")
	assert.Equal(t, "test", c.ID)
	assert.Empty(t, c.portfolios)
}

func TestAddPortfolio_shouldAddPortfolioToCustomer(t *testing.T) {
	testAddPortfolio := func(portfolios []*Portfolio, portfolioName string, expected []*Portfolio) {
		c := Customer{ID: "name", portfolios: portfolios}
		err := c.AddPortfolio(portfolioName)
		assert.Nil(t, err)
		assert.Equal(t, expected, c.portfolios)
	}

	testAddPortfolio([]*Portfolio{}, "p1", []*Portfolio{{Name: "p1"}})
	testAddPortfolio([]*Portfolio{{Name: "p1"}, {Name: "p2"}}, "p3", []*Portfolio{{Name: "p1"}, {Name: "p2"}, {Name: "p3"}})
}

func TestAddPortfolio_shouldReturnError_givenDuplicatePortfolioName(t *testing.T) {
	testAddPortfolio := func(portfolios []*Portfolio, portfolioName string, expected []*Portfolio) {
		c := Customer{ID: "name", portfolios: portfolios}
		err := c.AddPortfolio(portfolioName)
		assert.NotNil(t, err)
		assert.Equal(t, "portfolio with specfied name already added", err.Error())
		assert.Equal(t, expected, c.portfolios)
	}

	testAddPortfolio([]*Portfolio{{Name: "p1"}, {Name: "p2"}}, "p1", []*Portfolio{{Name: "p1"}, {Name: "p2"}})
}

func TestAddPortfolio_shouldReturnError_givenPortfolioNameMissing(t *testing.T) {
	testAddPortfolio := func(portfolios []*Portfolio, portfolioName string, expected []*Portfolio) {
		c := Customer{ID: "name", portfolios: portfolios}
		err := c.AddPortfolio(portfolioName)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid name", err.Error())
		assert.Equal(t, expected, c.portfolios)
	}

	testAddPortfolio([]*Portfolio{}, "", []*Portfolio{})
}

func TestStartSession_shouldStartANewSession(t *testing.T) {
	testStartSession := func() {
		c := Customer{}
		err := c.StartSession()
		assert.Nil(t, err)
		assert.NotNil(t, c.DepositSession)
		assert.Equal(t, &DepositSession{depositPlans: []DepositPlan{}, deposits: []float32{}}, c.DepositSession)
	}

	testStartSession()
}

func TestStartSession_shouldReturnError_givenAnotherSessionIsActive(t *testing.T) {
	testStartSession := func() {
		c := Customer{DepositSession: &DepositSession{}}
		err := c.StartSession()
		assert.NotNil(t, err)
		assert.Equal(t, "another transaction is still active", err.Error())
	}

	testStartSession()
}

func TestPayDepositPlan_shouldUpdatePlansInSession(t *testing.T) {
	testPayDepositPlan := func(initialPlans []DepositPlan, plan DepositPlan, expected []DepositPlan) {
		c := Customer{DepositSession: &DepositSession{depositPlans: initialPlans, deposits: []float32{}}}
		err := c.PayDepositPlan(plan)
		assert.Nil(t, err)
		assert.Equal(t, expected, c.DepositSession.depositPlans)
	}

	testPayDepositPlan(
		[]DepositPlan{},
		&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}},
		[]DepositPlan{&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}},
	)
	testPayDepositPlan(
		[]DepositPlan{&baseDepositPlan{name: "test1", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}},
		&baseDepositPlan{name: "test2", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}},
		[]DepositPlan{&baseDepositPlan{name: "test1", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}, &baseDepositPlan{name: "test2", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
	)
}

func TestPayDepositPlan_shouldReturnError_givenNoActiveSession(t *testing.T) {
	testPayDepositPlan := func(plan DepositPlan) {
		c := Customer{DepositSession: nil}
		err := c.PayDepositPlan(plan)
		assert.Error(t, err)
		assert.Equal(t, "no active session", err.Error())
	}

	testPayDepositPlan(&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}})
}

func TestPayDepositPlan_shouldReturnError_givenDuplicateDepositPlanName(t *testing.T) {
	testPayDepositPlan := func(initialPlans []DepositPlan, plan DepositPlan, expected []DepositPlan) {
		c := Customer{DepositSession: &DepositSession{depositPlans: initialPlans, deposits: []float32{}}}
		err := c.PayDepositPlan(plan)
		assert.Error(t, err)
		assert.Equal(t, "duplicate plan name in session", err.Error())
		assert.Equal(t, expected, c.DepositSession.depositPlans)
	}

	testPayDepositPlan(
		[]DepositPlan{&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}},
		&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}},
		[]DepositPlan{&baseDepositPlan{name: "test", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}},
	)
	testPayDepositPlan(
		[]DepositPlan{&baseDepositPlan{name: "test1", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}, &baseDepositPlan{name: "test2", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
		&baseDepositPlan{name: "test2", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}},
		[]DepositPlan{&baseDepositPlan{name: "test1", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 100}}, &baseDepositPlan{name: "test2", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
	)
}

func TestDeposit_shouldDepositAmount(t *testing.T) {
	testDeposit := func(initialDeposits []float32, amount float32, expected []float32) {
		c := Customer{DepositSession: &DepositSession{depositPlans: []DepositPlan{}, deposits: initialDeposits}}
		err := c.Deposit(amount)
		assert.Nil(t, err)
		assert.Equal(t, expected, c.DepositSession.deposits)
	}
	testDeposit([]float32{}, 10.10, []float32{10.1})
	testDeposit([]float32{20}, 0, []float32{20, 0})
}

func TestDeposit_shouldReturnError_givenNoActiveSession(t *testing.T) {
	testDeposit := func(amount float32) {
		c := Customer{DepositSession: nil}
		err := c.Deposit(amount)
		assert.Error(t, err)
		assert.Equal(t, "no active session", err.Error())
	}
	testDeposit(10.10)
}

func TestDeposit_shouldReturnError_givenAmountIsNegative(t *testing.T) {
	testDeposit := func(initialDeposits []float32, amount float32, expected []float32) {
		c := Customer{DepositSession: &DepositSession{depositPlans: []DepositPlan{}, deposits: initialDeposits}}
		err := c.Deposit(amount)
		assert.Error(t, err)
		assert.Equal(t, "amount is negative", err.Error())
	}
	testDeposit([]float32{20}, -1, []float32{20})
}

func TestPerformDeposit_shouldUpdatePortfolioBalance(t *testing.T) {
	testPerformDeposit := func(portfolios []*Portfolio, depositPlans []DepositPlan, deposits []float32, updatedPortfolios []*Portfolio) {
		c := Customer{portfolios: portfolios}
		err := c.PerformDeposit(depositPlans, deposits)
		assert.Nil(t, err)
		assert.Equal(t, updatedPortfolios, c.portfolios)
	}

	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
		[]float32{100},
		[]*Portfolio{{"Retirement", 100}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}, {"High Risk", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
		[]float32{100},
		[]*Portfolio{{"Retirement", 100}, {"High Risk", 0}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}, {"High Risk", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100, "High Risk": 100}}},
		[]float32{200},
		[]*Portfolio{{"Retirement", 100}, {"High Risk", 100}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}, {"High Risk", 0}},
		[]DepositPlan{
			&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100, "High Risk": 100}},
			&baseDepositPlan{name: "Plan B", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 50, "High Risk": 100}},
		},
		[]float32{350},
		[]*Portfolio{{"Retirement", 150}, {"High Risk", 200}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 100}, {"High Risk", 100}},
		[]DepositPlan{
			&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100, "High Risk": 100}},
			&baseDepositPlan{name: "Plan B", planType: "monthly", portfolioRatio: map[string]float32{"Retirement": 50, "High Risk": 100}},
		},
		[]float32{50, 50, 100, 100, 25, 25},
		[]*Portfolio{{"Retirement", 250}, {"High Risk", 300}},
	)
}

func TestPerformDeposit_shouldReturnError_givenCustomerDoesNotHaveTheSpecifiedPortfolio(t *testing.T) {
	testPerformDeposit := func(portfolios []*Portfolio, depositPlans []DepositPlan, deposits []float32, updatedPortfolios []*Portfolio) {
		c := Customer{portfolios: portfolios}
		err := c.PerformDeposit(depositPlans, deposits)
		assert.NotNil(t, err)
		assert.Equal(t, "deposit plan does not match customer portfolio", err.Error())
		assert.Equal(t, updatedPortfolios, c.portfolios)
	}

	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement1": 100}}},
		[]float32{100},
		[]*Portfolio{{"Retirement", 0}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 100}, {"High Risk", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement1": 100, "High Risk": 100}}},
		[]float32{200},
		[]*Portfolio{{"Retirement", 100}, {"High Risk", 0}},
	)
}

func TestPerformDeposit_shouldReturnError_givenDepositAmountDoesNotMeetDepositPlansAmount(t *testing.T) {
	testPerformDeposit := func(portfolios []*Portfolio, depositPlans []DepositPlan, deposits []float32, updatedPortfolios []*Portfolio) {
		c := Customer{portfolios: portfolios}
		err := c.PerformDeposit(depositPlans, deposits)
		assert.NotNil(t, err)
		assert.Equal(t, "deposits does not match the plan amounts", err.Error())
		assert.Equal(t, updatedPortfolios, c.portfolios)
	}

	testPerformDeposit(
		[]*Portfolio{{"Retirement", 0}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
		[]float32{99},
		[]*Portfolio{{"Retirement", 0}},
	)
	testPerformDeposit(
		[]*Portfolio{{"Retirement", 100}},
		[]DepositPlan{&baseDepositPlan{name: "Plan A", planType: "one-time", portfolioRatio: map[string]float32{"Retirement": 100}}},
		[]float32{101},
		[]*Portfolio{{"Retirement", 100}},
	)
}
