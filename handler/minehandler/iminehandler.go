package minehandler

type MinerWoker interface {
	GenerateOpBlock()

	GenerateNoOpBlock()

	AddOperation()

	// ValidateBlock()

	// ValidateOperation()
}
