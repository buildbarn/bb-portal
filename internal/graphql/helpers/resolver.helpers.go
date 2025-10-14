package helpers

import (
	"time"

	"entgo.io/contrib/entgql"
)

// StringSliceArrayToPointerArray takes an array of strings and returns an array of string pointers
func StringSliceArrayToPointerArray(strings []string) []*string {
	result := make([]*string, len(strings))
	for i, str := range strings {
		result[i] = &str
	}
	return result
}

func paginationCursorToUTC(cursor *entgql.Cursor[int64]) {
	if cursor == nil || cursor.Value == nil {
		return
	}
	switch v := cursor.Value.(type) {
	case time.Time:
		cursor.Value = v.UTC()
	case *time.Time:
		if v != nil {
			ut := v.UTC()
			cursor.Value = &ut
		}
	}
}

// PaginationCursorsToUTC converts pagination cursors that consist of
// timestamps to UTC instead of local time. When the backend sends the cursors
// to the frontend, they are in UTC. However, when the frontend sends them
// back, they are interpreted as local time. This causes issues since Sqlite
// cannot handle comparisons between timestamps in different timezones.
func PaginationCursorsToUTC(after, before *entgql.Cursor[int64]) {
	paginationCursorToUTC(after)
	paginationCursorToUTC(before)
}
