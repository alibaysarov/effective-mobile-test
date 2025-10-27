package repository

import (
	"database/sql"
	"effective-mobile/internal/dto"
	"effective-mobile/internal/model"
	"effective-mobile/internal/storage/repository"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type subscribeRepository struct {
	db *sql.DB
}

func NewSubscribeRepository(db *sql.DB) repository.SubscribeRepository {
	return &subscribeRepository{
		db: db,
	}
}

func (r *subscribeRepository) GetOne(id string) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date FROM subscriptions WHERE id=$1 LIMIT 1`
	var firstRow model.Subscription
	if err := r.db.QueryRow(query, id).Scan(&firstRow.ID,
		&firstRow.ServiceName,
		&firstRow.Price,
		&firstRow.UserId, &firstRow.StartDate); err != nil {
		fmt.Println("error while fetching one subscription", err)
		return nil, err
	}
	return &firstRow, nil
}

func (r *subscribeRepository) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("error while deleting subscription", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("error while getting rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with id %s not found", id)
	}

	return nil
}
func (r *subscribeRepository) Create(data dto.SubscriptionDto) (*model.Subscription, error) {
	var lastId string
	query := `INSERT INTO subscriptions (service_name, user_id, price, start_date) 
              VALUES($1, $2, $3, $4) RETURNING id`

	if err := r.db.QueryRow(query, data.ServiceName, data.UserId, data.Price, data.StartDate).Scan(&lastId); err != nil {
		fmt.Println("error while creating subscription", err)
		return nil, err
	}

	return r.GetOne(lastId)
}

func (r *subscribeRepository) Update(data dto.SubscriptionDto) (*model.Subscription, error) {
	query := `UPDATE subscriptions 
              SET service_name = $1, 
                  user_id = $2, 
                  price = $3, 
                  start_date = $4 
              WHERE id = $5`

	_, err := r.db.Exec(query, data.ServiceName, data.UserId, data.Price, data.StartDate, data.ID)
	if err != nil {
		fmt.Println("error while updating subscription", err)
		return nil, err
	}

	return r.GetOne(data.ID)
}

func (r *subscribeRepository) GetAll(filter dto.Filter) (*[]model.Subscription, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id,service_name, price, user_id, start_date FROM subscriptions ")
	filterString, arguments, err := getFilterQuery(filter)
	if err != nil {
		return nil, err
	}
	queryBuilder.WriteString(filterString)

	rows, err := r.db.Query(queryBuilder.String(), arguments...)
	if err != nil {
		fmt.Println("error while fetching subscriptions ", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := []model.Subscription{}
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &subscriptions, nil
}

func (r *subscribeRepository) GetSum(filter dto.Filter) (float64, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT COUNT(price) FROM subscriptions ")
	filterString, arguments, err := getFilterQuery(filter)
	if err != nil {
		return 0, err
	}
	queryBuilder.WriteString(filterString)
	var count float64
	if err := r.db.QueryRow(queryBuilder.String(), arguments...).Scan(&count); err != nil {
		fmt.Println("error while fetching subscriptions ", err)
		return 0, err
	}

	return count, nil
}

func getFilterQuery(filter dto.Filter) (string, []any, error) {
	if reflect.ValueOf(filter).Kind() != reflect.Struct {
		return "", []any{}, errors.New("type is not a struct")
	}
	value := reflect.ValueOf(filter)
	typeOf := reflect.TypeOf(filter)
	numberOfFields := value.NumField()
	var builder strings.Builder
	count := 0
	arguments := []any{}
	for i := 0; i < numberOfFields; i++ {
		field := typeOf.Field(i)
		fieldValue := value.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" {
			formTag = field.Tag.Get("json")
		}

		formName := strings.Split(formTag, ",")[0]

		if !fieldValue.IsNil() {
			actualValue := fieldValue.Elem().Interface()
			if actualValue != nil {
				count++
				if count == 1 {
					builder.WriteString("WHERE ")
				}

				if count > 1 {
					builder.WriteString(" AND ")
				}

				if strings.HasSuffix(formName, "_from") {
					cutString, found := strings.CutSuffix(formName, "_from")
					if found {
						builder.WriteString(fmt.Sprintf("%s > $%d", cutString, count))
						arguments = append(arguments, actualValue)
						continue
					}
				}

				if strings.HasSuffix(formName, "_to") {
					cutString, found := strings.CutSuffix(formName, "_to")
					if found {
						builder.WriteString(fmt.Sprintf("%s < $%d", cutString, count))
						continue
					}
				}

				builder.WriteString(fmt.Sprintf("%s=$%d", formName, count))
				arguments = append(arguments, actualValue)
			}
		}
	}
	return builder.String(), arguments, nil
}
