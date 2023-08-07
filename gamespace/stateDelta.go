package gamespace

import (
	g "Engee-Server/game"
	utils "Engee-Server/utils"
	"database/sql"
	"log"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

/**
 *  These functions are intended to buffer state changes and updates.
 *	By checking for cumulative changes the intention is to reduce the number of
 *	messages being sent by the server.
 *	The latency of .5s is unlikely to be significant in the expected game types.
 *  TODO: The cycle delay could be a user-specified variable (or gametype specific)
 */

const cycleDelay = 500 * time.Millisecond

func getGame(gid string) (utils.Game, error) {
	game, err := g.GetGame(gid)
	if err == sql.ErrNoRows {
		CleanUp(gid)
		return game, err
	} else if err != nil {
		return game, err
	}

	return game, nil
}

/**
 * This function concerns changes to the game struct itself.
 * Changes to status are most common, but the rules and leader could also change.
 */
func CheckGame(game utils.Game) {
	gid := game.GID

	ticker := time.NewTicker(cycleDelay)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			nuGame, err := getGame(gid)
			if err == sql.ErrNoRows {
				log.Printf("[Alert] No game found. Ending CheckGame")
				return
			} else if err != nil {
				log.Printf("[Error] CheckGame cycle failed: could not get game: %v", err)
				continue Loop
			}

			sChange := false

			if nuGame.Status != game.Status {
				sChange = true
			}

			//Check for non-status changes
			game.Status = nuGame.Status
			game.OldStatus = nuGame.OldStatus

			if nuGame != game {
				game = nuGame
				err := gameUpdateBC(game)
				if err != nil {
					log.Printf("[Error] Broadcast Game failed: %v", err)
				}
				continue Loop
			}

			if sChange {
				err := gameStatusBC(gid, game.Status)
				if err != nil {
					log.Printf("[Error] Broadcast game status failed: %v", err)
				}
			}

		case <-Shutdown[gid]:
			return
		}
	}
}

/**
 * This function concerns changes to the game players.
 * It is not unlikely that many players join, leave or change their status in a short time.
 * This polling compresses those changes into a single broadcast update.
 */
func CheckPlayers(game utils.Game) {

	gid := game.GID

	userList, err := g.GetGamePlayers(gid)
	if err != nil {
		log.Printf("[Error] CheckPlayers start failed: could not get game players: %v", err)
		return
	}

	sort.Slice(userList, func(i, j int) bool {
		return userList[i].UID < userList[j].UID
	})

	ticker := time.NewTicker(cycleDelay)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			nuList, err := g.GetGamePlayers(gid)
			if err == utils.ErrNoGame {
				log.Printf("[Alert] No game found. Ending CheckPlayers")
				return
			} else if err == sql.ErrNoRows {
				log.Printf("[Alert] No players found in game (%v), ending game", gid)

				msg := utils.GameMsg{
					GID:  gid,
					Type: "End",
				}

				_, _ = end(msg, game)

				return
			} else if err != nil {
				log.Printf("[Error] CheckPlayers cycle failed: could not get game players: %v", err)
				continue Loop
			}

			sort.Slice(nuList, func(i, j int) bool {
				return nuList[i].UID < nuList[j].UID
			})

			dif := false

			if len(nuList) != len(userList) {
				dif = true
			}

			//Check for any missing users from last update
			nid := []string{}
			for _, u := range nuList {
				nid = append(nid, u.UID)
			}

			for _, u := range userList {
				if !slices.Contains(nid, u.UID) {
					dif = true
					utils.RemoveConnection(gid, u.UID)
				}
			}

			//Check if there have been any other changes
			if !dif {
				for i, u := range nuList {
					o := userList[i]
					if o != u {
						dif = true
						break
					}
				}
			}

			if dif {
				userList = nuList
				err := pListUpdateBC(gid, userList)
				if err != nil {
					log.Printf("[Error] Broadcast Player change failed: %v", err)
				}

			}
		case <-Shutdown[gid]:
			return
		}
	}
}
