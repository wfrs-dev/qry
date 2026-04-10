package db

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
)

var DSNDrivers = map[string]string{
	"mysql":     "{username}:{password}@{host}:{port}/{dbname}[?param1=value1&...]",
	"postgres":  "postgres://{username}:{password}@{host}:{port}/{dbname}[?param1=value1&...]",
	"sqlserver": "sqlserver://{username}[:{password}]@{host}:{port}?database={database}[&param1=value1&...]",
}

func GetDriverInstances() map[string]Driver {
	return map[string]Driver{
		"mysql":     NewSQLStandardDriver("mysql", "mysql"),
		"postgres":  NewSQLStandardDriver("postgres", "pgx"),
		"sqlserver": NewSQLServerDriver(),
	}
}

func GetDriver(name string) (Driver, error) {
	instances := GetDriverInstances()
	if driver, ok := instances[name]; ok {
		return driver, nil
	}

	return nil, fmt.Errorf("driver `%s` no encontrado", name)
}

type QueryResult struct {
	Columns []string
	Rows    [][]any
	Summary string
}

type Driver interface {
	Connect(string) (*sql.DB, error)
	EnsureDefaultLimit(string) string
	Query(dbConn *sql.DB, query string) (*QueryResult, error)
	GetName() string
}

type SQLStandardDriver struct {
	name   string
	driver string
}

type SQLServerDriver struct{}

func NewSQLStandardDriver(name, driver string) *SQLStandardDriver {
	return &SQLStandardDriver{
		name:   name,
		driver: driver,
	}
}

func (d *SQLStandardDriver) Connect(dsn string) (*sql.DB, error) {
	return connect(d.driver, dsn)
}

func (d *SQLStandardDriver) GetName() string {
	return d.name
}

func (d *SQLStandardDriver) EnsureDefaultLimit(query string) string {
	q := strings.TrimSpace(query)

	qUpper := strings.ToUpper(q)
	if !strings.HasPrefix(qUpper, "SELECT") {
		return query
	}

	limitRegex := regexp.MustCompile(`(?i)\bLIMIT\b`)
	if limitRegex.MatchString(q) {
		return query
	}

	hasSemicolon := false
	if strings.HasSuffix(q, ";") {
		q = strings.TrimSuffix(q, ";")
		hasSemicolon = true
	}

	q = q + " LIMIT 200"
	if hasSemicolon {
		q += ";"
	}

	return q
}

func (d *SQLStandardDriver) Query(dbConn *sql.DB, query string) (*QueryResult, error) {
	modifiedQuery := d.EnsureDefaultLimit(query)

	return executeQuery(dbConn, modifiedQuery)
}

func NewSQLServerDriver() *SQLServerDriver {
	return &SQLServerDriver{}
}

func (d *SQLServerDriver) Connect(dsn string) (*sql.DB, error) {
	return connect("sqlserver", dsn)
}

func (d *SQLServerDriver) GetName() string {
	return "sqlserver"
}

func (d *SQLServerDriver) EnsureDefaultLimit(query string) string {
	q := strings.TrimSpace(query)

	qUpper := strings.ToUpper(q)
	if !strings.HasPrefix(qUpper, "SELECT") {
		return query
	}

	topRegex := regexp.MustCompile(`(?i)\bTOP\b`)
	if topRegex.MatchString(q) {
		return query
	}

	re := regexp.MustCompile(`(?i)^SELECT\s+(DISTINCT\s+)?`)
	match := re.FindString(q)
	if match != "" {
		return re.ReplaceAllString(q, match+"TOP 200 ")
	}

	return q
}

func (d *SQLServerDriver) Query(dbConn *sql.DB, query string) (*QueryResult, error) {
	modifiedQuery := d.EnsureDefaultLimit(query)

	return executeQuery(dbConn, modifiedQuery)
}

func connect(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("fallo al conectar a la base de datos: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("fallo de conectividad a la base de datos: %w", err)
	}

	return db, nil
}

func executeQuery(dbConn *sql.DB, query string) (*QueryResult, error) {
	qUpper := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(qUpper, "SELECT") ||
		strings.HasPrefix(qUpper, "WITH") ||
		strings.HasPrefix(qUpper, "SHOW") ||
		strings.HasPrefix(qUpper, "DESCRIBE") ||
		strings.HasPrefix(qUpper, "EXPLAIN")

	if isSelect {
		return executeSelect(dbConn, query)
	}

	return executeUpdate(dbConn, query)
}

func executeSelect(db *sql.DB, query string) (*QueryResult, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultRows [][]any
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		rowOut := make([]any, len(cols))
		for i := range columns {
			val := columnPointers[i].(*any)
			rowOut[i] = *val
		}
		resultRows = append(resultRows, rowOut)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	res := &QueryResult{
		Columns: cols,
		Rows:    resultRows,
		Summary: fmt.Sprintf("%d filas devueltas", len(resultRows)),
	}

	return res, nil
}

func executeUpdate(db *sql.DB, query string) (*QueryResult, error) {
	result, err := db.Exec(query)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	var summary string
	if err != nil {
		summary = "Consulta ejecutada correctamente. (Información de filas afectadas no soportada por el driver): " + err.Error()
	} else {
		summary = fmt.Sprintf("%d filas afectadas", affected)
	}

	return &QueryResult{
		Summary: summary,
	}, nil
}
