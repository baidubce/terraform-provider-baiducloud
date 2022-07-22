package baiducloud

import (
	"fmt"
	"regexp"

	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func validateReservationLength() schema.SchemaValidateFunc {
	return validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36})
}

func validateReservationUnit() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{"month"}, false)
}

func validatePaymentTiming() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{PaymentTimingPostpaid, PaymentTimingPrepaid}, false)
}

func validatePort() schema.SchemaValidateFunc {
	return validation.IntBetween(1, 65535)
}

func validateHttpMethod() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"GET",
		"PUT",
		"DELETE",
		"POST",
		"HEAD",
	}, false)
}

func validateStorageType() schema.SchemaValidateFunc {
	return validateStringFormat()
}

func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}

	return
}

func validateInstanceType() schema.SchemaValidateFunc {
	return validateStringFormat()
}

func validateStringFormat() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9]+`), "value must be alphanumeric and underscore")
}

func validateSubnetType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(vpc.SUBNET_TYPE_BCC),
		string(vpc.SUBNET_TYPE_BCCNAT),
		string(vpc.SUBNET_TYPE_BBC),
	}, false)
}

func validateBOSBucketACL() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		api.CANNED_ACL_PRIVATE,
		api.CANNED_ACL_PUBLIC_READ,
		api.CANNED_ACL_PUBLIC_READ_WRITE,
	}, false)
}

func validateBOSBucketStorageClass() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		api.STORAGE_CLASS_COLD,
		api.STORAGE_CLASS_STANDARD,
		api.STORAGE_CLASS_STANDARD_IA,
		STORAGE_CLASS_ARCHIVE,
	}, false)
}

func validateBOSBucketRCStorageClass() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		api.STORAGE_CLASS_COLD,
		api.STORAGE_CLASS_STANDARD,
		api.STORAGE_CLASS_STANDARD_IA,
	}, false)
}

func validateBOSBucketLifecycleRuleActionName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		ACTION_TRANSITION,
		ACTION_DELETEOBJECT,
		ACTION_ABORTMULTIPARTUPLOAD,
	}, false)
}

func validateBOSBucketLifecycleRuleActionStorage() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		api.STORAGE_CLASS_STANDARD_IA,
		api.STORAGE_CLASS_COLD,
		STORAGE_CLASS_ARCHIVE,
	}, false)
}

func validateBOSObjectACL() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		api.CANNED_ACL_PRIVATE,
		api.CANNED_ACL_PUBLIC_READ,
	}, false)
}

func validateBOSObjectCacheControl() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		BOS_BUCKET_OBJECT_CACHE_CONTROL_PRIVATE,
		BOS_BUCKET_OBJECT_CACHE_CONTROL_NO_CACHE,
		BOS_BUCKET_OBJECT_CACHE_CONTROL_MAX_AGE,
		BOS_BUCKET_OBJECT_CACHE_CONTROL_MUST_REVALIDATE,
	}, false)
}

func validateBOSObjectContentDisposition() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		BOS_BUCKET_OBJECT_CONTENT_DISPOSITION_INLINE,
		BOS_BUCKET_OBJECT_CONTENT_DISPOSITION_ATTACHMENT,
	}, false)
}

func validateCFCMemorySize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value%128 != 0 {
		errors = append(errors, fmt.Errorf("%q is not an integer multiple of 128", k))
	}

	if value > 1024 || value < 128 {
		errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, 128, 1024, value))
	}

	return
}
