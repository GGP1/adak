package review

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/GGP1/adak/internal/token"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Reviews implements the reviews service.
type Reviews struct {
	db *sqlx.DB
}

// NewService returns a new reviews server.
func NewService(db *sqlx.DB) *Reviews {
	return &Reviews{
		db: db,
	}
}

// Run starts the server.
func (r *Reviews) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrapf(err, "reviews: failed listening on port %d", port)
	}

	srv := grpc.NewServer()
	RegisterReviewsServer(srv, r)

	return srv.Serve(lis)
}

// Create creates a review.
func (r *Reviews) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	q := `INSERT INTO reviews(id, stars, comment, user_id, product_id, shop_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	id := token.GenerateRunes(30)
	req.Review.CreatedAt.Seconds = time.Now().Unix()

	_, err := r.db.ExecContext(ctx, q, id, req.Review.Stars, req.Review.Comment, req.UserID, req.Review.ProductID, req.Review.ShopID, req.Review.CreatedAt, req.Review.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create the review")
	}

	return nil, nil
}

// Delete permanently deletes a review from the database.
func (r *Reviews) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	_, err := r.db.ExecContext(ctx, "DELETE FROM reviews WHERE id=$1", req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't delete the review")
	}

	return nil, nil
}

// Get returns a list with all the reviews stored in the database.
func (r *Reviews) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	var reviews []*Review

	if err := r.db.SelectContext(ctx, &reviews, "SELECT * FROM reviews"); err != nil {
		return nil, errors.Wrap(err, "couldn't find the reviews")
	}

	return &GetResponse{Reviews: reviews}, nil
}

// GetByID retrieves the review requested from the database.
func (r *Reviews) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	var review *Review

	if err := r.db.GetContext(ctx, &review, "SELECT * FROM reviews WHERE id=$1", req.ID); err != nil {
		return nil, errors.Wrap(err, "couldn't find the review")
	}

	return &GetByIDResponse{Reviews: review}, nil
}
