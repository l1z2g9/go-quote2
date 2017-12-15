package rthk

import (
	"log"
	"testing"
)

func TestGetRthkNews(t *testing.T) {
	log.Println(GetRthkNews("NULL"))
}

/*func TestGetRthkNewsDetail(t *testing.T) {
	log.Println(GetRthkNewsDetail("1291808-20161020", "ch"))
}
*/
