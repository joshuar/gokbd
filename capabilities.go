// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"os/user"
	"strconv"

	"github.com/rs/zerolog/log"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

// taken from https://git.kernel.org/pub/scm/libs/libcap/libcap.git/tree/goapps/setid/setid.go#n32
func setIDsWithCaps(setgid, setuid int, gids []int) {

	if err := cap.SetGroups(setgid, gids...); err != nil {
		log.Fatal().Err(err).Msgf("Unable to set gid to %d.", setgid)
	}
	if err := cap.SetUID(setuid); err != nil {
		log.Fatal().Err(err).Msgf("Unable to set uid to %d.", setuid)
	}
}

func getInputGroupGid() int {
	inputGroup, err := user.LookupGroup("input")
	if err != nil {
		log.Fatal().Err(err).Msg("No input group defined.")
	}
	inputGroupGid, err := strconv.Atoi(inputGroup.Gid)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not convert gid string to int.")
	}
	return inputGroupGid
}

func getUserIds() (int, int) {
	userDetails, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not retrieve user details.")
	}
	uid, err := strconv.Atoi(userDetails.Uid)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not convert uid string to int.")
	}
	gid, err := strconv.Atoi(userDetails.Gid)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not convert gid string to int.")
	}
	return uid, gid
}
