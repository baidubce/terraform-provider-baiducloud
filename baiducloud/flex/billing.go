package flex

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	PaymentTimingPostpaid = "Postpaid"
	PaymentTimingPrepaid  = "Prepaid"
)

func SchemaPaymentTiming() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Payment timing of billing. Valid values: `Prepaid`, `Postpaid`. Defaults to `Postpaid`.",
		Optional:     true,
		Default:      PaymentTimingPostpaid,
		ValidateFunc: ValidatePaymentTiming(),
	}
}

func ComputedSchemaPaymentTiming() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Payment timing of billing. Possible values: `Prepaid`, `Postpaid`.",
		Computed:    true,
	}
}

func SchemaReservationLength() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeInt,
		Description: "The reservation length (month) will pay. Effective when `payment_timing` is `Prepaid`. " +
			"Valid values: `1`~`9`, `12`, `24`, `36`. Defaults to `1`.",
		Optional:         true,
		Default:          1,
		ValidateFunc:     ValidateReservationLength(),
		DiffSuppressFunc: PostPaidDiffSuppressFunc,
	}
}

func SchemaAutoRenewLength() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeInt,
		Description: "The automatic renewal time (month). Effective when `payment_timing` is `Prepaid`. " +
			"Valid values: `1`~`9`, `12`, `24`, `36`. Defaults to `1`.",
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: ValidateReservationLength(),
	}
}

func SchemaAutoRenewTimeUnit() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Auto renew time unit, currently only supports monthly.",
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"month"}, false),
	}
}
func ValidatePaymentTiming() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{PaymentTimingPostpaid, PaymentTimingPrepaid}, false)
}

func ValidateReservationLength() schema.SchemaValidateFunc {
	return validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36})
}

func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("payment_timing").(string) == PaymentTimingPostpaid
}
