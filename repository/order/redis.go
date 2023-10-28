package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmvdr-iscte/Orders/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error { // meter o context em primeiro numafunção de go que precise de um context é boa pratica
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	txn := r.Client.TxPipeline() // uma transaction serve para juntar dois acontecimentos e torna-los um atómico

	res := txn.SetNX(ctx, key, string(data), 0) // vai inserir o valor na bd sem repetir as entradas
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil { // serve para adicionar a chave a um set no redis
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) { // esta função devolve dois resultados o model.Order e erro
	key := orderIDKey(id)

	value, err := r.Client.Get(ctx, key).Result() // se existir, devolve a entrada, senão devolve um erro do redis
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist // devolve o erro e order vazia
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err) // devolve o erro e order vazia
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order) // altera a instancia original da order através da referencia
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}
	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	txn := r.Client.TxPipeline() // uma transaction serve para juntar dois acontecimentos e torna-los um atómico

	err := r.Client.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get order: %w", err) // devolve o erro
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil { // remover a key do orders set
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set %w", err) // devolve o erro
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err() // Só altera se a entrada já existir
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("get order: %w", err) // devolve o erro
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) { // como o set não está ordenadovamos receber os resultados ao calhas
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size)) // vai buscar tudo no set, de acordo com o page size

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}
	xs, err := r.Client.MGet(ctx, keys...).Result() // permite ir buscar todas as keys ao passar os ...
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, len(xs)) // criamos uma slice do tamanho de todas as order slices

	for i, x := range xs { // itera-se sobre todos os resultados e transforma-se em string
		x := x.(string)
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}
		orders[i] = order
	}
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
