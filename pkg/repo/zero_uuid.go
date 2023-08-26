package repo

import "github.com/google/uuid"

func isZeroUUID(id uuid.UUID) bool {
	// Maybe an easy way to do this in uuid but didn't see it off hand...
	str := id.String()
	return str == "00000000-0000-0000-0000-000000000000" || str == ""
}
