package zetabase

import "github.com/zetabase/zetabase-client/zbprotocol"

type IndexedField struct {
	FieldName string
	LangCode  string
	IndexType zbprotocol.QueryOrdering
}

func NewIndexedField(fieldName string, indexTyp zbprotocol.QueryOrdering) *IndexedField {
	return &IndexedField{
		FieldName: fieldName,
		LangCode:  "",
		IndexType: indexTyp,
	}
}

func (f *IndexedField) SetLanguageCode(code string) {
	f.LangCode = code
}

func (f *IndexedField) ToProtocol() *zbprotocol.TableIndexField {
	return &zbprotocol.TableIndexField{
		Field:        f.FieldName,
		Ordering:     f.IndexType,
		LanguageCode: f.LangCode,
	}
}

func indexedFieldsToProtocol(ifs []*IndexedField) *zbprotocol.TableIndexFields {
	var arr []*zbprotocol.TableIndexField
	for _, x := range ifs {
		arr = append(arr, x.ToProtocol())
	}
	return &zbprotocol.TableIndexFields{
		Fields: arr,
	}
}
