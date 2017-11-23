package rthk

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
)

func TestGetCityShapshots(t *testing.T) {
	fmt.Println(GetCitySnapshots("England"))
}

func TestGetYCantonese(t *testing.T) {
	fmt.Println(GetYCantonese())
}

func TestDrivePath(t *testing.T) {
	c, _ := ioutil.ReadFile("googledrive_sourcecode.html")
	re := regexp.MustCompile(`window\['_DRIVE_ivd'\] = '\[\[\[(.+)1\]\\n';`)
	items := re.FindStringSubmatch(string(c))[1]
	fmt.Println("items", len(items))
}
