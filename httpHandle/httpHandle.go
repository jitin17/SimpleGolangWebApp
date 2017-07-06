package httpHandle

import(
"net/http"
"fmt"
"database/sql"
_ "github.com/go-sql-driver/mysql"
"golang.org/x/crypto/bcrypt"
"github.com/gorilla/mux"
"strconv"
)

var IndexHtml []byte
var Login []byte
var db *sql.DB 
var err error

func Handler(){

fmt.Println("listening on 5000")

 db, err = sql.Open("mysql", "root:Admin123@/hello")
    if err != nil {
        panic(err.Error())    
    }
    // sql.DB should be long lived "defer" closes it once this function ends
    defer db.Close()

    // Test the connection to the database
    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

router:=mux.NewRouter()
router.HandleFunc("/add",add)
router.HandleFunc("/login",loginPage)
router.HandleFunc("/signup",signupPage)
http.ListenAndServe(":5000",router)


}

func add(res http.ResponseWriter, req *http.Request){

   if req.Method != "POST" {
		http.ServeFile(res, req, "./template/add.html")
		return
	}



num1 := req.FormValue("num1")
num2 := req.FormValue("num2")

numa, _:= strconv.Atoi(num1)
numb, _ := strconv.Atoi(num2)

var num3 int
num3=numa+numb

fmt.Println("value",num3)

}

func loginPage(res http.ResponseWriter, req *http.Request) {
  
	if req.Method != "POST" {
		http.ServeFile(res, req, "./template/login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
        
	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}



func signupPage(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.ServeFile(res, req, "./template/signup.html")
        return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

    err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

    switch {
    case err == sql.ErrNoRows:
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)    
            return
        } 

        _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)    
            return
        }

        res.Write([]byte("User created!"))
        return
    case err != nil: 
        http.Error(res, "Server error, unable to create your account.", 500)    
        return
    default: 
        http.Redirect(res, req, "/", 301)
    }
}




/*
func IndexHandler(w http.ResponseWriter,r *http.Request){
 w.Write(IndexHtml)

}

func Page2(w http.ResponseWriter,r *http.Request){
r.ParseForm()
 log.Println("creatring redirect",r.Form)
 w.Write(Login)

}





func connectToDB(){

  

}


func init(){
 
IndexHtml,_= ioutil.ReadFile("./template/index.html") 

Login,_= ioutil.ReadFile("./template/page2.html")
 
}*/
