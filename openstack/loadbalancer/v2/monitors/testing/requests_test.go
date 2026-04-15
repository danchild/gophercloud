package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListHealthmonitors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorListSuccessfully(t, fakeServer)

	pages := 0
	err := monitors.List(fake.ServiceClient(fakeServer), monitors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := monitors.ExtractMonitors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 healthmonitors, got %d", len(actual))
		}
		th.CheckDeepEquals(t, HealthmonitorWeb, actual[0])
		th.CheckDeepEquals(t, HealthmonitorDb, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllHealthmonitors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorListSuccessfully(t, fakeServer)

	allPages, err := monitors.List(fake.ServiceClient(fakeServer), monitors.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := monitors.ExtractMonitors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, HealthmonitorWeb, actual[0])
	th.CheckDeepEquals(t, HealthmonitorDb, actual[1])
}

func TestCreateHealthmonitor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorCreationSuccessfully(t, fakeServer, SingleHealthmonitorBody)

	actual, err := monitors.Create(context.TODO(), fake.ServiceClient(fakeServer), monitors.CreateOpts{
		Type:           "HTTP",
		DomainName:     "www.example.com",
		Name:           "db",
		PoolID:         "84f1b61f-58c4-45bf-a8a9-2dafb9e5214d",
		ProjectID:      "453105b9-1754-413f-aab1-55f1af620750",
		Delay:          20,
		Timeout:        10,
		MaxRetries:     5,
		MaxRetriesDown: 4,
		HTTPVersion:    "1.1",
		Tags:           []string{},
		URLPath:        "/check",
		ExpectedCodes:  "200-299",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, HealthmonitorDb, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := monitors.Create(context.TODO(), fake.ServiceClient(fakeServer), monitors.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = monitors.Create(context.TODO(), fake.ServiceClient(fakeServer), monitors.CreateOpts{Type: monitors.TypeHTTP})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetHealthmonitor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := monitors.Get(context.TODO(), client, "5d4b5228-33b0-4e60-b225-9b727c1a20e7").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, HealthmonitorDb, *actual)
}

func TestDeleteHealthmonitor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorDeletionSuccessfully(t, fakeServer)

	res := monitors.Delete(context.TODO(), fake.ServiceClient(fakeServer), "5d4b5228-33b0-4e60-b225-9b727c1a20e7")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateHealthmonitor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleHealthmonitorUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := monitors.Update(context.TODO(), client, "5d4b5228-33b0-4e60-b225-9b727c1a20e7", monitors.UpdateOpts{
		Name:           ptr.To("NewHealthmonitorName"),
		Delay:          ptr.To(3),
		Timeout:        ptr.To(20),
		MaxRetries:     ptr.To(10),
		MaxRetriesDown: ptr.To(8),
		URLPath:        ptr.To("/another_check"),
		ExpectedCodes:  ptr.To("301"),
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, HealthmonitorUpdated, *actual)
}

func TestDelayMustBeGreaterOrEqualThanTimeout(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	_, err := monitors.Create(context.TODO(), fake.ServiceClient(fakeServer), monitors.CreateOpts{
		Type:          "HTTP",
		PoolID:        "d459f7d8-c6ee-439d-8713-d3fc08aeed8d",
		Delay:         1,
		Timeout:       10,
		MaxRetries:    5,
		URLPath:       "/check",
		ExpectedCodes: "200-299",
	}).Extract()

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	_, err = monitors.Update(context.TODO(), fake.ServiceClient(fakeServer), "453105b9-1754-413f-aab1-55f1af620750", monitors.UpdateOpts{
		Delay:   ptr.To(1),
		Timeout: ptr.To(10),
	}).Extract()

	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}
