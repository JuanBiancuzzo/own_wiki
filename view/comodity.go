package view

func AcumulateEvents(amount int, yield FnYield) ([]Event, bool) {
	events := []Event{}
	for range amount {
		if currentEvents, ok := <-yield([]*SceneOperation{}); !ok {
			return events, false

		} else {
			events = append(events, currentEvents...)
		}
	}

	return events, true
}
