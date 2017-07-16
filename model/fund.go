package model

type Fund map[string]float64

func (f Fund) Invert() Fund {
	fund := Fund{}
	for id, income := range f {
		fund[id] = -income
	}
	return fund
}
