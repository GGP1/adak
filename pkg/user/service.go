package user

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/GGP1/adak/internal/token"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

// Users implements the users interface.
type Users struct {
	db *sqlx.DB

	orderingClient ordering.OrderingClient
	shoppingClient cart.ShoppingClient
}

// NewService returns a new users server.
func NewService(db *sqlx.DB, orderingConn, shoppingConn *grpc.ClientConn) *Users {
	return &Users{
		db:             db,
		orderingClient: ordering.NewOrderingClient(orderingConn),
		shoppingClient: cart.NewShoppingClient(shoppingConn),
	}
}

// Run starts the server.
func (u *Users) Run(port int) error {
	srv := grpc.NewServer()
	RegisterUsersServer(srv, u)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "users: failed listening on port %d", port)
	}

	log.Info().Msgf("User service listening on port %d", port)
	return srv.Serve(lis)
}

// Create creates a user.
func (u *Users) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	cartQuery := `INSERT INTO carts
	(id, counter, weight, discount, taxes, subtotal, total)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	userQuery := `INSERT INTO users
	(id, cart_id, username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := u.GetByEmail(ctx, &GetByEmailRequest{Email: req.User.Email})
	if err == nil {
		return nil, errors.New("email is already taken")
	}

	_, err = u.GetByUsername(ctx, &GetByUsernameRequest{Username: req.User.Username})
	if err == nil {
		return nil, errors.New("username is already taken")
	}

	// Non default cost blocks forever (check bcrypt issues)
	hash, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.User.Password = string(hash)

	// Create a cart for each user
	cartID := token.GenerateRunes(30)
	req.User.CartID = cartID

	new, _ := u.shoppingClient.New(ctx, &cart.NewRequest{ID: req.User.CartID})

	// Create user cart
	_, err = u.db.ExecContext(ctx, cartQuery, new.Cart.ID, new.Cart.Counter, new.Cart.Weight,
		new.Cart.Discount, new.Cart.Taxes, new.Cart.Subtotal, new.Cart.Total)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the cart")
	}

	userID := token.GenerateRunes(30)
	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	// Create user
	_, err = u.db.ExecContext(ctx, userQuery, userID, new.Cart.ID, req.User.Username, req.User.Email,
		req.User.Password, createdAt, updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the user")
	}

	return &CreateResponse{}, nil
}

// Delete permanently deletes a user from the database.
func (u *Users) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the user")
	}

	return &DeleteResponse{}, nil
}

// Get returns a list with all the users stored in the database.
func (u *Users) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var users []*ListUser

	if err := u.db.SelectContext(ctx, &users, "SELECT id, cart_id, username, email, created_at FROM users"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	list, err := getRelationships(ctx, u.db, users)
	if err != nil {
		return nil, err
	}

	return &GetResponse{Users: list}, nil
}

// GetByEmail retrieves the user requested from the database.
func (u *Users) GetByEmail(ctx context.Context, req *GetByEmailRequest) (*GetByEmailResponse, error) {
	var user ListUser
	q := "SELECT id, cart_id, email, username, password, created_at FROM users WHERE email=$1"
	if err := u.db.GetContext(ctx, &user, q, req.Email); err != nil {
		return nil, errors.Wrap(err, "couldn't find the user")
	}

	return &GetByEmailResponse{User: &user}, nil
}

// GetByID retrieves the user requested from the database.
func (u *Users) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	var (
		user    ListUser
		reviews []*review.Review
	)

	if err := u.db.GetContext(ctx, &user, "SELECT id, cart_id, username, email, created_at FROM users WHERE id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the user")
	}

	if err := u.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE user_id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	getByUserID, err := u.orderingClient.GetByUserID(ctx, &ordering.GetByUserIDRequest{UserID: req.ID})
	if err != nil {
		return nil, err
	}

	user.Orders = getByUserID.Orders

	return &GetByIDResponse{User: &user}, nil
}

// GetByUsername retrieves the user requested from the database.
func (u *Users) GetByUsername(ctx context.Context, req *GetByUsernameRequest) (*GetByUsernameResponse, error) {
	var user ListUser

	if err := u.db.GetContext(ctx, &user, "SELECT id, cart_id, username, email, created_at FROM users WHERE username=$1", req.Username); err != nil {
		return nil, errors.Wrap(err, "couldn't find the user")
	}

	return &GetByUsernameResponse{User: &user}, nil
}

// Search looks for the users that contain the value specified. (Only text fields)
func (u *Users) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if strings.ContainsAny(req.Search, ";-\\|@#~€¬<>_()[]}{¡'") {
		return nil, errors.New("invalid search")
	}

	users := []*ListUser{}
	q := `SELECT * FROM users WHERE
	to_tsvector(id || ' ' || username || ' ' || email) 
	@@ to_tsquery($1)`

	if err := u.db.SelectContext(ctx, &users, q, req.Search); err != nil {
		return nil, errors.Wrap(err, "couldn't find the users")
	}

	list, err := getRelationships(ctx, u.db, users)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{Users: list}, nil
}

// Update sets new values for an already existing user.
func (u *Users) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	_, err := u.db.ExecContext(ctx, "UPDATE users SET username=$2 WHERE id=$1", req.ID, req.User.Username)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't update the user")
	}

	return &UpdateResponse{}, nil
}

func getRelationships(ctx context.Context, db *sqlx.DB, users []*ListUser) ([]*ListUser, error) {
	ch, errCh := make(chan *ListUser), make(chan error, 1)

	for _, user := range users {
		go func(user *ListUser) {
			var (
				reviews []*review.Review
				orders  []*ordering.Order
			)

			if err := db.Select(&reviews, "SELECT * FROM reviews WHERE user_id=$1", user.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the reviews")
			}

			if err := db.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", user.ID); err != nil {
				errCh <- errors.Wrap(err, "couldn't find the orders")
			}

			user.Orders = orders
			user.Reviews = reviews

			ch <- user
		}(user)
	}

	list := make([]*ListUser, len(users))
	for i := range users {
		select {
		case user := <-ch:
			list[i] = user
		case err := <-errCh:
			return nil, err
		}
	}

	return list, nil
}
