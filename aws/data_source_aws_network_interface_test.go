package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAwsNetworkInterface_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsNetworkInterface_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_network_interface.test", "private_ips.#", "1"),
					resource.TestCheckResourceAttr("data.aws_network_interface.test", "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAwsNetworkInterface_basic(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test" {
  cidr_block = "10.0.0.0/24"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_security_group" "test" {
  name = "tf-sg-%s"
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_network_interface" "test" {
  subnet_id = "${aws_subnet.test.id}"
  private_ips = ["10.0.0.50"]
  security_groups = ["${aws_security_group.test.id}"]
}

data "aws_network_interface" "test" {
  id = "${aws_network_interface.test.id}"
}
    `, rName)
}

func TestAccDataSourceAwsNetworkInterface_filters(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsNetworkInterface_filters(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_network_interface.test", "private_ips.#", "1"),
					resource.TestCheckResourceAttr("data.aws_network_interface.test", "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAwsNetworkInterface_filters(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test" {
  cidr_block = "10.0.0.0/24"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_security_group" "test" {
  name = "tf-sg-%s"
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_network_interface" "test" {
  subnet_id = "${aws_subnet.test.id}"
  private_ips = ["10.0.0.60"]
  security_groups = ["${aws_security_group.test.id}"]
}

data "aws_network_interface" "test" {
  filter {
	name = "subnet-id"
	values = ["${aws_subnet.test.id}"]
  }

  filter {
	name = "group-id"
	values = ["${aws_security_group.test.id}"]
  }

  depends_on = ["aws_network_interface.test"]
}
    `, rName)
}

func TestAccDataSourceAwsNetworkInterface_error_required(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsNetworkInterface_error_required(),
				ExpectError: regexp.MustCompile("One of id or filter must be assigned"),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsNetworkInterface_error_required_non_empty_id(),
				ExpectError: regexp.MustCompile("One of id or filter must be assigned"),
			},
		},
	})
}

func testAccDataSourceAwsNetworkInterface_error_required() string {
	return `data "aws_network_interface" "test" {}`
}

func testAccDataSourceAwsNetworkInterface_error_required_non_empty_id() string {
	return `
data "aws_network_interface" "test" {
	id = ""
}`
}
