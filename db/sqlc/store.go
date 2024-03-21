package db

import (
	"context"
	"database/sql"
	"fmt"
)

//store provide all funtions to execute db queries and data trival and transfers
type Store struct {
	*Queries
	db *sql.DB
}

//create NewStore
func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

//execTX execute a funtion within a database action
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q :=New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//CreateTodo give access to the API caling to perfrom deletion in the databse 
//CreateTodoTxParams contains the input parameters of the Createing of the data 
type CreateTodosParams struct{
	AccountID int64  `json:"account_id"`
	Title     string `json:"title"`
	Time      string `json:"time"`
	Date      string `json:"date"`
	Complete  string `json:"complete"`
}


//CreateTodoResult contains the result of the Createing of the data
type CreateTodoResult struct{
	Todo Todo `json:"todo"`
}


//CreateTodos give access to API call to perfrom and create in the databse 
//it contains title, time, date, completion of the todo event with the database storation
func (store *Store) CreateTodo(ctx context.Context, arg CreateTodosParams)(CreateTodoResult, error){
	var result CreateTodoResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		result.Todo, err = q.CreateTodo(ctx, CreateTodoParams{
			AccountID: arg.AccountID,
			Title: 		arg.Title,
			Time: 		arg.Time,
			Date: 		arg.Date,
			Complete: 	arg.Complete,
		})
		if err != nil {
			return err
		}

		
		return nil
	})
	return result, err
}


//DeleteTodo give access to API call to perfrom deletion in the databse 
//it contains title, time, date, completion of the todo event with the database removal
func (store *Store) DeleteTodo(ctx context.Context, id int64) error {
    return store.Queries.DeleteTodo(ctx, id)
}

//UpdateteTodoTxParams contains the input parameters of the Updating of the data 
type UpdateTodosParams struct{
	ID       int64  `json:"id"`
	Title    string      `json:"title"`
	Time     string      `json:"time"`
	Date     string      `json:"date"`
	Complete string `json:"complete"`
}


//UpdateTodoResult contains the result of the Updating of the data
type UpdateTodoResult struct{
	Todo Todo `json:"todo"`
}


//UpdateTodo give access to API call to perfrom update data in the databse 
//it contains title, time, date, completion of the todo event with the database update
func (store *Store) UpdateTodo(ctx context.Context, arg UpdateTodosParams)(UpdateTodoResult, error){
	var result UpdateTodoResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error
		updateTodo, err := q.UpdateTodo(ctx, UpdateTodoParams{
			ID: arg.ID,
			Title: 		arg.Title,
			Time: 		arg.Time,
			Date: 		arg.Date,
			Complete: 	arg.Complete,
		})

		if err != nil {
			return err
		}

		if updateTodo.ID == 0 {
			return err
		}
		
		result.Todo = updateTodo
		return nil
	})
	return result, err
}

//ListTodo give access to API call to perfrom list data in the databse 
// ListTodoParams contains the input parameters for listinging an owner
func (store *Store) ListTodo(ctx context.Context, params ListTodoParams) ([]Todo, error) {
    return store.Queries.ListTodo(ctx, params)
}

//GetTodoParams contains the input parameters of the Geting of the data 
type GetTodoParams struct{
	AccountID int64 `uri:"id" binding:"required,min=1"`
}

//GetTodoResult contains the result of the Geting of the data
type GetTodoResult struct{
	Todo Todo `json:"Todo"`
}


//GetTodo perfrom data transfer
//it contains title, time, date, completion of the todo event with the database storation
func (store *Store) GetTodo(ctx context.Context, arg GetTodoParams)(GetTodoResult, error){
	var result GetTodoResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		result.Todo, err = q.GetTodo(ctx, arg.AccountID)
		
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

//CreateAccountTxParams contains the input parameters of the Createing of the data 
type CreateAccountsParams struct{
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	UserName  string `json:"user_name" binding:"required,alphanum,min=8"`
	Password  string `json:"password" binding:"required"`
}

//CreateAccountResult contains the result of the Createing of the data
type CreateAccountResult struct{
	Account Account `json:"account"`
}

//CreateAccount perfrom data transfer
//it contains title, time, date, completion of the todo event with the database storation
func (store *Store) CreateAccount(ctx context.Context, arg CreateAccountsParams)(CreateAccountResult, error){
	var result CreateAccountResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		result.Account, err = q.CreateAccount(ctx, CreateAccountParams{
			FirstName: 		arg.FirstName,
			LastName: 		arg.LastName,
			UserName: 		arg.UserName,
			Password: 		arg.Password,
		});
		
		if err != nil {
			return err
		}

	
		return nil
	})
	return result, err
}


//GetAccountTxParams contains the input parameters of the Geting of the data 
type GetAccountsParams struct{
	UserName  string `json:"user_name" binding:"required,alphanum,min=8"`
	Password  string `json:"password" binding:"required"`
}

//GetAccountResult contains the result of the Geting of the data
type GetAccountResult struct{
	Account Account `json:"account"`
}


//GetAccount perfrom data transfer
//it contains title, time, date, completion of the todo event with the database storation
func (store *Store) GetAccount(ctx context.Context, arg GetAccountsParams)(GetAccountResult, error){
	var result GetAccountResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		result.Account, err = q.GetAccount(ctx, arg.UserName)
		
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

//deleteAccount 
//it contains title, time, date, completion of the todo event with the database removal
func (store *Store) DeleteAccount(ctx context.Context, id int64) error {
    return store.Queries.DeleteAccount(ctx, id)
}

// ListAccountParams contains the input parameters for listinging an owner
func (store *Store) ListAccount(ctx context.Context, params ListAccountParams) ([]Account, error) {
    return store.Queries.ListAccount(ctx, params)
}

//UpdateAccountTxParams contains the input parameters of the Updating of the data 
type UpdateAccountsParams struct{
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	UserName  string `json:"user_name" binding:"required,alphanum,min=8"`
	Password  string `json:"password" binding:"required"`
}

//UpdateTodoResult contains the result of the Updating of the data
type UpdateAccountResult struct{
	Account Account `json:"todo"`
}


//UpdateTodos perfrom data transfer
//it contains title, time, date, completion of the todo event with the database storation
func (store *Store) UpdateAccount(ctx context.Context, arg UpdateAccountsParams)(UpdateAccountResult, error){
	var result UpdateAccountResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		result.Account, err = q.UpdateAccount(ctx, UpdateAccountParams{
			FirstName: 		arg.FirstName,
			LastName: 		arg.LastName,
			UserName: 		arg.UserName,
			Password: 	arg.Password,
		})

		if err != nil {
			return err
		}
		
		return nil
	})
	return result, err
}


