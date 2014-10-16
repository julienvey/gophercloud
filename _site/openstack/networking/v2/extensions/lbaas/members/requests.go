package members

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status       string `q:"status"`
	Weight       int    `q:"weight"`
	AdminStateUp *bool  `q:"admin_state_up"`
	TenantID     string `q:"tenant_id"`
	PoolID       string `q:"pool_id"`
	Address      string `q:"address"`
	ProtocolPort int    `q:"protocol_port"`
	ID           string `q:"id"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
}

// List returns a Pager which allows you to iterate over a collection of
// pools. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those pools that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.LastHTTPResponse) pagination.Page {
		return MemberPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	})
}

// CreateOpts contains all the values needed to create a new pool member.
type CreateOpts struct {
	// Only required if the caller has an admin role and wants to create a pool
	// for another tenant.
	TenantID string

	// Required. The IP address of the member.
	Address string

	// Required. The port on which the application is hosted.
	ProtocolPort int

	// Required. The pool to which this member will belong.
	PoolID string
}

// Create accepts a CreateOpts struct and uses the values to create a new
// load balancer pool member.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	type member struct {
		TenantID     string `json:"tenant_id"`
		ProtocolPort int    `json:"protocol_port"`
		Address      string `json:"address"`
		PoolID       string `json:"pool_id"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{
		Address:      opts.Address,
		TenantID:     opts.TenantID,
		ProtocolPort: opts.ProtocolPort,
		PoolID:       opts.PoolID,
	}}

	var res CreateResult
	_, res.Err = perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})
	return res
}

// Get retrieves a particular pool member based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOpts contains the values used when updating a pool member.
type UpdateOpts struct {
	// The administrative state of the member, which is up (true) or down (false).
	AdminStateUp bool
}

// Update allows members to be updated.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	type member struct {
		AdminStateUp bool `json:"admin_state_up"`
	}
	type request struct {
		Member member `json:"member"`
	}

	reqBody := request{Member: member{AdminStateUp: opts.AdminStateUp}}

	// Send request to API
	var res UpdateResult
	_, res.Err = perigee.Request("PUT", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	return res
}

// Delete will permanently delete a particular member based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", resourceURL(c, id), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
