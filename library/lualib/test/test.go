package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type User struct {
	Name  string
	token string
}

func (u *User) SetToken(t string) {
	u.token = t
}

func (u *User) Token() string {
	return u.token
}

const script = `
print("Hello from Lua, " .. u.token .. "!")
u:SetToken("12345")
`

func Example_basic() {
	L := lua.NewState()
	defer L.Close()

	u := &User{
		Name: "Tim",
	}
	L.SetGlobal("u", luar.New(L, u))
	if err := L.DoString(script); err != nil {
		panic(err)
	}

	fmt.Println("Lua set your token to:", u.Token())
	// Output:
	// Hello from Lua, Tim!
	// Lua set your token to: 12345
}

func main() {
	Example_basic()
}