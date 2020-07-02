/*package middleware

 import (
	"fmt"
	"net/http"

	"github.com/GGP1/palo/internal/cfg"

	"golang.org/x/net/context"
)

// AdminOnly verifies if the user is an admin and gives him special permissions
func AdminOnly(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !currentUser(ctx, r) {
			http.NotFound(w, r)
			fmt.Println("YOU ARE NOT AN ADMIN!")
			return
		}
		f(w, r)
	})
}*/

/*func currentUser(ctx context.Context, r *http.Request) bool {
	if ___ {
		return true
	}
	return false
}*/
