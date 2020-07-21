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
