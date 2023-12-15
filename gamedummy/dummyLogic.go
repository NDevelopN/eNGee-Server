package gamedummy

import "fmt"

const dummyRules = "Rules"

const (
	NEW    int = 0
	ACTIVE int = 1
	PAUSED int = 2
	RESET  int = 3
	END    int = 4
)

type GameDummy struct {
	Rules  string
	Status int
}

func CreateDefaultGame() GameDummy {
	newDummy := GameDummy{
		Rules:  dummyRules,
		Status: NEW,
	}

	return newDummy
}

func (dummy *GameDummy) SetRules(rules string) error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Rules = rules

	return nil

}

func (dummy *GameDummy) StartGame() error {
	err := checkValidGame(dummy, []int{NEW, RESET})
	if err != nil {
		return err
	}

	dummy.Status = ACTIVE

	return nil
}

func (dummy *GameDummy) PauseGame() error {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return err
	}

	if dummy.Status == ACTIVE {
		dummy.Status = PAUSED
	} else {
		dummy.Status = ACTIVE
	}

	return nil
}

func (dummy *GameDummy) ResetGame() error {
	err := checkValidGame(dummy, []int{ACTIVE, PAUSED})
	if err != nil {
		return err
	}

	dummy.Status = RESET

	return nil
}

func (dummy *GameDummy) EndGame() error {
	err := checkValidGame(dummy, []int{})
	if err != nil {
		return err
	}

	return nil
}

func checkValidGame(dummy *GameDummy, status []int) error {
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
