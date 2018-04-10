package cpu

type action struct {
	action     int8
	location   int8
	isRegister bool
	params     []value
}

type value struct {
	value      int8
	isRegister bool
}

func decode(payload []int8) action {
	action := action{
		params: make([]value, 0),
	}

	for i, v := range payload {
		if i == 0 {
			action.action = v
		} else if i == 1 {
			action.location = -(v + 1)
			action.isRegister = true
		} else {
			vl := value{}
			vl.isRegister = v < 0

			if vl.isRegister {
				vl.value = -(v + 1)
			} else {
				vl.value = v
			}

			action.params = append(action.params, vl)
		}
	}

	return action
}
