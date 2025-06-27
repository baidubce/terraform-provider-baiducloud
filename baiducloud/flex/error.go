package flex

import (
	"errors"
	
	"github.com/baidubce/bce-sdk-go/bce"
)

func IsResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	var bceSdkErr *bce.BceServiceError
	if errors.As(err, &bceSdkErr) {
		return bceSdkErr.StatusCode == 404
	}
	return false
}
