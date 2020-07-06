package zetabase

import (
	"github.com/zetabase/zetabase-client/zbprotocol"
	"strings"
)

type PermConstraint struct {
	Field    string
	ReqValue string
}

type PermEntry struct {
	Level        zbprotocol.PermissionLevel
	AudienceType zbprotocol.PermissionAudienceType
	AudienceId   string
	Constraints  []*PermConstraint
}

func NewPermissionConstraint(field, reqValue string) *PermConstraint {
	return &PermConstraint{
		Field:    field,
		ReqValue: reqValue,
	}
}

func NewPermConstraintUserId(field string) *PermConstraint {
	return &PermConstraint{
		Field:    field,
		ReqValue: "@uid",
	}
}

func NewPermissionEntry(level zbprotocol.PermissionLevel, typ zbprotocol.PermissionAudienceType, audId string) *PermEntry {
	return &PermEntry{
		Level:        level,
		AudienceType: typ,
		AudienceId:   audId,
		Constraints:  nil,
	}
}

func (p *PermEntry) AddConstraint(c *PermConstraint) {
	p.Constraints = append(p.Constraints, c)
}

func toFieldConstraint(uid, tblId string, cs *PermConstraint) *zbprotocol.PermissionConstraint {
	fTyp := zbprotocol.FieldConstraintValueType_CONSTANT
	fVal := cs.ReqValue
	if strings.ToLower(cs.ReqValue) == "@uid" {
		fTyp = zbprotocol.FieldConstraintValueType_UID
		fVal = ""
	}
	return &zbprotocol.PermissionConstraint{
		ConstraintType: zbprotocol.PermissionConstraintType_FIELD,
		FieldConstraint: &zbprotocol.FieldConstraint{
			ConstraintType: zbprotocol.FieldConstraintType_EQUALS_VALUE,
			FieldKey:       cs.Field,
			ValueType:      fTyp,
			RequiredValue:  fVal,
		},
	}
}

func toFieldConstraints(uid, tblId string, cs []*PermConstraint) []*zbprotocol.PermissionConstraint {
	var fcs []*zbprotocol.PermissionConstraint
	for _, x := range cs {
		fcs = append(fcs, toFieldConstraint(uid, tblId, x))
	}
	return fcs
}

func (p *PermEntry) ToProtocol(uid, tblId string) *zbprotocol.PermissionsEntry {
	return &zbprotocol.PermissionsEntry{
		Id:           uid,
		TableId:      tblId,
		AudienceType: p.AudienceType,
		AudienceId:   p.AudienceId,
		Level:        p.Level,
		Nonce:        0,
		Credential:   nil,
		Constraints:  toFieldConstraints(uid, tblId, p.Constraints),
	}
}
