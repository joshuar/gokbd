// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"github.com/rs/zerolog/log"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

// taken from https://git.kernel.org/pub/scm/libs/libcap/libcap.git/tree/goapps/setid/setid.go#n32
// TODO: re-write to just elevate to root privileges
func setIDsWithCaps(setUID, setGID int, gids []int) {
	if err := cap.SetGroups(setGID, gids...); err != nil {
		log.Fatal().Err(err).Msg("Unable to raise gid to root.")
	}
	if err := cap.SetUID(setUID); err != nil {
		log.Fatal().Err(err).Msg("Unable to raise uid to root.")
	}
}
