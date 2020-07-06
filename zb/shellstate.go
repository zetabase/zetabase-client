package main

import (
	"github.com/zetabase/zetabase-client/zbprotocol"
	"sync"
)

type ShellState struct {
	keyBuffer      []string
	keyBufferTblId *string
	tblIdHistMap   map[string]bool
	lock           *sync.Mutex
}

func NewShellState() *ShellState {
	return &ShellState{
		keyBuffer:      nil,
		keyBufferTblId: nil,
		tblIdHistMap:   map[string]bool{},
		lock:           &sync.Mutex{},
	}
}

func (s *ShellState) SetKeyBuffer(tblId string, ks []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.keyBuffer = ks
	s.keyBufferTblId = &tblId
}

func (s *ShellState) AddTableIdToHistory(tblId string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.tblIdHistMap[tblId] = true
}

func (s *ShellState) GetKeyBufferTable() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.keyBufferTblId != nil {
		return *s.keyBufferTblId
	}
	return ""
}

func (s *ShellState) GetKeyBuffer() []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.keyBuffer
}

func (s *ShellState) GetUsedTableIds() []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	var res []string
	for k, _ := range s.tblIdHistMap {
		res = append(res, k)
	}
	return res
}


func (s *ShellState) IngestTablesData(res []*zbprotocol.TableCreate) {
	if res != nil {
		for _, defn := range res {
			s.AddTableIdToHistory(defn.GetTableId())
		}
	}
}


