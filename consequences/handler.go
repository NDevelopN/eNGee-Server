package consequences

import (
	g "Engee-Server/game"
	u "Engee-Server/user"
	"Engee-Server/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var ready = map[string]bool{}

func timer(gid string, cVars ConVars) {
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

	startState := cVars.State

TickLoop:
	for {
		select {
		case <-ticker.C:
			cVars, k := CVars[gid]

			if !k || cVars.State == LOBBY || cVars.State != startState {
				break TickLoop
			}

			if cVars.Paused {
				continue TickLoop
			}

			t -= 1
			cVars.Timer = t

			CVars[gid] = cVars

			if t <= 0 {
				nextState(gid, CVars[gid])
				break TickLoop
			}
		}
	}
}

func sendPrompts(gid string) error {
	cVars := CVars[gid]
	prompts, err := json.Marshal(cVars.Settings.Prompts)
	if err != nil {
		return fmt.Errorf("error marshalling prompts: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Prompts",
		GID:     gid,
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

			reset(upd, true)
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

func nextState(gid string, cVars ConVars) {
	cVars.State++

	cVars.Ready = 0

	if cVars.State > POSTSTORIES {
		cVars.State = LOBBY
	}

	switch cVars.State {
	case LOBBY:
		msg := utils.GameMsg{
			Type: "Reset",
			GID:  gid,
		}
		reset(msg, true)
	case PROMPTS:
		err := sendPrompts(gid)
		if err != nil {
			log.Printf("[Error] Failed to send prompts after state transition: %v", err)
			cVars.State = ERROR
		}
		cVars.Timer = cVars.Settings.Timer1
		go timer(gid, cVars)
	case POSTPROMPTS:
		cVars.Timer = cVars.Settings.Timer1
		go timer(gid, cVars)
	case STORIES:
		err := sendStories(gid, cVars)
		if err != nil {
			log.Printf("[Error] Failed to send stories after state transition: %v", err)
			cVars.State = ERROR
		}
		cVars.Timer = cVars.Settings.Timer2

		go timer(gid, cVars)
	case POSTSTORIES:
		cVars.Timer = cVars.Settings.Timer1
		go timer(gid, cVars)
	}

	uMsg := utils.GameMsg{
		GID:     gid,
		Type:    "ConState",
		Content: fmt.Sprintf("%d", cVars.State),
	}

	err := utils.Broadcast(uMsg)
	if err != nil {
		log.Printf("[Error] Failed to send state update: %v", err)
		cVars.State = ERROR
	}

	CVars[gid] = cVars
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

	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		fmt.Printf("[Error] Could not get players for init: %v", err)
		return "Error", "Could not get players for init."
	}

	act := 0
	for _, p := range plrs {
		if p.Status != "Observing" && p.Status != "Leaving" {
			act++
			p.Status = "Not Ready"
			err := u.UpdateUser(p)
			if err != nil {
				log.Printf("[Error] Could not update user in init: %v", err)
				return "Error", "Could not update user for init."
			}
		}
	}

	cVar := ConVars{
		State:    0,
		Settings: settings,
		Timer:    settings.Timer1,
		Active:   act,
		Ready:    0,
		Stories:  map[string][]string{},
	}

	CVars[msg.GID] = cVar
	ready[msg.GID] = true

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

	if cVars.State != LOBBY {
		log.Printf("[Error] Can only accept Start request when in Lobby State")
		return "Error", "Cannot accept Start when not in Lobby State"
	}

	nextState(msg.GID, cVars)

	if CVars[msg.GID].State == ERROR {
		return "Error", "Game state not valid."
	}

	go checkPhaseChange(msg.GID)

	return "", ""
}

func reset(msg utils.GameMsg, init bool) (string, string) {

	cVars := CVars[msg.GID]
	cVars.State = LOBBY
	cVars.Stories = make(map[string][]string)

	CVars[msg.GID] = cVars

	ready[msg.GID] = false

	cause, resp := "", ""

	if init {
		cause, resp = initialize(msg)
		if cause != "" {
			log.Printf("[Error] Could not reset game state to current settings: %v", resp)
		}
	}

	return cause, resp
}

func end(msg utils.GameMsg) (string, string) {
	cause, resp := reset(msg, false)
	if cause != "" {
		log.Printf("[Error] Could not reset game state before ending: %v", resp)
	}

	delete(CVars, msg.GID)

	return "", ""
}

func pause(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]

	cVars.Paused = !cVars.Paused

	CVars[msg.GID] = cVars

	return "", ""
}

func setActivePlayers(gid string, status string, cVars ConVars) error {
	plrs, err := g.GetGamePlayers(gid)
	if err != nil {
		return fmt.Errorf("couold not get game players: %v", err)
	}

	for _, p := range plrs {
		if p.Status == "Ready" || p.Status == "Not Ready" {
			p.Status = status
			err = u.UpdateUser(p)
			if err != nil {
				return fmt.Errorf("could not update user status: %v", err)
			}
			cVars.Ready--
		}
	}

	CVars[gid] = cVars

	return nil
}

func checkPhaseChange(gid string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cVars := CVars[gid]
			if cVars.State == ERROR || cVars.State == LOBBY {
				return
			}

			plrs, err := g.GetGamePlayers(gid)
			if err != nil {
				log.Printf("[Error] Could not get game players in phase loop: %v", err)
				return
			}
			ready := 0
			for _, p := range plrs {
				if p.Status == "Ready" {
					ready++
				}
			}

			cVars.Ready = ready

			if ready > cVars.Active/2 {
				err := setActivePlayers(gid, "Not Ready", cVars)
				if err != nil {
					log.Printf("[Error] Could not reset active players to 'Not Ready': %v", err)
					return
				}

				nextState(gid, cVars)
			} else {
				CVars[gid] = cVars
			}
		}
	}
}

func status(msg utils.GameMsg) (string, string) {
	return "", ""
}

func leave(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	cVars.Active--

	delete(cVars.Stories, msg.UID)
	CVars[msg.GID] = cVars

	return "", ""
}

func reply(msg utils.GameMsg) (string, string) {
	cVars := CVars[msg.GID]
	if cVars.State != PROMPTS {
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

	user, err := u.GetUser(msg.UID)
	if err != nil {
		log.Printf("[Error] Could not get user to update status: %v", err)
		return "Error", "Could not get user to update status after reply."
	}

	user.Status = "Ready"
	err = u.UpdateUser(user)
	if err != nil {
		log.Printf("[Error] COuld not update user status: %v", err)
		return "Error", "Could not update user status after reply."
	}

	cVars.Stories[msg.UID] = replies
	CVars[msg.GID] = cVars

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

	if msg.Type == "Init" {
		return initialize(msg)
	} else if !ready[msg.GID] {
		time.Sleep(time.Second * 2)
		if !ready[msg.GID] {
			return "Error", "Game has not been initialized."
		}
	}

	switch msg.Type {
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
			return reset(msg, true)
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
