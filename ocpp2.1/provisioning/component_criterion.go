package provisioning

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// ComponentCriterionType represents criteria for filtering components.
type ComponentCriterionType string

const (
	ComponentCriterionActive    ComponentCriterionType = "Active"
	ComponentCriterionAvailable ComponentCriterionType = "Available"
	ComponentCriterionEnabled   ComponentCriterionType = "Enabled"
	ComponentCriterionProblem   ComponentCriterionType = "Problem"
)

func isValidComponentCriterion(fl validator.FieldLevel) bool {
	c := ComponentCriterionType(fl.Field().String())
	switch c {
	case ComponentCriterionActive, ComponentCriterionAvailable, ComponentCriterionEnabled, ComponentCriterionProblem:
		return true
	default:
		return false
	}
}

func init() {
	_ = types.Validate.RegisterValidation("componentCriterion21", isValidComponentCriterion)
}
