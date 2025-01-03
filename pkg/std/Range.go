package std

func LazyRange(start, step int, end ...int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		current := start
		// Check if it's an infinite range
		hasEnd := len(end) > 0
		limit := 0
		if hasEnd {
			limit = end[0]
		}

		for {
			// If finite range, stop when the limit is reached
			if hasEnd && ((step > 0 && current >= limit) || (step < 0 && current <= limit)) {
				break
			}
			ch <- current
			current += step
		}
	}()

	return ch
}
