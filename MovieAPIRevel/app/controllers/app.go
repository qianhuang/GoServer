package controllers

import (
	"net/http"
	"strconv"

	"github.com/qianhuang/GoServer/MovieAPIRevel/app"
	"github.com/revel/revel"

	"log"
)

// App balabala
type App struct {
	*revel.Controller
}

// Index balabala
func (c App) Index() revel.Result {
	return c.Render()
}

//ProductResource balabala
type ProductResource struct {
	ProductNo int     `json:"id"`
	Name      string  `json:"product_name"`
	Price     float32 `json:"product_price"`
}

//GetProduct balabala
func (c App) GetProduct() revel.Result {
	var product ProductResource
	id := c.Params.Route.Get("product-id")

	product.ProductNo, _ = strconv.Atoi(id)
	product.Name = "default"
	product.Price = 1.99

	// use this ID to query from database and fill train table....
	rows, err := app.DB.Query("select product_no, name, price from products where product_no =  ?", product.ProductNo)

	revel.INFO.Println(err)

	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.ProductNo, &product.Name, &product.Price)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(product.ProductNo, product.Name, product.Price)

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	//product.ProductNo, _ = strconv.Atoi(id)
	//product.Name = "Rice"
	//product.Price = 5.99

	c.Response.Status = http.StatusOK
	return c.RenderJSON(product)
}

//CreateProduct balabala
func (c App) CreateProduct() revel.Result {
	var product ProductResource
	c.Params.BindJSON(&product)

	//insert into table ...
	stmt, err := app.DB.Prepare("INSERT INTO products VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(product.ProductNo, product.Name, product.Price)
	if err != nil {
		log.Fatal(err)
	}
	/*
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
	*/
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		revel.INFO.Println("\n row affected: ", rowCnt)
	}

	//product.ProductNo = 2
	c.Response.Status = http.StatusCreated
	return c.RenderJSON(product)

}

//RemoveProduct balabala
func (c App) RemoveProduct() revel.Result {
	id := c.Params.Route.Get("product-id")
	prodID, _ := strconv.Atoi(id)

	// Use ID to delete record from product table....
	stmt, err := app.DB.Prepare("DELETE FROM products where product_no = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(prodID)
	if err != nil {
		log.Fatal(err)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	} else {
		revel.INFO.Println("\n row affected: ", rowCnt)
	}

	log.Println("Successfully deleted the resource:", id)
	c.Response.Status = http.StatusOK
	return c.RenderText("\nDeleted!!!\n")
}
