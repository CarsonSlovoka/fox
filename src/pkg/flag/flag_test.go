package flag_test

import (
	"encoding/json"
	"flag"
	"fmt"
	flag2 "github.com/CarsonSlovoka/fox/pkg/flag"
	"testing"
)

// 這是一種比較囉嗦的方式，自己寫Init, Parse, MainFunc，實際上已經有寫好預設的動作，如果不需要客製化，可以省略，請參考TestCmd2
func TestCmd1(t *testing.T) {
	cmdTest := &flag2.Command{
		FlagSet: flag.NewFlagSet("test", flag.ContinueOnError),
		Fields: map[string][]flag2.CmdField{
			"string": {
				{"new", "", "Create Room"},
				{"delete", "", "Delete Room"},
				{"query", "", "Query Room"},
			},
			"bool": {
				{"list", false, "List Room"},
				{"private", true, "is private?"},
			},
		}}

	cmdTest.Init = func() {
		for typeName, cmdFields := range cmdTest.Fields {
			switch typeName {
			case "string":
				for _, curCmdField := range cmdFields {
					cmdTest.FlagSet.String(curCmdField.Name, fmt.Sprintf("%v", curCmdField.DefaultValue), curCmdField.Usage)
				}
			case "int":
				for _, curCmdField := range cmdFields {
					cmdTest.FlagSet.Int(curCmdField.Name, curCmdField.DefaultValue.(int), curCmdField.Usage)
				}
			case "bool":
				for _, curCmdField := range cmdFields {
					cmdTest.FlagSet.Bool(curCmdField.Name, curCmdField.DefaultValue.(bool), curCmdField.Usage)
				}
			}
		}
	}

	cmdTest.Parse = func(args []string, reset bool) error {
		if reset {
			for _, cmdFields := range cmdTest.Fields {
				for _, curCmdField := range cmdFields {
					field := cmdTest.FlagSet.Lookup(curCmdField.Name)
					if field != nil {
						if err := field.Value.Set(fmt.Sprintf("%v", curCmdField.DefaultValue)); err != nil {
							return err
						}
					}
				}
			}
		}
		if err := cmdTest.FlagSet.Parse(args[1:]); err != nil {
			return err
		}
		return nil
	}
	cmdTest.MainFunc = func(args []string) error {
		if err := cmdTest.Parse(args, true); err != nil {
			return err
		}
		return nil
	}

	testData := [][]string{
		{"myCmd",
			"-new", "myItem",
			"-query", "qItem",
			"-private", "true",
		},

		{"myCmd",
			"-delete", "dd",
			"-list", "true",
		},
	}
	cmdTest.Init()

	for _, args := range testData {
		if err := cmdTest.MainFunc(args); err != nil {
			t.FailNow()
		}
		fieldNew := cmdTest.Lookup("new")
		val := fieldNew.Value.(flag.Getter) // Field.Value只有String和Set的方法，而Getter也相容Field.Value只是多實現了Set方法，官方的Value除了實現了Set, String以外還多了Get，所以才可以這樣用，如果您有自定義Value要這樣用，就要注意有沒有再額外定義Get
		fmt.Printf("new:%s\n", val.Get().(string))

		fieldQuery := cmdTest.Lookup("query")
		fmt.Printf("query:%s\n", fieldQuery.Value.(flag.Getter).Get().(string))

		fieldDelete := cmdTest.Lookup("delete")
		fmt.Printf("delete:%s\n", fieldDelete.Value.(flag.Getter).Get().(string))

		fieldPrivate := cmdTest.Lookup("private")
		fmt.Printf("private:%t\n", fieldPrivate.Value.(flag.Getter).Get().(bool))

		fieldList := cmdTest.Lookup("list")
		fmt.Printf("list:%t\n", fieldList.Value.(flag.Getter).Get().(bool))
	}

	if err := cmdTest.MainFunc([]string{"myCmd",
		"-new", "myItem",
		"-query", "qItem",
	}); err != nil {
		t.FailNow()
	}
	if (cmdTest.Lookup("new")).Value.(flag.Getter).Get().(string) != "myItem" {
		t.FailNow()
	}
	if (cmdTest.Lookup("private")).Value.(flag.Getter).Get().(bool) != true {
		t.FailNow()
	}

	if err := cmdTest.MainFunc([]string{"myCmd",
		"-new", "myItem",
		"-private=false", // 注意bool的設定不能分開 "-private", "false" 是錯的
	}); err != nil {
		t.FailNow()
	}
	if (cmdTest.Lookup("private")).Value.(flag.Getter).Get().(bool) != false {
		t.FailNow()
	}
}

func TestCmd2(t *testing.T) {
	cmd := flag2.NewCommand(
		flag.NewFlagSet("echo", flag.ContinueOnError),
		map[string][]flag2.CmdField{
			"string": {
				{"msg", "No msg", "echo a message to you"},
			},
			"int": {
				{"maxCount", 5, "maxCount:int"},
			},
			"bool": {
				{"debug", false, "debug:bool"},
				{"active", true, "debug:bool"},
			},
		})

	type MyStruct struct {
		Msg      string
		MaxCount int
		Debug    bool
		Active   bool
	}

	for testIdx, curTestData := range []struct {
		args     []string
		expected MyStruct
	}{
		{[]string{"echo", "-msg", "hello world", "-debug=true"}, MyStruct{"hello world", 5, true, true}},
		{[]string{"echo", "-maxCount", "30"}, MyStruct{"No msg", 30, false, true}},
		{[]string{"echo", "-debug", "-active"}, MyStruct{"No msg", 5, true, true}},
		{[]string{"echo", "-debug=true", "-active=false"}, MyStruct{"No msg", 5, true, false}},
	} {

		if err := cmd.Parse(curTestData.args[1:], true); err != nil {
			t.FailNow()
		}
		myObj := &MyStruct{}
		myObj.Msg = (cmd.Lookup("msg")).Value.(flag.Getter).Get().(string)
		myObj.MaxCount = (cmd.Lookup("maxCount")).Value.(flag.Getter).Get().(int)
		myObj.Debug = (cmd.Lookup("debug")).Value.(flag.Getter).Get().(bool)
		myObj.Active = (cmd.Lookup("active")).Value.(flag.Getter).Get().(bool)
		expectedByte, _ := json.MarshalIndent(curTestData.expected, "", "")
		actualByte, _ := json.MarshalIndent(myObj, "", "")
		if string(expectedByte) != string(actualByte) {
			t.Fatalf("idx:%d: Expected:\n%s\nActual\n%s\n", testIdx, expectedByte, actualByte)
		}
	}
}

func TestDefaultParse(t *testing.T) {
	cmd := flag2.NewCommand(
		flag.NewFlagSet("echo", flag.ContinueOnError),
		map[string][]flag2.CmdField{
			"int": {
				{"age", 18, "usage"},
			},
		})


	if err := cmd.Parse([]string{"-notExistVar", "val"}, false); err == nil {
		t.FailNow()
	}


	if err := cmd.Parse([]string{"-age", "36"}, true); err != nil {
		t.FailNow()
	}
	if (cmd.Lookup("age")).Value.(flag.Getter).Get().(int) != 36 {
		t.FailNow()
	}

	cmd.Fields["int"][0].DefaultValue = "foo"
	err := cmd.Parse([]string{"-age", "36"}, true)
	if err == nil {
		t.FailNow()
	}
	// t.Logf(err.Error())
}
