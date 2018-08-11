// @APIVersion 1.0.0
// @Title SPG Test API
// @Description This API documents is created for friend management
// @Contact sbpeng2010@yahoo.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"spgapi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/friend",
			beego.NSInclude(
				&controllers.FriendController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
