package soundcloud

import (
	//"fmt"
	"../util"
	"testing"
)

func TestGetPlaylists(t *testing.T) {
	util.ShowLog()
	s := Soundcloud{}
	//s.Login()
	//util.Info.Println(">> ", s.Login())

	//s = Soundcloud{"1-274177-278918219-efb057192def2"}
	//fmt.Println(string(s.GetPlaylists()))

	util.Info.Println(string(s.GetPlaylist("305138956")))
}
