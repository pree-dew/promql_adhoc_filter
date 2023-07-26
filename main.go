package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/VictoriaMetrics/metricsql"
)

func addLabelsToExpr(existingLabels [][]metricsql.LabelFilter, byLabels []string, labels map[string]string) [][]metricsql.LabelFilter {
	for k, val := range labels {
		for _, by := range byLabels {
			if by == k {
				existingLabels[0] = append(existingLabels[0], metricsql.LabelFilter{Label: k, Value: val})
			}
		}
	}

	return existingLabels
}

func addLabelsToPromql(expr interface{}, labels map[string]string) {
	switch v := expr.(type) {
	case *metricsql.AggrFuncExpr:
		byLabels := v.Modifier.Args
		for _, a := range v.Args {
			switch arg := a.(type) {
			case *metricsql.FuncExpr:
				re := arg.Args[0].(*metricsql.RollupExpr)
				me := re.Expr.(*metricsql.MetricExpr)
				me.LabelFilterss = addLabelsToExpr(me.LabelFilterss, byLabels, labels)
				re.Expr = me
				v.Args[0] = re

			case *metricsql.MetricExpr:
				arg.LabelFilterss = addLabelsToExpr(arg.LabelFilterss, byLabels, labels)
				v.Args[0] = arg
			case metricsql.Expr:
				addLabelsToPromql(a, labels)
			}
		}

		expr = v

	case *metricsql.BinaryOpExpr:
		addLabelsToPromql(v.Left, labels)
		addLabelsToPromql(v.Right, labels)
	case *metricsql.NumberExpr:
		fmt.Printf("Number: %s \n", v.AppendString(nil))
	default:
		fmt.Printf("unknown expr.: %+v", reflect.TypeOf(expr))
	}
}

func main() {
	expr, err := metricsql.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("error in parsing : %+v\n", err)
		os.Exit(1)
	}

	labels := map[string]string{}
	if len(os.Args) > 2 {
		if err = json.Unmarshal([]byte(os.Args[2]), &labels); err != nil {
			fmt.Printf("error in unmarshaling: %+v\n", err)
			os.Exit(1)
		}
	}

	addLabelsToPromql(expr, labels)

	fmt.Printf("parsed expr: %s\n", expr.AppendString(nil))
}
