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
	instance, found := instances[rid]
	if !found {
		return fmt.Errorf("game does not exist for room %s", rid)
	}

	err := instance.EndGame()
	if err != nil {
		return err
	}

	delete(instances, rid)

	return nil
}
