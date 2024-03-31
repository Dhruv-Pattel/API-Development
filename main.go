package main

func main() {
	app := App{}
	app.initialize(dbuser, dbpwd, dbname)
	app.run("localhost:10000")
}
