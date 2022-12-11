package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"honnef.co/go/tools/staticcheck"
)

// main использует несколько анализаторов.
//
// OsExitAnalyzer - кастомный анализатор, который проверяет наличие os.Exit в main пакете.
//
// printf.Analyzer - проверяет на соответствие спецификаторы шаблона и типы аргументов.
//
// shadow.Analyzer - помогает найти затенённые переменные.
//
// structtag.Analyzer - проверяет теги полей структур на соответствие reflect.StructTag.Get.
//
// shift.Analyzer - проверяет сдвиги, которые превышают ширину целого числа.
//
// tests.Analyzer - проверяет распространенные ошибки использования тестов и примеров.
func main() {
	analyzers := []*analysis.Analyzer{
		OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		tests.Analyzer,
	}

	for _, sca := range staticcheck.Analyzers {
		if strings.Index(sca.Analyzer.Name, "SA") == 0 {
			analyzers = append(analyzers, sca.Analyzer)
		}
	}

	multichecker.Main(analyzers...)
}
