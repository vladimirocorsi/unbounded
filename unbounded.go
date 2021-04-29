package unbounded

import "log"

const initialLenght = 1

//Same as https://github.com/jonbodner/unbounded but using a ring buffer.
//Has an issue with backing array never being shrinked when space is not needed anymore, e.g. queue gets emptied.
func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{} = make([]interface{}, initialLenght)
		readIdx := 0
		writeIdx := 0
		curVal := func() interface{} {
			if readIdx == writeIdx {
				return nil
			}
			val := inQueue[readIdx]
			return val
		}
		outCh := func() chan interface{} {
			if readIdx == writeIdx {
				return nil
			}
			return out
		}
		for readIdx != writeIdx || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue[writeIdx] = v
					if (writeIdx+1)%len(inQueue) == readIdx {
						newLength := 1 + 3*len(inQueue)/2
						log.Printf("reallocating with size %v", newLength)
						newInQueue := make([]interface{}, 0, newLength)
						if readIdx < writeIdx {
							newInQueue = append(newInQueue, inQueue...)
						} else {
							newInQueue = append(newInQueue, inQueue[readIdx:]...)
							newInQueue = append(newInQueue, inQueue[:readIdx]...)
						}
						newInQueue = newInQueue[:newLength]
						readIdx = 0
						writeIdx = len(inQueue)
						inQueue = newInQueue
					} else {
						writeIdx = (writeIdx + 1) % len(inQueue)
					}
				}
			case outCh() <- curVal():
				inQueue[readIdx] = nil
				readIdx = (readIdx + 1) % len(inQueue)
			}
		}
		close(out)
	}()
	return in, out
}
