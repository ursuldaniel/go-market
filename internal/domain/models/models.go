package models

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
}

type Purchase struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	ProductId int    `json:"productId"`
	Quantity  int    `json:"quantity"`
	Timestamp string `json:"timestamp"`
}
