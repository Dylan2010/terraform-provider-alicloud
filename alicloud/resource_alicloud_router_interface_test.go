package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRouterInterface_basic(t *testing.T) {
	var vpcInstance vpc.DescribeVpcAttributeResponse
	var ri vpc.RouterInterfaceTypeInDescribeRouterInterfaces
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_router_interface.interface",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouterInterfaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"alicloud_vpc.foo", &vpcInstance),
					testAccCheckRouterInterfaceExists(
						"alicloud_router_interface.interface", &ri),
					resource.TestCheckResourceAttr(
						"alicloud_router_interface.interface", "name", "testAccRouterInterfaceConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_router_interface.interface", "role", "InitiatingSide"),
				),
			},
		},
	})

}

func testAccCheckRouterInterfaceExists(n string, ri *vpc.RouterInterfaceTypeInDescribeRouterInterfaces) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeRouterInterface(client.RegionId, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		ri = &response
		return nil
	}
}

func testAccCheckRouterInterfaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*AliyunClient)

		ri, err := client.DescribeRouterInterface(client.RegionId, rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.RouterInterfaceId == rs.Primary.ID {
			return fmt.Errorf("Interface %s still exists.", rs.Primary.ID)
		}
	}
	return nil
}

const testAccRouterInterfaceConfig = `
resource "alicloud_vpc" "foo" {
  name = "tf_test_foo12345"
  cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "current_regions" {
  current = true
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "testAccRouterInterfaceConfig"
  description = "testAccRouterInterfaceConfig"
}`
