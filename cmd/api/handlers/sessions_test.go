package handlers_test

import (
	"testing"

	test "github.com/eyecuelab/kit/web/testing/integration"
)

func TestGetSession(t *testing.T) {
	resp, jsonAPI, _ := test.Get(t, "", "")
	test.AssertStatusOK(t, resp)
	test.AssertAttr(t, jsonAPI.Data.Attributes, map[string]interface{}{
		"something_extra": "abc",
	})
	test.AssertAction(t, jsonAPI.Data.Meta, "login", "some_other_common_action")
	test.AssertLink(t, jsonAPI.Data.Links, "self")
}

func TestAuthGetSession(t *testing.T) {
	token := authToken(t, "user1@example.com")
	resp, jsonAPI, _ := test.Get(t, "", token)
	test.AssertStatusOK(t, resp)
	test.AssertAttr(t, jsonAPI.Data.Attributes, map[string]interface{}{
		"something_extra": "abc",
	})
	test.AssertAction(t, jsonAPI.Data.Meta, "logout", "update_profile", "some_other_common_action")
	test.AssertLink(t, jsonAPI.Data.Links, "self", "users", "profile")
}
