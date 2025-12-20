package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/InfoSecured/globalscape-eft-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"globalscapeeft": providerserver.NewProtocol6WithError(New()),
}

func TestAccSitesDataSource_basic(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() + `
data "globalscapeeft_sites" "all" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.globalscapeeft_sites.all", "sites.#"),
				),
			},
		},
	})
}

func TestAccSiteUser_basic(t *testing.T) {
	testAccPreCheck(t)
	siteID := os.Getenv("EFT_TEST_SITE_ID")
	if siteID == "" {
		t.Skip("EFT_TEST_SITE_ID must be set for site user acceptance tests")
	}

	loginName := fmt.Sprintf("tf-acctest-%d", os.Getpid())
	resourceName := "globalscapeeft_site_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSiteUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteUserConfig(siteID, loginName, "Terraform Example", "tf@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "login_name", loginName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Example"),
				),
			},
			{
				Config: testAccSiteUserConfig(siteID, loginName, "Terraform Updated", "tf-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Updated"),
					resource.TestCheckResourceAttr(resourceName, "email", "tf-updated@example.com"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("EFT_TEST_HOST") == "" ||
		os.Getenv("EFT_TEST_USERNAME") == "" ||
		os.Getenv("EFT_TEST_PASSWORD") == "" {
		t.Skip("EFT_TEST_HOST, EFT_TEST_USERNAME, and EFT_TEST_PASSWORD must be set for acceptance tests")
	}
}

func testAccProviderConfig() string {
	authType := os.Getenv("EFT_TEST_AUTHTYPE")
	if authType == "" {
		authType = "EFT"
	}
	insecure := os.Getenv("EFT_TEST_INSECURE")
	if insecure == "" {
		insecure = "false"
	}

	return fmt.Sprintf(`
provider "globalscapeeft" {
  host                 = %q
  username             = %q
  password             = %q
  auth_type            = %q
  insecure_skip_verify = %s
}
`, os.Getenv("EFT_TEST_HOST"), os.Getenv("EFT_TEST_USERNAME"), os.Getenv("EFT_TEST_PASSWORD"), authType, insecure)
}

func testAccSiteUserConfig(siteID, loginName, displayName, email string) string {
	return fmt.Sprintf(`
%s

resource "globalscapeeft_site_user" "test" {
  site_id         = %q
  login_name      = %q
  password        = "TerraformP@ssw0rd!"
  password_type   = "Default"
  display_name    = %q
  email           = %q
  account_enabled = "yes"
}
`, testAccProviderConfig(), siteID, loginName, displayName, email)
}

func testAccClient() (*client.Client, error) {
	authType := os.Getenv("EFT_TEST_AUTHTYPE")
	if authType == "" {
		authType = "EFT"
	}

	insecure := false
	if os.Getenv("EFT_TEST_INSECURE") == "true" {
		insecure = true
	}

	return client.NewClient(context.Background(), client.Config{
		BaseURL:            os.Getenv("EFT_TEST_HOST"),
		Username:           os.Getenv("EFT_TEST_USERNAME"),
		Password:           os.Getenv("EFT_TEST_PASSWORD"),
		AuthType:           authType,
		InsecureSkipVerify: insecure,
	})
}

func testAccCheckSiteUserDestroy(s *terraform.State) error {
	siteID := os.Getenv("EFT_TEST_SITE_ID")
	if siteID == "" {
		return nil
	}

	c, err := testAccClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "globalscapeeft_site_user" {
			continue
		}

		_, err := c.GetSiteUser(context.Background(), siteID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("user %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
