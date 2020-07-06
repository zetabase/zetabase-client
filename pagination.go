package zetabase

import (
	"sync"
)

type paginationRequester func(int64) (map[string][]byte, bool, error)

// Type PaginationHandler manages pagination of Zetabase responses.
type PaginationHandler struct {
	requester   paginationRequester
	curData     map[string][]byte
	curPage     int64
	curError    error
	hasNextPage bool
	lock        *sync.Mutex
}

// Get data fetched so far (for data retrieval responses)
func (p *PaginationHandler) Data() (map[string][]byte, error) {
	if p.curError != nil {
		return nil, p.curError
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	return p.curData, nil
}

// Get keys fetched so far
func (p *PaginationHandler) Keys() ([]string, error) {
	if p.curError != nil {
		return nil, p.curError
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	var ks []string
	for k, _ := range p.curData {
		ks = append(ks, k)
	}
	return ks, nil
}

// Get all data from all pages (for data retrieval responses)
func (p *PaginationHandler) DataAll() (map[string][]byte, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var i int64
	i = 1
	data := p.curData
	for p.hasNextPage {
		dat, nxt, err := p.requester(i)
		if err != nil {
			return nil, err
		}
		for k, v := range dat {
			data[k] = v
		}
		i++
		p.hasNextPage = nxt
		p.curPage = i
		if len(dat) == 0 || (!nxt) {
			break
		}
	}
	p.hasNextPage = false
	p.curData = data
	return p.curData, nil
}

// Get all keys from all pages (for data retrieval responses)
func (p *PaginationHandler) KeysAll() ([]string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Page 0
	var ks []string
	for k, _ := range p.curData {
		ks = append(ks, k)
	}
	if !p.hasNextPage {
		return ks, nil
	}
	var i int64
	i = 1
	for {
		dat, nxt, err := p.requester(i)
		if err != nil {
			return nil, err
		}
		for k, _ := range dat {
			ks = append(ks, k)
			p.curData[k] = nil
		}
		i++
		p.hasNextPage = nxt
		p.curPage = i
		if len(dat) == 0 || (!nxt) {
			break
		}
	}
	p.hasNextPage = false
	return ks, nil
}

// Fetch next page worth of data
func (p *PaginationHandler) Next() {
	if !p.hasNextPage {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.curPage += 1
	dat, nxt, err := p.requester(p.curPage)
	if err != nil {
		p.curError = err
		p.hasNextPage = false
	} else {
		p.curData = dat
		p.hasNextPage = nxt
		//log.Printf("pagination: hasNextPage = %v", p.hasNextPage)
		//log.Printf("pagination: dat = %v", dat)
	}
}

// Build a PaginationHandler given a function that takes as input a page index and
// returns the page data
func StandardPaginationHandlerFor(f func(int64) (map[string][]byte, bool, error)) *PaginationHandler {
	ph := &PaginationHandler{
		requester:   f,
		curData:     nil,
		curError:    nil,
		curPage:     -1,
		hasNextPage: true,
		lock:        &sync.Mutex{},
	}
	ph.Next()
	return ph
}
