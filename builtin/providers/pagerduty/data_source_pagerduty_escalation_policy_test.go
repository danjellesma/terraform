package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourcePagerDutyEscalationPolicy_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePagerDutyEscalationPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEscalationPolicy("pagerduty_escalation_policy.test", "data.pagerduty_escalation_policy.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyEscalationPolicy(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a escalation policy ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the escalation policy %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEscalationPolicyConfig(rName string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "TF User %[1]s"
  email = "tf.%[1]s@example.com"
}

resource "pagerduty_escalation_policy" "test" {
  name        = "TF Escalation Policy %[1]v"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = "${data.pagerduty_user.test.id}"
    }
  }
}

data "pagerduty_escalation_policy" "by_name" {
  name = "${pagerduty_escalation_policy.test.name}"
}
`, rName)
}
