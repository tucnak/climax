package climax

import (
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {
	app := &Application{}
	check := func(c string, f []Flag, a []string, exp Context) {
		exp.app = app

		ctx, err := app.parseContext(f, a)
		if err != nil {
			t.Errorf(`case "%s" didn't finish well:`, c)
			t.Logf(`error: %s`, err)
			return
		}

		if !reflect.DeepEqual(*ctx, exp) {
			t.Errorf(`case "%s" didn't finish well:`, c)
			t.Logf("- expected:\n%s", exp.String())
			t.Logf("- recieved:\n%s", ctx.String())
		}
	}

	mustFail := func(c string, f []Flag, a []string) {
		_, err := app.parseContext(f, a)
		if err == nil {
			t.Errorf(`invalid case "%s" resulted in valid context`, c)
		}
	}

	// PASS TESTS
	// ==========

	check("no args", []Flag{}, []string{}, *newContext(app))

	check("no flags", []Flag{}, []string{"argument", "a thing"}, Context{
		Args:        []string{"argument", "a thing"},
		NonVariable: map[string]bool{},
		Variable:    map[string]string{},
	})

	check("single non-var flag", []Flag{
		Flag{Name: "force"},
	}, []string{"--force", "hard life in ghetto"}, Context{
		Args: []string{"hard life in ghetto"},
		NonVariable: map[string]bool{
			"force": true,
		},
		Variable: map[string]string{},
	})

	check("single shortened non-var flag", []Flag{
		Flag{Name: "force", Short: "f"},
	}, []string{"-f", "hard life in ghetto"}, Context{
		Args: []string{"hard life in ghetto"},
		NonVariable: map[string]bool{
			"force": true,
		},
		Variable: map[string]string{},
	})

	check("separated variable flag", []Flag{
		Flag{Name: "filter", Variable: true},
	}, []string{"--filter", "token here"}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "token here",
		},
		NonVariable: map[string]bool{},
	})

	check("joined empty variable flag", []Flag{
		Flag{Name: "filter", Variable: true},
	}, []string{`--filter=`}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "",
		},
		NonVariable: map[string]bool{},
	})

	check("joined single-word variable flag", []Flag{
		Flag{Name: "filter", Variable: true},
	}, []string{`--filter=token`}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "token",
		},
		NonVariable: map[string]bool{},
	})

	check("joined multi-word variable flag", []Flag{
		Flag{Name: "filter", Variable: true},
	}, []string{`--filter=token here`}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "token here",
		},
		NonVariable: map[string]bool{},
	})

	check("joined single-word shortened variable flag", []Flag{
		Flag{Name: "filter", Short: "f", Variable: true},
	}, []string{`-f=token`}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "token",
		},
		NonVariable: map[string]bool{},
	})

	check("joined multi-word shortened variable flag", []Flag{
		Flag{Name: "filter", Short: "f", Variable: true},
	}, []string{`-f=token here`}, Context{
		Args: []string{},
		Variable: map[string]string{
			"filter": "token here",
		},
		NonVariable: map[string]bool{},
	})

	check("sophisticated", []Flag{
		Flag{Name: "force", Short: "f", Variable: false},
		Flag{Name: "slug", Variable: true},
	}, []string{"-f", "--slug", "dog_03", "Dog Doggson"}, Context{
		Args: []string{"Dog Doggson"},
		NonVariable: map[string]bool{
			"force": true,
		},
		Variable: map[string]string{
			"slug": "dog_03",
		},
	})

	// FAIL TESTS
	// ==========

	mustFail("setting non-var flag", []Flag{
		Flag{Name: "force", Variable: false},
	}, []string{"--force=value"})

	mustFail("missing var flag value", []Flag{
		Flag{Name: "filter", Variable: true},
	}, []string{"--filter"})
}
