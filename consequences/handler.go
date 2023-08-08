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

var tickerStop = map[string]chan byte{}

func timer(msg utils.GameMsg, cVars ConVars) {
	gid := msg.GID

	if cVars.Timer == 0 {
		return
	}

	upd := utils.GameMsg{
		Type:    "ConTimer",
		GID:     gid,
		Content: fmt.Sprintf("%d", cVars.Timer),
	}

	err := utils.Broadcast(upd)
	if err != nil {
		log.Printf("[Error] Could not broadcast new timer: %v", err)
		return
	}

	t := cVars.Timer

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

TickLoop:
	for {
		select {
		case <-ticker.C:
			cVars = CVars[gid]
			if cVars.Paused {
				continue TickLoop
			}

			t -= 1
			cVars.Timer = t

			CVars[gid] = cVars

			if t <= 0 {
				nextState(msg, CVars[gid])
				break TickLoop
			}

		case <-tickerStop[gid]:
			delete(tickerStop, gid)
			tickerStop[gid] = make(chan byte)

			ticker.Stop()
			break TickLoop
		}

	}
}

func sendPrompts(msg utils.GameMsg) error {
	cVars := CVars[msg.GID]
	prompts, err := json.Marshal(cVars.Settings.Prompts)
	if err != nil {
		return fmt.Errorf("error marshalling prompts: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Prompts",
		GID:     msg.GID,
		Content: string(prompts),
	}

	utils.Broadcast(upd)

	return nil
}

func shuffle(stories map[string][]string) map[string][]string {

	//Used to allow int based index into map
	key := map[int]string{}
	shuffled := map[string][]string{}

	i := 0
	for k := range stories {
		key[i] = k
		shuffled[k] = []string{}
		i++
	}

	i = 0
	for i = 0; i < len(stories); i++ {
		list := []string{}
		for j := 0; j < len(stories[key[i]]); j++ {
			list = append(list, stories[key[(i+j)%len(stories)]][j])
		}
		shuffled[key[i]] = list
	}

	return shuffled
}

func sendStories(gid string, cVars ConVars) error {
	shuffled := shuffle(cVars.Stories)

	for k, s := range shuffled {
		story, err := json.Marshal(s)
		if err != nil {
			upd := utils.GameMsg{
				Type:    "ConErr",
				GID:     gid,
				Content: "Could not send shuffled story.",
			}

			reset(upd)
			return fmt.Errorf("could not marshal stories: %v", err)
		}

		sMsg := utils.GameMsg{
			Type:    "Story",
			UID:     k,
			GID:     gid,
			Content: string(story),
		}

		err = utils.SingleMessage(sMsg)
		if err != nil {
			return fmt.Errorf("could not send story: %v", err)
		}
	}

	return nil
}

func nextState(msg utils.GameMsg, cVars ConVars) ConVars {
	cVars.State++
	if cVars.State > POSTSTORIES {
		cVars.State = LOBBY
	}

	switch cVars.State {
	case LOBBY:
		reset(msg)
	case PROMPTS:
		err := sendPrompts(msg)
		if err != nil {
			log.Printf("[Error] Failed to send prompts after state transition: %v", err)
			cVars.State = ERROR
		}
		go timer(msg, cVars)
	case POSTPROMPTS:

	case STORIES:
		err := sendStories(msg.GID, cVars)
		if err != nil {
			log.Printf("[Error] Failed to send stories after state transition: %v", err)
			cVars.State = ERROR
		}

		go timer(msg, cVars)
	case POSTSTORIES:

	}

	uMsg := utils.GameMsg{
		GID:     msg.GID,
		Type:    "ConState",
		Content: fmt.Sprintf("%d", cVars.State),
	}

	err := utils.Broadcast(uMsg)
	if err != nil {
		log.Printf("[Error] Failed to send state update: %v", err)
		cVars.State = ERROR
	}

	return cVars
}

func initialize(msg utils.GameMsg) (string, string) {
	if len(CVars) == 0 {
		CVars = make(map[string]ConVars)
	}

	game, _ := g.GetGame(msg.GID)

	var settings ConSettings

	if game.AdditionalRules != "" {
		decoder := json.NewDecoder(strings.NewReader(game.AdditionalRules))
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&settings)
		if err != nil {
			log.Printf("[Error] Could not read additional rules")
			return "Error", "Could not read additional rules."
		}
	} else {
		settings = DefSettings
	}

	if settings.Rounds < 0 {
		log.Printf("[Error] Rounds %d is less than 0", settings.Rounds)
		return "Error", "Rounds value must not be less than 0."
	}

	//TODO
	highestShuffle := 3

	if settings.Shuffle < 0 {
		log.Printf("[Error] Shuffle %d is less than 0", settings.Shuffle)
		return "Error", "Shuffle value must not be less than 0."
	} else if settings.Shuffle > highestShuffle {
		log.Printf("[Error] Shuffle %d is greater than maxVal (%d)", settings.Shuffle, highestShuffle)
		return "Error", fmt.Sprintf("Shuffle value must not be greater than %d.", highestShuffle)
	} else if settings.Shuffle == 0 {
		settings.Shuffle = 1
	}

	if settings.Timer1 < 0 {
		log.Printf("[Error] Timer1 %d is less than 0", settings.Timer1)
		return "Error", "Timer1 value must not be less than 0."
	}
	if settings.Timer2 < 0 {
		log.Printf("[Error] Timer2 %d is less than 0", settings.Timer2)
		return "Error", "Timer2 value must not be less than 0."
	}

	pLen := len(settings.Prompts)
	if pLen == 0 {
		settings.Prompts = DefPrompts
	} else if pLen == 1 {
		log.Printf("[Error] Only one prompt provided: %v", settings.Prompts)
		return "Error", "2 or more prompts must be provided. Empty prompts will result in the default set being used."
	}

	cVar := ConVars{
		State:    0,
		Settings: settings,
		Timer:    settings.Timer1,
		Stories:  map[string][]string{},
	}

	tickerStop[msg.GID] = make(chan byte)

	CVars[msg.GID] = cVar

	return "", ""
}

func start(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	cVars.Stories = make(map[string][]string)

	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		log.Printf("[Error] Could not get game players: %v", err)
		return "Error", "Could not get the game players."
	}

	for _, p := range plrs {
		cVars.Stories[p.UID] = []string{}
	}

	if cVars.State != 0 {
		log.Printf("[Error] Can only accept Start request when in Lobby State")
		return "Error", "Cannot accept Start when not in Lobby State"
	}

	cVars = nextState(msg, cVars)
	CVars[msg.GID] = cVars

	if cVars.State == ERROR {
		return "Error", "Game state not valid."
	}

	return "", ""
}

func reset(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	cVars.State = 0
	cVars.Stories = make(map[string][]string)

	CVars[msg.GID] = cVars

	tickerStop[msg.GID] <- 1

	cause, resp := initialize(msg)
	if cause != "" {
		log.Printf("[Error] Could not reset game state to current settings: %v", resp)
	}

	return cause, resp
}

func end(msg utils.GameMsg) (string, string) {
	cause, resp := reset(msg)
	if cause != "" {
		log.Printf("[Error] Could not reset game state before ending: %v", resp)
	}

	tickerStop[msg.GID] <- 1

	delete(CVars, msg.GID)

	return "", ""
}

func pause(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]

	cVars.Paused = !cVars.Paused

	CVars[msg.GID] = cVars

	return "", ""
}

func status(msg utils.GameMsg) (string, string) {
	cVar := CVars[msg.GID]

	if cVar.State == LOBBY {
		return "", ""
	}

	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		log.Printf("[Error] Could not get game players when checking for state change: %v", err)
		return "Error", "Could not check for state change."
	}

	count := len(plrs)
	ready := 0

	for _, p := range plrs {
		if p.Status == "Ready" {
			ready++
		}
	}

	if ready >= count/2 {
		tickerStop[msg.GID] <- 1
		nextState(msg, cVar)
	}

	return "", ""
}

func leave(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	delete(cVars.Stories, msg.UID)
	CVars[msg.GID] = cVars

	return "", ""
}

func reply(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	if cVars.State != 1 {
		log.Printf("[Error] Game state does not accept replies")
		return "Error", "Game is not currently accepting replies."
	}

	var replies []string

	err := json.Unmarshal([]byte(msg.Content), &replies)
	if err != nil {
		log.Printf("[Error] Could not parse reply: %v", err)
		return "Error", "Could not parse reply."
	}

	rLen := len(replies)
	pLen := len(cVars.Settings.Prompts)

	if rLen != pLen {
		log.Printf("[Error] Number of replies %d is not equal to number of prompts %d", rLen, pLen)
		return "Error", "Incorrect number of replies."
	}

	for i, r := range replies {
		if r == "" {
			log.Printf("[Error] Empty reply [%d] %s", i, r)
			return "Error", "One or more empty replies received."
		}
	}

	s, f := cVars.Stories[msg.UID]
	if f && len(s) > 0 {
		log.Printf("[Error] There have already been replies received from this user.")
		return "Error", "There have already been replies received from this user."
	}

	cVars.Stories[msg.UID] = replies
	CVars[msg.GID] = cVars

	statMsg := utils.GameMsg{
		Type:    "Status",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: "Replied",
	}

	cause, resp := status(statMsg)
	if cause != "" {
		log.Printf("[Error] Could not handle status change after reply: %v", resp)
		return "Error", "Could not update player status after reply."
	}

	return "", ""
}

func remove(msg utils.GameMsg) (string, string) {
	msg.UID = msg.Content
	cause, resp := leave(msg)
	if cause != "" {
		log.Printf("[Error] Could not remove player: %v", resp)
		return "Error", "Could not remove player: " + resp
	}

	return "", ""
}

func leaderError(mType string) (string, string) {
	cause := "Error"
	msg := "Must be a leader for " + mType + " requests."

	return cause, msg
}

func Handle(msg utils.GameMsg) (string, string) {
	game, _ := g.GetGame(msg.GID)

	leader := msg.UID == game.Leader

	switch msg.Type {
	case "Init":
		if !leader {
			return leaderError("Init")
		} else {
			return initialize(msg)
		}
	case "Start":
		if !leader {
			return leaderError("Start")

		} else {
			return start(msg)
		}
	case "Reset":
		if !leader {
			return leaderError("Reset")

		} else {
			return reset(msg)
		}
	case "End":
		if !leader {
			return leaderError("End")

		} else {
			return end(msg)
		}
	case "Pause":
		if !leader {
			return leaderError("Pause")

		} else {
			return pause(msg)
		}
	case "Remove":
		if !leader {
			return leaderError("Remove")

		} else {
			return remove(msg)
		}
	case "Status":
		return status(msg)
	case "Leave":
		return leave(msg)
	case "Reply":
		return reply(msg)
	default:

		//TODO
	}

	return "", ""
}

func GetConState(gid string) (ConVars, error) {
	cv, f := CVars[gid]
	if !f {
		return cv, fmt.Errorf("could not find variables for consequences game: %v", gid)
	}

	return cv, nil
}
