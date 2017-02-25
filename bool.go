package sat

// Tribool is a tri-state boolean with undefined as the 3rd state.
type Tribool uint8

const (
	True  Tribool = 0
	False         = 1
	Undef         = 2
)

func BoolToTri(b bool) Tribool {
	if b {
		return True
	}

	return False
}
