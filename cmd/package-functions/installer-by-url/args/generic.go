package args

import (
	"github.com/spf13/cobra"
)

type GenericArgs struct {
	BaseArgs // Extend
}

func BuildGenericFlags(inputCmd *cobra.Command) {
	buildBaseFlags(inputCmd)
}

func (g *GenericArgs) ValidateGeneric() {
	g.validateBase()
}

func (g *GenericArgs) FillGenericValues(cmd *cobra.Command) {
	g.fillBaseValues(cmd)
}
