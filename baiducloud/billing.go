package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform/helper/schema"
)

func setBilling(d *schema.ResourceData, paymentTiming string) {
	billingMap := map[string]interface{}{"payment_timing": paymentTiming}
	billings := []interface{}{}
	billings = append(billings, billingMap)
	d.Set("billing", billings)
}

func createBillingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
				Optional:     true,
				Default:      bbc.PaymentTimingPostPaid,
				ValidateFunc: validatePaymentTiming(),
			},
			"reservation": {
				Type:        schema.TypeMap,
				Description: "Reservation of the instance.",
				Optional:    true,
				//DiffSuppressFunc: postPaidDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_length": {
							Type:         schema.TypeInt,
							Description:  "The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:     true,
							Default:      1,
							ValidateFunc: validateReservationLength(),
							//DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type:         schema.TypeString,
							Description:  "The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
							Required:     true,
							Default:      "Month",
							ValidateFunc: validateReservationUnit(),
							//DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
		},
	}
}
