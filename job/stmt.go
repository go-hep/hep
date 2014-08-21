package job

import (
	"fmt"
)

// Stmt represents a job options statement.
type Stmt struct {
	Type StmtType // type of the statement
	Data C        // the configuration data associated with that statement
}

// StmtType represents the type of a job-options statement.
type StmtType int

// String returns the string representation of a StmtType
func (stmt StmtType) String() string {
	switch stmt {
	case StmtCreate:
		return "CREATE"
	case StmtSetProp:
		return "SetProp"
	}
	panic(fmt.Errorf("fwk: invalid StmtType value (%d)", int(stmt)))
}

const (
	StmtCreate StmtType = iota
	StmtSetProp
)
