package sqlparser

import (
	"fmt"
	"log"

	"github.com/xwb1989/sqlparser"
)

/*
	sqlparser represents sql query string in AST (Abstract Syntax Tree)
	ASTs are used to represent structure of program/code, in this case
	a sql query
	It is "abstract" in the sense it does not capture all the details
	of the structure
*/

func ParseQuery(sql string) {
	// we can parse sql string to sqlparser.Statement type
	var stmt sqlparser.Statement
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatalf("error parsing sql: %v", err)
	}

	var selectStmt *sqlparser.Select
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		log.Fatal("error converting to Select statement")
	}

	fmt.Println("SelectExprs length: ", len(selectStmt.SelectExprs))
	fmt.Println("From: ", selectStmt.From[0])

	fmt.Println("GroupBy: ", selectStmt.GroupBy)
	fmt.Println("Having: ", selectStmt.Having)
	fmt.Println("Limit: ", selectStmt.Limit)
	fmt.Println("OrderBy: ", selectStmt.OrderBy)
	fmt.Println("Where: ", *selectStmt.Where)

	parseSelectExp(selectStmt)
	parseTableName(selectStmt)
	parseOrderBy(selectStmt)
	parseWhere(selectStmt)

	createComparisonExpr()
	createComplexConparisonExpr()
}

func parseSelectExp(selectStmt *sqlparser.Select) error {
	if len(selectStmt.SelectExprs) < 1 {
		return fmt.Errorf("invalid select statement")
	}

	for _, expr := range selectStmt.SelectExprs {
		aliasExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			return fmt.Errorf("error parsing select statement")
		}

		colName, ok := aliasExpr.Expr.(*sqlparser.ColName)
		if !ok {
			return fmt.Errorf("error parsing column name")
		}
		fmt.Println("columnName: ", colName.Name.String())
	}

	return nil
}

func parseTableName(selectStmt *sqlparser.Select) error {
	tableExpr, ok := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return fmt.Errorf("error parsing table name")
	}

	tableName := sqlparser.String(tableExpr.Expr)
	fmt.Println("tableName: ", tableName)
	return nil
}

func parseOrderBy(selectStmt *sqlparser.Select) {
	for _, order := range selectStmt.OrderBy {
		fmt.Println("orderBy field: ", sqlparser.String(order.Expr))
		fmt.Println("orderBy direction: ", order.Direction)
	}
}

func parseWhere(selectStmt *sqlparser.Select) {
	if selectStmt.Where == nil {
		return
	}
	fmt.Println("where expr: ", sqlparser.String(selectStmt.Where.Expr))
	fmt.Println("where type: ", selectStmt.Where.Type)
}

func createComparisonExpr() {
	expr := sqlparser.ComparisonExpr{
		Left: &sqlparser.ColName{
			Name: sqlparser.NewColIdent("name"),
		},
		Operator: sqlparser.EqualStr,
		Right:    sqlparser.NewStrVal([]byte("sam")),
	}

	buf := sqlparser.NewTrackedBuffer(nil)
	expr.Format(buf)
	fmt.Println("comparison expr: ", buf.String())
}

func createComplexConparisonExpr() {
	expr := sqlparser.OrExpr{
		Left: &sqlparser.ComparisonExpr{
			Left: &sqlparser.ColName{
				Name: sqlparser.NewColIdent("name"),
			},
			Operator: sqlparser.EqualStr,
			Right:    sqlparser.NewStrVal([]byte("sam")),
		},
		Right: &sqlparser.ComparisonExpr{
			Left: &sqlparser.ColName{
				Name: sqlparser.NewColIdent("email"),
			},
			Operator: sqlparser.EqualStr,
			Right:    sqlparser.NewStrVal([]byte("sam@gmail.com")),
		},
	}

	buf := sqlparser.NewTrackedBuffer(nil)
	expr.Format(buf)
	fmt.Println("complex comparison expr: ", buf.String())
}
