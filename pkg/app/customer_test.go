package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
