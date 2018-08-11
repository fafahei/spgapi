package models

import (
	"testing"
	"github.com/astaxie/beego/utils"
)

func TestConnect(t *testing.T) {
	var f = [2]*Info{&Info{0, "x@example.com"}, &Info{0, "xxx@example.com"}}
	status, err := Connect(f)
	if status == false && err.Error() == "This friend x@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "andy@example.com"}, &Info{0, "xxx@example.com"}}
	status, err = Connect(f)
	if status == false && err.Error() == "This friend xxx@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "andy@example.com"}, &Info{0, "john@example.com"}}
	status, err = Connect(f)
	if status == false && err.Error() == "This connection relationship has been created!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "mnt@example.com"}, &Info{0, "ytu@example.com"}}
	status, err = Connect(f)
	if status == false && err.Error() == "mnt@example.com blocks ytu@example.com, so can't create friends connection!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "ytu@example.com"}, &Info{0, "mnt@example.com"}}
	status, err = Connect(f)
	if status == false && err.Error() == "ytu@example.com blocks mnt@example.com, so can't create friends connection!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "111@example.com"}, &Info{0, "222@example.com"}}
	status, err = Connect(f)
	if status == true && err == nil {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}

func TestGetFriends(t *testing.T) {
	var friend = &Info{0, "xxx@example.com"}
	status, _, err := GetFriends(friend)
	if status == false && err.Error() == "This friend xxx@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	friend = &Info{0, "kate@example.com"}
	status, _, err = GetFriends(friend)
	if status == true && err.Error() == "Can't find any friend for kate@example.com!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	friend = &Info{0, "andy@example.com"}
	status, f, err := GetFriends(friend)
	if utils.InSlice("john@example.com", f) && utils.InSlice("common@example.com", f) {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}

func TestGetCommonFriends(t *testing.T) {
	var f = [2]*Info{&Info{0, "x@example.com"}, &Info{0, "xxx@example.com"}}
	status, _, err := GetCommonFriends(f)
	if status == false && err.Error() == "This friend x@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "andy@example.com"}, &Info{0, "xxx@example.com"}}
	status, _, err = GetCommonFriends(f)
	if status == false && err.Error() == "This friend xxx@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "andy@example.com"}, &Info{0, "kate@example.com"}}
	status, _, err = GetCommonFriends(f)
	if status == true && err.Error() == "Can't find any friend for kate@example.com!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{&Info{0, "andy@example.com"}, &Info{0, "abc@example.com"}}
	status, _, err = GetCommonFriends(f)
	if status == true && err.Error() == "No common friends!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	f = [2]*Info{{0, "andy@example.com"}, {0, "john@example.com"}}
	status, c, err := GetCommonFriends(f)
	if utils.InSlice("common@example.com", c) {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}

func TestSubscribe(t *testing.T) {
	var r = &Relationship{0, &Info{0, "x@example.com"}, &Info{0, "xxx@example.com"}, 0}
	status, err := Subscribe(r)
	if status == false && err.Error() == "This friend x@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "andy@example.com"}
	r.Target = &Info{0, "xxx@example.com"}
	status, err = Subscribe(r)
	if status == false && err.Error() == "This friend xxx@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "mnt@example.com"}
	r.Target = &Info{0, "ytu@example.com"}
	status, err = Subscribe(r)
	if status == true && err == nil {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "lisa@example.com"}
	r.Target = &Info{0, "john@example.com"}
	status, err = Subscribe(r)
	if status == false && err.Error() == "This subscription relationship has been created!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "111@example.com"}
	r.Target = &Info{0, "222@example.com"}
	status, err = Subscribe(r)
	if status == true && err == nil {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}

func TestBlock(t *testing.T) {
	var r = &Relationship{0, &Info{0, "x@example.com"}, &Info{0, "xxx@example.com"}, 0}
	status, err := Block(r)
	if status == false && err.Error() == "This friend x@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "andy@example.com"}
	r.Target = &Info{0, "xxx@example.com"}
	status, err = Block(r)
	if status == false && err.Error() == "This friend xxx@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "111@example.com"}
	r.Target = &Info{0, "222@example.com"}
	status, err = Block(r)
	if status == true && err == nil {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "andy@example.com"}
	r.Target = &Info{0, "lisa@example.com"}
	status, err = Block(r)
	if status == false && err.Error() == "This block relationship has been created!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	r.Requestor = &Info{0, "333@example.com"}
	r.Target = &Info{0, "444@example.com"}
	status, err = Block(r)
	if status == true && err == nil {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}

func TestGetUpdateFriends(t *testing.T) {
	var target = &Info{0, "x@example.com"}
	var r = []*Info{&Info{0, "kate@example.com"}}
	status, _, err := GetUpdateFriends(target, r)
	if status == false && err.Error() == "This friend x@example.com is non-exist!" {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	target = &Info{0, "john@example.com"}
	r = []*Info{&Info{0, "kate@example.com"}, {0, "xxx@example.com"}}
	status, f, err := GetUpdateFriends(target, r)
	if utils.InSlice("kate@example.com", f) && utils.InSlice("lisa@example.com", f) {
		t.Log("pass")
	} else {
		t.Log("failed")
	}

	target = &Info{0, "john@example.com"}
	r = []*Info{&Info{0, "kate@example.com"}, {0, "mnt@example.com"}}
	status, f, err = GetUpdateFriends(target, r)
	if utils.InSlice("kate@example.com", f) && utils.InSlice("lisa@example.com", f) && utils.InSlice("mnt@example.com", f) {
		t.Log("pass")
	} else {
		t.Log("failed")
	}
}