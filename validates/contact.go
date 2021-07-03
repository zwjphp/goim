package validates

import (
	"github.com/gookit/validate"
	"goim/model"
)

type ContactValidate struct {
}

func (validatec *ContactValidate) ContactValidates(userid, dstid int64) (string, error) {
	contact := &model.Contact{
		Ownerid: userid,
		Dstobj:  dstid,
	}
	v := validate.Struct(contact)
	if v.Validate() {
		return "", nil
	} else {
		return v.Errors.One(), v.Errors
	}
}