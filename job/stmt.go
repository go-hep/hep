package job

// Stmt represents a job options statement.
type Stmt struct {
	Type StmtType // type of the statement
	Data C        // the configuration data associated with that statement
}

// StmtType represents the type of a job-options statement.
type StmtType int

const (
	StmtCreate StmtType = iota
	StmtSetProp
)
