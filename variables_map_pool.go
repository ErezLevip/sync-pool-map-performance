package sync_pool_performance

import "sync"

type variablesMapPool struct {
	pool    *sync.Pool
	channel chan map[string]interface{}
}

type variablesMap map[string]interface{}

func newVariablesMapPool(hotSize int, prefill bool, initialMapSize int) *variablesMapPool {
	newInstance := func() variablesMap {
		return make(map[string]interface{}, initialMapSize)
	}

	p := &variablesMapPool{
		pool: &sync.Pool{
			New: func() interface{} { return newInstance() },
		},
	}
	//init a channel to reduce sync variablesMapPool extention
	if hotSize > 0 {
		p.channel = make(chan map[string]interface{}, hotSize)

		if prefill {
			// We prefill our hot data at half the size so that we don't waste time allocating in the variablesMapPool initially.
			for i := 0; i < hotSize; i++ {
				p.channel <- newInstance()
			}
		}
	}
	return p
}

//reset using a compiler optimization
//https://go-review.googlesource.com/c/go/+/110055/
func (m variablesMap) reset() {
	for k := range m {
		delete(m, k)
	}
}

func (bp *variablesMapPool) Get() (b variablesMap) {
	select {
	case b = <-bp.channel:
	default:
		b = bp.pool.Get().(variablesMap)
	}

	b.reset()
	return b
}

func (bp *variablesMapPool) Put(m variablesMap) {
	m.reset()
	select {
	case bp.channel <- m:
	default:
		bp.pool.Put(m)
	}
}
