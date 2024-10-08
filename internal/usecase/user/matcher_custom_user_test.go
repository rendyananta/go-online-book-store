package user

import (
	"fmt"
	"github.com/rendyananta/example-online-book-store/internal/entity/user"
)

type userMatcher struct {
	input interface{}
}

func (u userMatcher) Matches(x interface{}) bool {
	input, okInput := u.input.(user.User)
	actual, okActual := x.(user.User)

	if !okInput && !okActual {
		return false
	}

	if input.ID != actual.ID {
		return false
	}

	if input.Name != actual.Name {
		return false
	}

	if input.Email != actual.Email {
		return false
	}

	return true
}

func (u userMatcher) String() string {
	return fmt.Sprintf("is matched to %v", u.input)
}
