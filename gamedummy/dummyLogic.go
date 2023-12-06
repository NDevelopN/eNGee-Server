package gamedummy

import "fmt"

const dummyAddress = "http://localhost:8099"
const dummyRules = "Rules"

const (
	NEW    int = 0
	ACTIVE int = 1
	PAUSED int = 2
	RESET  int = 3
	END    int = 4
)

type GameDummy struct {
	Address string
	Rules   string
	Status  int
}

func CreateDefaultGame() GameDummy {
	newDummy := GameDummy{
		Address: dummyAddress,
		Rules:   dummyRules,
		Status:  NEW,
	}

	return newDummy
}

func (dummy GameDummy) SetRules(rules string) (GameDummy, error) {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return dummy, err
	}

	dummy.Rules = rules

	return dummy, nil

}

func (dummy GameDummy) StartGame() (GameDummy, error) {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return dummy, err
	}

	dummy.Status = ACTIVE

	return dummy, nil
}

func (dummy GameDummy) PauseGame() (GameDummy, error) {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return dummy, err
	}

	if dummy.Status == ACTIVE {
		dummy.Status = PAUSED
	} else {
		dummy.Status = ACTIVE
	}

	return dummy, nil
}

func (dummy GameDummy) ResetGame() (GameDummy, error) {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return dummy, err
	}

	dummy.Status = RESET

	return dummy, nil
}

func (dummy GameDummy) EndGame() error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	return nil
}

func checkValidGame(dummy GameDummy, status []int) error {
	if dummy.Address == "" {
		return fmt.Errorf("address is not set")
	}

	if dummy.Rules == "" {
		return fmt.Errorf("rules are not set")
	}

	if len(status) > 0 {
		validStatus := false

		for _, s := range status {
			if dummy.Status == s {
				validStatus = true
				continue
			}
		}

		if !validStatus {
			return fmt.Errorf("game is not in a valid state")
		}
	}

	return nil
}
