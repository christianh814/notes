package main

import "fmt"

func main() {
	//maps are key value pairs. (unordered). Unlike other langages; it needs the trailing comma
	// syntax is `map[keytype]valuetype{element1:value1, element2:value2, ...}`
	x := map[string]int{
		"James": 42,
		"MoneyPenny": 24,
	}
	fmt.Println(x)
	// ask for the key and you get the value
	fmt.Println(x["James"])

	// Asks if a key exists (if you do it with println; you'll get 0 if there is no key...you may not want that...
	// ...this is the prefered method. Called "comma ok idiom"
	v, ok := x["foobar"]
	// this will print out the value...which is 0
	fmt.Println(v)
	// this prints out `false` since it doesn't exist
	fmt.Println(ok)


	// This is how you would check it in a if
	if v, ok := x["foobar"]; ok {
		fmt.Println("This should not print but the default value of a missing key is:", v)
	} else {
		fmt.Println("Key does not exist")
	}

	// Just check if it exists. throw away _ as v because I don't care about the value
	// You're saying if it's "NOT okay"...hence the !
	if _, ok := x["foobar"]; !ok {
		fmt.Println("Key does not exist (throw away)")
	}
}
