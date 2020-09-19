package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseDepositPlan_shouldReturnError_givenInvalidType(t *testing.T) {
	dp, err := newBaseDepositPlan("TestName", "some-other-type", map[string]float32{"retirement": 100})
	assert.Error(t, err)
	assert.Equal(t, "invalid plan type", err.Error())
	assert.Nil(t, dp)
}

func TestNewMonthlyDepositPlan_shouldCreateNewDepositPlan(t *testing.T) {
	dp, err := NewMonthlyDepositPlan("TestName", map[string]float32{"retirement": 100})
	assert.NoError(t, err)
	assert.Equal(t, dp.Name(), "TestName")
	assert.Equal(t, dp.PlanType(), "monthly")
	assert.Equal(t, dp.PortfolioRatio(), map[string]float32{"retirement": 100})
}

func TestNewMonthlyDepositPlan_shouldReturnError_givenNoNameProvided(t *testing.T) {
	dp, err := NewMonthlyDepositPlan("", map[string]float32{"retirement": 100})
	assert.Error(t, err)
	assert.Equal(t, "name cannot be empty", err.Error())
	assert.Nil(t, dp)
}

func TestNewMonthlyDepositPlan_shouldReturnError_givenNoPortfolioGiven(t *testing.T) {
	dp, err := NewMonthlyDepositPlan("TestName", map[string]float32{})
	assert.Error(t, err)
	assert.Equal(t, "no portfolio defined", err.Error())
	assert.Nil(t, dp)
}

func TestNewOneTimeDepositPlan_shouldCreateNewDepositPlan(t *testing.T) {
	dp, err := NewOneTimeDepositPlan("TestName", map[string]float32{"retirement": 100})
	assert.NoError(t, err)
	assert.Equal(t, dp.Name(), "TestName")
	assert.Equal(t, dp.PlanType(), "one-time")
	assert.Equal(t, dp.PortfolioRatio(), map[string]float32{"retirement": 100})
}

func TestNewOneTimeDepositPlan_shouldReturnError_givenNoNameProvided(t *testing.T) {
	dp, err := NewOneTimeDepositPlan("", map[string]float32{"retirement": 100})
	assert.Error(t, err)
	assert.Equal(t, "name cannot be empty", err.Error())
	assert.Nil(t, dp)
}

func TestNewOneTimeDepositPlan_shouldReturnError_givenNoPortfolioGiven(t *testing.T) {
	dp, err := NewOneTimeDepositPlan("TestName", map[string]float32{})
	assert.Error(t, err)
	assert.Equal(t, "no portfolio defined", err.Error())
	assert.Nil(t, dp)
}
