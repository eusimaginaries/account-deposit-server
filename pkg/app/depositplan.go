package app

import "errors"

// DepositPlan is the plan used to determine the splitting of deposits to the various portfolio
type DepositPlan interface {
	Name() string
	PlanType() string
	PortfolioRatio() map[string]float32
	DepositTotal() float32
}

type baseDepositPlan struct {
	name           string
	planType       string
	portfolioRatio map[string]float32
}

func newBaseDepositPlan(name string, planType string, portfolioRatio map[string]float32) (DepositPlan, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if planType != "monthly" && planType != "one-time" {
		return nil, errors.New("invalid plan type")
	}

	if len(portfolioRatio) == 0 {
		return nil, errors.New("no portfolio defined")
	}

	for k, v := range portfolioRatio {
		if k == "" || v < 0 {
			return nil, errors.New("invalid portfolio ratio")
		}
	}

	return &baseDepositPlan{name: name, planType: planType, portfolioRatio: portfolioRatio}, nil
}

func (dp *baseDepositPlan) Name() string {
	return dp.name
}

func (dp *baseDepositPlan) PlanType() string {
	return dp.planType
}

func (dp *baseDepositPlan) PortfolioRatio() map[string]float32 {
	return dp.portfolioRatio
}

func (dp *baseDepositPlan) DepositTotal() float32 {
	var sum float32
	for _, v := range dp.portfolioRatio {
		sum += v
	}
	return sum
}

type monthlyDepositPlan struct {
	baseDepositPlan
}

// NewMonthlyDepositPlan creates a new monthly deposit plan
func NewMonthlyDepositPlan(name string, portfolioRatio map[string]float32) (DepositPlan, error) {
	return newBaseDepositPlan(name, "monthly", portfolioRatio)
}

type onetimeDepositPlan struct {
	baseDepositPlan
}

// NewOneTimeDepositPlan creates a new one-time deposit plan
func NewOneTimeDepositPlan(name string, portfolioRatio map[string]float32) (DepositPlan, error) {
	return newBaseDepositPlan(name, "one-time", portfolioRatio)
}
