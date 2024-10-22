package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/roquib/go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatal(err)
	}
	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock created Successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		log.Fatal("Error getting stocks")
	}
	json.NewEncoder(w).Encode(stocks)

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatal("Error getting stock")
	}
	json.NewEncoder(w).Encode(stock)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("Unable to decode the request body")
	}
	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated Successfully. Total rows/records affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}
	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stock deleted Successfully. total rows/records affected. %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := `insert into stocks(name,price,company) values ($1,$2,$3) returning stockid`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatal("Unable to execute the query", err)
	}
	fmt.Printf("Inserted a single record %v", id)
	return id

}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	var stock models.Stock
	sqlStatement := `select * from stocks where stockid=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row. %v", err)
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	var stocks []models.Stock
	sqlStatement := `select * from stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer db.Close()
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatal("Unable to scan the row: ", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()
	sqlStatement := `update stocks set name=$2, price=$3, company=$4 where stockid=$1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer db.Close()
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}
func deleteStock(id int64) int64 {
	db := CreateConnection()
	sqlStatement := `delete from stocks where stockid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer db.Close()
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected

}
