package mapper

// effectiveEmployeeIDs prefers the plural field, even when empty, and falls
// back to the legacy singular employee_id of imported records.
func effectiveEmployeeIDs(single *uint64, many []uint64) []uint64 {
	if many != nil {
		return many
	}
	if single != nil {
		return []uint64{*single}
	}
	return []uint64{}
}
