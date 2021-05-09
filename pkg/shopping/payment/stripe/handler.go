package stripe

import (
	"net/http"

	"github.com/GGP1/adak/internal/response"

	"github.com/go-chi/chi/v5"
)

// Handler manages stripe endpoints.
type Handler struct{}

// GetBalance responds with the account balance.
func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balance, err := GetBalance()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, balance)
	}
}

// GetEvent looks for the details of the events.
func (h *Handler) GetEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := chi.URLParam(r, "event")

		event, err := GetEvent(e)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, event)
	}
}

// GetTxBalance responds with the transaction balance.
func (h *Handler) GetTxBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txID := chi.URLParam(r, "txID")

		tx, err := GetTxBalance(txID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		response.JSON(w, http.StatusOK, tx)
	}
}

// ListEvents retrieves a list of all the stripe events within the last 30 days.
func (h *Handler) ListEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventList := ListEvents()

		response.JSON(w, http.StatusOK, eventList)
	}
}

// ListTxs responds with a list of stripe transactions.
func (h *Handler) ListTxs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txList := ListTxs()

		response.JSON(w, http.StatusOK, txList)
	}
}
