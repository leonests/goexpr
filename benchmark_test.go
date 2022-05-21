package rulengine

import "testing"

func BenchmarkRulengine(b *testing.B) {
	benchmarks := []ParseAstTest{
		{
			Name:  "const",
			Input: "1",
		},
		{
			Name:  "literal",
			Input: "(2) > (1)",
		},
		{
			Name:  "modifier",
			Input: "(2) + (2) == (4)",
		},
		{
			Name:  "single parameter",
			Input: "param_key",
			Params: map[string]interface{}{
				"param_key": "param_val",
			},
		},
		{
			Name:  "parameter",
			Input: "param1 < param2",
			Params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
			},
		},
	}

	for _, benchmark := range benchmarks {
		expr, err := NewExpr(benchmark.Input)
		if err != nil {
			b.Fatal(err)
		}

		_, err = expr.Eval(benchmark.Params)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(benchmark.Name+"_eval", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				expr.Eval(nil)
			}
		})

		b.Run(benchmark.Name+"_parsing", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewExpr(benchmark.Input)
			}
		})
	}
}
