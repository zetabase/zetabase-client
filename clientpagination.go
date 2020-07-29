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
	KeyGroups     [][]string 
	Client        *ZetabaseClient 
	ItemsPerPage  int64
	KeyIndex      int
	TableOwnerId  string 
	TableId       string 
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

func MakeGetPages(client *ZetabaseClient, dataKeys []string, maxItemSize int64, tableOwnerId, tableId string) *getPages {
	maxPageSize := int64(2000000)
	itemsPerPage := maxPageSize/maxItemSize
	kgs := breakKeys(dataKeys, itemsPerPage)

	return &getPages{
		KeyGroups:    kgs,
		Client:       client,
		ItemsPerPage: itemsPerPage,
		KeyIndex:     0,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
	}
}

func breakKeys(keys []string, itemsPerPage int64) ([][]string){
	var keyGroups [][]string 
	var i int64
	var kg []string

	lenKeys := int64(len(keys))

	for i=0; i<lenKeys; i++ {
		kg = append(kg, keys[i])
		if int64(len(kg)) == itemsPerPage {
			keyGroups = append(keyGroups, kg)
			kg = []string{}
		} else if i == lenKeys-1 {
			keyGroups = append(keyGroups, kg)
		}
	}

	return keyGroups
}

func addData(data map[string][]byte, curData map[string][]byte) {
	for k, v := range(curData) {
		data[k] = v
	}
}

func (p *getPages) getCurPag() *PaginationHandler{
	pag := p.Client.getPag(p.TableOwnerId, p.TableId, p.KeyGroups[p.KeyIndex])
	return pag 
}

func (p *getPages) DataAll() (map[string][]byte, error) {
	dataAll := make(map[string][]byte)

	for p.KeyIndex < len(p.KeyGroups) {
		var curData map[string][]byte

		curPag := p.getCurPag()
		curData, _ = curPag.DataAll()

		addData(dataAll, curData)

		p.KeyIndex ++ 
	}
	return dataAll, nil
}

func (p *getPages) KeysAll() ([]string, error) {
	var keys []string 

	for p.KeyIndex < len(p.KeyGroups) {
		curPag := p.getCurPag()
		pagKeys, _ := curPag.KeysAll()

		keys = append(keys, pagKeys...)
		p.KeyIndex ++
	}
	return keys, nil
}

func (p *getPages) Data() (map[string][]byte, error) {
	curPag := p.getCurPag()
	return curPag.Data()
}

func (p *getPages) Keys() ([]string, error) {
	curPag := p.getCurPag()
	return curPag.Keys()
}

func (p *getPages) GetFirstNPages(numPages int) (map[string][]byte, error) {
	dataAll := make(map[string][]byte)

	for p.KeyIndex < numPages {
		var curData map[string][]byte

		curPag := p.getCurPag()
		curData, _ = curPag.DataAll()

		addData(dataAll, curData)

		p.KeyIndex ++ 
	}
	return dataAll, nil
}

func (p *getPages) Next() {
	if p.KeyIndex < len(p.KeyGroups) - 1 {
		p.KeyIndex ++ 
	}
}
