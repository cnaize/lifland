package model

import "fmt"

type Fund map[*Player]float64

func (f Fund) Revert() error {
	fmt.Println("reverting funds")
	for pl, income := range f {
		if err := pl.IncrBalance(-income); err == nil {
			fmt.Printf("funds %f for player %s reverted, balance %f\n",
				-income, pl.Id(), pl.Balance())
			delete(f, pl)
		}
	}
	// if not all points reverted
	if len(f) != 0 {
		return fmt.Errorf("reverting failed")
	}
	return nil
}
