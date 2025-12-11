package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func GetPostgresConnection () *sql.DB {


	connstr := os.Getenv("CONN_STR")


    postgresdb,driveropenerr := sql.Open("postgres",connstr)

	if driveropenerr != nil {
          
		  log.Println("Open driver connection error",driveropenerr)  
	}  

	// connection pooling 
    postgresdb.SetMaxIdleConns(50)
	postgresdb.SetConnMaxIdleTime(100)
	postgresdb.SetMaxIdleConns(20)
	
	return postgresdb


}






