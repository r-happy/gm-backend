package main

func main() {
	router := newRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
