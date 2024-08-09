package utils

type LBStrategy int

const (
	RoundRobin LBStrategy = iota
	WeightedRoundRobin
	LeastConnection
)

func GetCurrentLBStrategy(givenStrategy string) LBStrategy {
	switch givenStrategy {
	case "weighted-round-robin":
		return WeightedRoundRobin
	case "least-connection":
		return LeastConnection
	default:
		return RoundRobin
	}
}
