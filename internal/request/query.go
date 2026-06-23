package request

import (
	"fmt"
	"net/http"
	"strconv"
)

func PositiveIntQuery(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return 0, fmt.Errorf("%s cannot be empty", key)
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be integer", key)
	}

	if id <= 0 {
		return 0, fmt.Errorf("%s must be positive", key)
	}

	return id, nil
}
