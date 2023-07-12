package consequences

import (
	g "Engee-Server/game"
	"Engee-Server/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func changeState(state string, gid string) {
	cVar := CVars[gid]
	cVar.State = state
	CVars[gid] = cVar
}

func timer(gid string) {

	if CVars[gid].Timer == 0 {
		return
	}

	upd := utils.GameMsg{
		Type: "Timer",
		GID:  gid,
	}

	for CVars[gid].Timer > 0 {
		cVars := CVars[gid]
		t := time.Now()
		if cVars.State == "Lobby" {
			return
		}

		if cVars.State == "Pause" {
			time.Sleep(time.Millisecond * 10)
			continue
		}

		time.Sleep(time.Second * 1)

		elapsed := time.Since(t)
		cVars.Timer -= int(elapsed)

		CVars[gid] = cVars

		upd.Content = fmt.Sprintf("%d", cVars.Timer)

		updatePlayers(upd)
	}

	nextState(gid)
}

func nextState(gid string) {
	//TODO
}

func updatePlayers(utils.GameMsg) {
	//TODO
}

func initialize(msg utils.GameMsg) (utils.GameMsg, error) {
	game, err := g.GetGame(msg.GID)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("could not get game from GID: %v", err))
	}

	var settings ConSettings
	decoder := json.NewDecoder(strings.NewReader(game.AdditionalRules))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&settings)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("could not parse additional rules: %v", err))
	}

	if settings.Rounds < 0 {
		return utils.ReplyError(msg, fmt.Errorf("rounds must not be less than 0: %v", settings.Rounds))
	}

	log.Printf("DEBUG: %v", settings)

	//TODO
	highestShuffle := 3

	if settings.Shuffle < 0 {
		return utils.ReplyError(msg, fmt.Errorf("shuffle must not be less than 0: %v", settings.Shuffle))
	} else if settings.Shuffle > highestShuffle {
		return utils.ReplyError(msg, fmt.Errorf("shuffle option not recognised: %v", settings.Shuffle))
	} else if settings.Shuffle == 0 {
		settings.Shuffle = 1
	}

	if settings.Timer1 < 0 || settings.Timer2 < 0 {
		return utils.ReplyError(msg,
			fmt.Errorf("timers must not be less than 0: (1: %v, 2: %v)", settings.Timer1, settings.Timer2))
	}

	pLen := len(settings.Prompts)
	if pLen == 0 {
		settings.Prompts = defPrompts
	} else if pLen == 1 {
		return utils.ReplyError(msg,
			fmt.Errorf("prompts must be empty or have 2 or more values: %v", pLen))
	}

	cVar := ConVars{
		State:    "Lobby",
		Settings: settings,
		Timer:    settings.Timer1,
		Stories:  map[string][]string{},
	}

	CVars[msg.GID] = cVar

	return utils.ReplyACK(msg), nil
}

func start(msg utils.GameMsg) (utils.GameMsg, error) {
	cVars := CVars[msg.GID]
	cVars.Stories = make(map[string][]string)
	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("error getting game players: %v", err))
	}

	for _, p := range plrs {
		cVars.Stories[p.UID] = []string{}
	}

	cVars.State = "Prompts"
	cVars.Timer = cVars.Settings.Timer1

	CVars[msg.GID] = cVars

	go timer(msg.GID)

	prompts, err := json.Marshal(cVars.Settings.Prompts)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("error marshalling prompts: %v", err))
	}

	upd := utils.GameMsg{
		Type:    "Prompt",
		GID:     msg.GID,
		Content: string(prompts),
	}

	updatePlayers(upd)

	return utils.ReplyACK(msg), nil
}

func reset(msg utils.GameMsg) (utils.GameMsg, error) {
	cVars := CVars[msg.GID]
	cVars.State = "Lobby"
	cVars.Stories = make(map[string][]string)

	CVars[msg.GID] = cVars

	resp, err := initialize(msg)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("error resetting to current settings: %v", err))
	}

	return resp, nil
}

func end(msg utils.GameMsg) (utils.GameMsg, error) {
	_, err := reset(msg)
	if err != nil {
		return utils.ReplyError(msg, fmt.Errorf("error ending settings: %v", err))
	}

	delete(CVars, msg.GID)

	return utils.ReplyACK(msg), nil
}

func pause(msg utils.GameMsg) (utils.GameMsg, error) {
	cVars := CVars[msg.GID]

	if cVars.State == "Pause" {
		cVars.State = cVars.SusState
		cVars.SusState = ""
	} else {
		cVars.SusState = cVars.State
		cVars.State = "Pause"
	}

	CVars[msg.GID] = cVars

	upd := utils.GameMsg{
		Type:    "State",
		GID:     msg.GID,
		Content: cVars.State,
	}

	updatePlayers(upd)

	return utils.ReplyACK(msg), nil
}

func rules(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func status(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func leave(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func reply(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func Handle(msg utils.GameMsg) (utils.GameMsg, error) {
	switch msg.Type {
	case "Init":
		return initialize(msg)
	case "Start":
		return start(msg)
	case "Reset":
		return reset(msg)
	case "End":
		return end(msg)
	case "Pause":
		return pause(msg)
	case "Remove":
		//Can be treated as leave from game mode point of view
		msg.UID = msg.Content
		return leave(msg)
	case "Rules":
		return rules(msg)
	case "Status":
		return status(msg)
	case "Leave":
		return leave(msg)
	case "Reply":
		return reply(msg)
	default:
		return utils.ReplyError(msg, fmt.Errorf("unknown message type: %v", msg.Type))
	}
}

func GetConState(gid string) (ConVars, error) {
	cv, f := CVars[gid]
	if !f {
		return cv, fmt.Errorf("could not find variables for consequences game: %v", gid)
	}

	return cv, nil
}
