package news

import (
	"fmt"
	"testing"
)

func TestFeedlyGetCategories(t *testing.T) {
	f := Feedly{}
	fmt.Println(string(f.GetCategories()))
}

func TestFeedlySubscriptions(t *testing.T) {
	f := Feedly{}
	fmt.Println(string(f.GetSubscriptions()))
}

func TestGetListByFeed(t *testing.T) {
	f := Feedly{}
	f.GetListByFeed("feed/http://dave.cheney.net/feed")
}

func TestFeedlyDetail(t *testing.T) {
	f := Feedly{}
	fmt.Println(string(f.GetEntryContent("74r0EQBzVZk4gOe3iYARQrNXuQxwM4qcgVbV4TVmzUg=_14e7534d049:123f5b71:cd74fcc6")))
}

func TestFeedlyProfile(t *testing.T) {
	f := Feedly{}
	fmt.Println(string(f.GetProfile()))
}
