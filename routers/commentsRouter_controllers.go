package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "BlockUpdate",
			Router: `/block`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "CreateConnection",
			Router: `/create`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "ListFriends",
			Router: `/list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "ListCommonFriends",
			Router: `/listCommon`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "ListUpdate",
			Router: `/listUpdate`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["spgapi/controllers:FriendController"] = append(beego.GlobalControllerRouter["spgapi/controllers:FriendController"],
		beego.ControllerComments{
			Method: "SubscribeUpdate",
			Router: `/subscribe`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
