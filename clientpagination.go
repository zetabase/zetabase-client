package zetabase

import (
	"errors"
	"log"
)

type putPages struct {
	Client          *ZetabaseClient
	Keys            []string
	Data            [][]byte
	MaxBytesPerPage uint64
}

type getPages struct {
	DataKeys     []string 
	MaxItemSize  int64 
	MaxPageSize  int64
	Client       *ZetabaseClient 
	PagHandlers  []*PaginationHandler 
	PagIndex     int
}

func makePutPages(client *ZetabaseClient, keys []string, valus [][]byte, maxBytes uint64) *putPages {
	return &putPages{
		Client:          client,
		Keys:            keys,
		Data:            valus,
		MaxBytesPerPage: maxBytes,
	}
}

func (p *putPages) putAll(tblOwnerId, tblId string, overwrite bool) error {
	keyPgs, valuPgs, err := p.pagify()
	if err != nil {
		return err
	}

	// TODO- removeme
	log.Printf("putAll: have %d pages (%d objects in page 0)\n", len(keyPgs), len(keyPgs[0]))

	for i := 0; i < len(keyPgs); i++ {
		keys := keyPgs[i]
		valus := valuPgs[i]
		err := p.Client.putMultiRaw(tblOwnerId, tblId, keys, valus, overwrite)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *putPages) pagify() ([][]string, [][][]byte, error) {
	pageBytes := uint64(0)
	var curPage []string
	var curPageValus [][]byte
	var pages [][]string
	var pageValues [][][]byte

	for i := 0; i < len(p.Keys); i++ {
		dlen := uint64(len(p.Data[i]))
		if dlen > p.MaxBytesPerPage {
			return nil, nil, errors.New("IndividualObjectTooLarge")
		}
		if (pageBytes + dlen) > p.MaxBytesPerPage {
			// Start a new page
			pages = append(pages, curPage)
			pageValues = append(pageValues, curPageValus)
			curPage = []string{}
			curPageValus = [][]byte{}
			pageBytes = 0
		}
		curPage = append(curPage, p.Keys[i])
		curPageValus = append(curPageValus, p.Data[i])
		pageBytes += dlen
	}
	if len(curPage) > 0 {
		pages = append(pages, curPage)
		pageValues = append(pageValues, curPageValus)
	}
	return pages, pageValues, nil
}


func MakeGetPages(client *ZetabaseClient, dataKeys []string, maxItemSize int64) *getPages {
	return &getPages{
		DataKeys:    dataKeys,
		MaxItemSize: maxItemSize,
		Client:      client,
		MaxPageSize: int64(2000000),
		PagIndex:    0,
	}
}

func (p *getPages) breakKeys() ([][]string){
	var keyGroups [][]string 
	var itemsPerPage int64
	itemsPerPage = p.MaxPageSize/p.MaxItemSize

	var lenKeys int64
	lenKeys = int64(len(p.DataKeys))
	var i int64
	var kg []string
	for i=0; i<lenKeys; i++ {
		kg = append(kg, p.DataKeys[i])
		if int64(len(kg)) == itemsPerPage {
			keyGroups = append(keyGroups, kg)
			kg = []string{}
		} else if i == lenKeys-1 {
			keyGroups = append(keyGroups, kg)
		}
	}

	return keyGroups
}

func (p *getPages) GetAll(tableOwnerId, tableId string) {
	keyGroups := p.breakKeys()
	
	for i:=0; i<len(keyGroups); i++ {
		kg := keyGroups[i]
		pag := p.Client.Get(tableOwnerId, tableId, kg)
		p.PagHandlers = append(p.PagHandlers, pag)
	}
}

func (p *getPages) DataAll() (map[string][]byte, error) {
	data := make(map[string][]byte)
	for i:=0; i<len(p.PagHandlers); i++ {
		var curData map[string][]byte

		curPag := p.PagHandlers[i]
		curData, _ = curPag.DataAll()

		for k, v := range(curData) {
			data[k] = v
		}
	}
	return data, nil
}

func (p *getPages) KeysAll() ([]string, error) {
	var keys []string 
	for i:=0; i<len(p.PagHandlers); i++ {
		curPag := p.PagHandlers[i]
		pagKeys, _ := curPag.KeysAll()

		keys = append(keys, pagKeys...)
	}
	return keys, nil
}

func (p *getPages) Data() (map[string][]byte, error) {
	curPag := p.PagHandlers[p.PagIndex]
	return curPag.Data()
}

func (p *getPages) Keys() ([]string, error) {
	var keys []string 
	curPag := p.PagHandlers[p.PagIndex]
	pagData, _ := curPag.Data()
	for k := range(pagData) {
		keys = append(keys, k)
	}
	return keys, nil
}

func (p *getPages) Next() {
	curPag := p.PagHandlers[p.PagIndex]

	if curPag.hasNextPage {
		curPag.Next()
	} else if p.PagIndex < len(p.PagHandlers) - 1 {
		p.PagIndex ++ 
	}
}
