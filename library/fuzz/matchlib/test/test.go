package main

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	lua "github.com/yuin/gopher-lua"
	"gowebfuzz/library/lualib"
	"log"
	"net/http"
	"strings"
)

func isIteratorDone(mp *map[int][]int) bool {
	if (*mp)[0][0] == (*mp)[0][1] {
		return true
	}
	return false
}
func iterator() {
	mi := map[int][]int {}

	usernames := []string{
		"username1",
		"username2",
		"username3",
	}

	passwords := []string{
		"password1",
		"password2",
		"password3",
	}

	passwords2 := []string{
		"password21",
		"password22",
		"password23",
	}

	mi[0] = []int{len(usernames),0}
	mi[1] = []int{len(passwords),0}
	mi[2] = []int{len(passwords2),0}
	for ;; {
		if isIteratorDone(&mi) {
			break
		}
		fmt.Println(usernames[mi[0][1]],":",passwords[mi[1][1]],":",passwords2[mi[2][1]])

		lastKey:=len(mi)-1
		mi[lastKey][1]++
		index := lastKey
		for ;index>0;index-- {
			if mi[index][1] == mi[index][0] {
				mi[index][1] = 0
				mi[index-1][1]++
			}
		}

	}
}

func replaceBytes(bytesToBeReplaced *[]byte,bytesToReplace []byte,start int,end int) *[]byte {
	tmp := (*bytesToBeReplaced)[end:]
	*bytesToBeReplaced = append((*bytesToBeReplaced)[:start],bytesToReplace...)
	*bytesToBeReplaced = append(*bytesToBeReplaced,tmp...)
	return bytesToBeReplaced
}

func testCel() {
	d := cel.Declarations(decls.NewVar("ame", decls.String))
	env, err := cel.NewEnv(d)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	ast, iss := env.Compile(`"Hello world! I'm " + ame + "."`)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalln(err)
	}

		out, _, err := prg.Eval(map[string]interface{}{
			"ame": "CEL",
		})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(out)

}

func testLua() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
}

type Person struct {
	Name string
}

const luaPersonTypeName = "person"

// Registers my person type to given L.
func registerPersonType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaPersonTypeName)
	L.SetGlobal("person", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newPerson))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), personMethods))
}

// Constructor
func newPerson(L *lua.LState) int {
	person := &Person{L.CheckString(1)}
	ud := L.NewUserData()
	ud.Value = person
	L.SetMetatable(ud, L.GetTypeMetatable(luaPersonTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkPerson(L *lua.LState) *Person {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Person); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}

var personMethods = map[string]lua.LGFunction{
	"name": personGetSetName,
}

// Getter and setter for the Person#Name
func personGetSetName(L *lua.LState) int {
	p := checkPerson(L)
	if L.GetTop() == 2 {
		p.Name = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Name))
	return 1
}


func run(s string) {
	L := lua.NewState()
	defer L.Close()
	registerPersonType(L)
	L.SetGlobal("response",lua.LString(string(s)))
	if err := L.DoString(`
		p = person.new("Steeve")
		print(p:name()) -- "Steeve"
		p:name("Alice")
		print(p:name()) -- "Alice"
		print(response)
    `); err != nil {
		panic(err)
	}

}

func testMatchResponse() {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	request,_ := http.NewRequest("GET","http://www.baidu.com",strings.NewReader(""))
	client := http.Client{}
	res,_ := client.Do(request)
	lualib.RegisterRequestTable(L,request)
	lualib.RegisterResponseTable(L,res)
	if err := L.DoFile("./test.lua"); err != nil {
		panic(err)
	}
	i := L.GetGlobal("result").(lua.LNumber)
	fmt.Println(i)
}
func main() {
	testMatchResponse()
}
