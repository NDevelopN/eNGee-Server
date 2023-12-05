package gamedummy

import "fmt"

var instances map[string]GameDummy

func PrepareInstancing() {
	instances = make(map[string]GameDummy)
}

func CreateNewInstance(rid string) (string, error) {
	if rid == "" {
		return "", fmt.Errorf("empty RID provided")
	}

	_, found := instances[rid]
	if found {
		return "", fmt.Errorf("game already exists for room %s", rid)
	}

	instances[rid] = CreateDefaultGame()

	return instances[rid].Address, nil
}

func DeleteInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	err = instance.EndGame()
	if err != nil {
		return err
	}

	delete(instances, rid)

	return nil
}

func StartInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	instance, err = instance.StartGame()
	if err != nil {
		return err
	}

	instances[rid] = instance

	return nil
}

func PauseInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	instance, err = instance.PauseGame()
	if err != nil {
		return err
	}

	instances[rid] = instance

	return nil
}

func ResetInstance(rid string) error {
	instance, err := getInstance(rid)
	if err != nil {
		return err
	}

	instance, err = instance.ResetGame()
	if err != nil {
		return err
	}

	instances[rid] = instance
	return nil
}

func getInstance(rid string) (GameDummy, error) {
	instance, found := instances[rid]
	if !found {
		return instance, fmt.Errorf("game does not exist for room %s", rid)
	}

	return instance, nil
}
