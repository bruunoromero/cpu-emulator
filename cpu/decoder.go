package cpu

type action struct {
	action     int
	location   int
	isRegister bool
	params     []int
}

func decode(payload []int) action {
	action := action{
		params: make([]int, 0),
	}

	for i, v := range payload {
		if i == 0 {
			action.action = v
		} else if i == 1 {
			if v > 0 {
				// If this branch matches, than the value is a memory
				action.location = v
				action.isRegister = false
			} else {
				// If this branch matches, than the value is a regsiter
				action.isRegister = true
				action.location = -(v + 1)
			}
		} else {
			action.params = append(action.params, v)
		}
	}

	return action
}
