package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/ovh/go-ovh/ovh"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccOVHClient *ovh.Client

func init() {
	log.SetOutput(os.Stdout)
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ovh": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	v := os.Getenv("OVH_ENDPOINT")
	if v == "" {
		t.Fatal("OVH_ENDPOINT must be set for acceptance tests")
	}

	v = os.Getenv("OVH_APPLICATION_KEY")
	if v == "" {
		t.Fatal("OVH_APPLICATION_KEY must be set for acceptance tests")
	}

	v = os.Getenv("OVH_APPLICATION_SECRET")
	if v == "" {
		t.Fatal("OVH_APPLICATION_SECRET must be set for acceptance tests")
	}

	v = os.Getenv("OVH_CONSUMER_KEY")
	if v == "" {
		t.Fatal("OVH_CONSUMER_KEY must be set for acceptance tests")
	}

	v = os.Getenv("OVH_VRACK")
	if v == "" {
		t.Fatal("OVH_VRACK must be set for acceptance tests")
	}

	v = os.Getenv("OVH_PUBLIC_CLOUD")
	if v == "" {
		t.Fatal("OVH_PUBLIC_CLOUD must be set for acceptance tests")
	}

	v = os.Getenv("OVH_ZONE")
	if v == "" {
		t.Fatal("OVH_ZONE must be set for acceptance tests")
	}

	v = os.Getenv("OVH_IPLB_SERVICE")
	if v == "" {
		t.Fatal("OVH_IPLB_SERVICE must be set for acceptance tests")
	}

	v = os.Getenv("OVH_IP_BLOCK")
	if v == "" {
		t.Fatal("OVH_IP_BLOCK must be set for acceptance tests")
	}

	v = os.Getenv("OVH_IP")
	if v == "" {
		t.Fatal("OVH_IP must be set for acceptance tests")
	}

	v = os.Getenv("OVH_IP_REVERSE")
	if v == "" {
		t.Fatal("OVH_IP_REVERSE must be set for acceptance tests")
	}

	if testAccOVHClient == nil {
		config := Config{
			Endpoint:          os.Getenv("OVH_ENDPOINT"),
			ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
			ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
			ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
		}

		if err := config.loadAndValidate(); err != nil {
			t.Fatalf("couln't load OVH Client: %s", err)
		} else {
			testAccOVHClient = config.OVHClient
		}
	}
}

func testAccCheckVRackExists(t *testing.T) {
	type vrackResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	r := vrackResponse{}

	endpoint := fmt.Sprintf("/vrack/%s", os.Getenv("OVH_VRACK"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read VRack %s -> name:'%s', desc:'%s' ", endpoint, r.Name, r.Description)

}

func testAccCheckPublicCloudExists(t *testing.T) {
	type cloudProjectResponse struct {
		ID          string `json:"project_id"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}

	r := cloudProjectResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s", os.Getenv("OVH_PUBLIC_CLOUD"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read Cloud Project %s -> status: '%s', desc: '%s'", endpoint, r.Status, r.Description)

}

func testAccCheckIpLoadbalancingExists(t *testing.T) {
	type iplbResponse struct {
		ServiceName string `json:"serviceName"`
		State       string `json:"state"`
	}

	r := iplbResponse{}

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s", os.Getenv("OVH_IPLB_SERVICE"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read IPLB service %s -> state: '%s', serviceName: '%s'", endpoint, r.State, r.ServiceName)

}

func testAccCheckDomainZoneExists(t *testing.T) {
	type domainZoneResponse struct {
		NameServers []string `json:"nameServers"`
	}

	r := domainZoneResponse{}

	endpoint := fmt.Sprintf("/domain/zone/%s", os.Getenv("OVH_ZONE"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}

	t.Logf("Read Domain Zone %s -> nameservers: '%v'", endpoint, r.NameServers)

}
