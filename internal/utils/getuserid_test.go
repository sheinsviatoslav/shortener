package utils

//import (
//	"github.com/stretchr/testify/assert"
//	"golang.org/x/net/context"
//	"net/http"
//	"testing"
//)
//
//func TestGetUserID(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		ctx := context.WithValue(context.Background(), "userID", "user")
//		req, err := http.NewRequestWithContext(ctx, "GET", "test", nil)
//		assert.NoError(t, err)
//
//		userID := GetUserID(req)
//		assert.Equal(t, "user", userID)
//	})
//
//	t.Run("no value", func(t *testing.T) {
//		req, err := http.NewRequest("GET", "test", nil)
//		assert.NoError(t, err)
//		userID := GetUserID(req)
//		assert.Equal(t, "", userID)
//	})
//}
