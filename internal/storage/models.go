// описание структур базы данных
package models

type User struct {
    ID 			int64  	`db:"id" 		  json:"id"`
    Name  		string 	`db:"name"		  json:"name"`
    Currency 	string 	`db:"currency"  json:"currency"`
    Quantity 	int 	`db:"quantity"  json:"quantity"`
}
