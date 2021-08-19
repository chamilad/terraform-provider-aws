package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53recoverycontrolconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSRoute53RecoveryControlConfigControlPanel_basic(t *testing.T) {
	rClusterName := acctest.RandomWithPrefix("tf-acc-test-cluster")
	rControlPanelName := acctest.RandomWithPrefix("tf-acc-test-control-panel")
	resourceName := "aws_route53recoverycontrolconfig_control_panel.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, route53recoverycontrolconfig.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsRoute53RecoveryControlConfigControlPanelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryControlConfigControlPanelConfig(rClusterName, rControlPanelName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryControlConfigControlPanelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rControlPanelName),
					resource.TestCheckResourceAttr(resourceName, "status", "DEPLOYED"),
					resource.TestCheckResourceAttr(resourceName, "default_control_panel", "false"),
					resource.TestCheckResourceAttr(resourceName, "routing_control_count", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAwsRoute53RecoveryControlConfigControlPanelDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).route53recoverycontrolconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53recoverycontrolconfig_control_panel" {
			continue
		}

		input := &route53recoverycontrolconfig.DescribeControlPanelInput{
			ControlPanelArn: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeControlPanel(input)

		if err == nil {
			return fmt.Errorf("Route53RecoveryControlConfig Control Panel (%s) not deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccAwsRoute53RecoveryControlConfigClusterSetUp(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53recoverycontrolconfig_cluster" "test" {
  name = %[1]q
}
`, rName)
}

func testAccAwsRoute53RecoveryControlConfigControlPanelConfig(rName, rName2 string) string {
	return composeConfig(testAccAwsRoute53RecoveryControlConfigClusterSetUp(rName), fmt.Sprintf(`
resource "aws_route53recoverycontrolconfig_control_panel" "test" {
  name        = %q
  cluster_arn = aws_route53recoverycontrolconfig_cluster.test.cluster_arn
}
`, rName2))
}

func testAccCheckAwsRoute53RecoveryControlConfigControlPanelExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).route53recoverycontrolconfigconn

		input := &route53recoverycontrolconfig.DescribeControlPanelInput{
			ControlPanelArn: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeControlPanel(input)

		return err
	}
}
