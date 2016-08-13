package wherebouts

import (
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/wherebouts/calc"
	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

type datastoreDB struct {
	client *datastore.Client
}

var _ VendorDatabase = &datastoreDB{}

func newDatastoreDB(client *datastore.Client) (VendorDatabase, error) {
	ctx := context.Background()

	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}

	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}

	return &datastoreDB{
		client: client,
	}, nil
}

func (db *datastoreDB) Close() {
	// nop
}

func (db *datastoreDB) datastoreKey(id int64) *datastore.Key {
	ctx := context.Background()
	return datastore.NewKey(ctx, "Vendor", "", id, nil)
}

func (db *datastoreDB) lookupByEmail(email string) (*Vendor, *datastore.Key) {
	ctx := context.Background()
	var v []Vendor
	query := datastore.NewQuery("Vendor").Filter("Email =", email).Limit(1)
	key, err := db.client.GetAll(ctx, query, &v)

	if err != nil {
		fmt.Print("There's some error in the key lookup?")
		return nil, nil
	}

	return &v[0], key[0]
}

// GetVendorByID retrieves a vendor by their ID
func (db *datastoreDB) GetVendorByID(id int64) (*Vendor, error) {
	ctx := context.Background()
	k := db.datastoreKey(id)
	vendor := &Vendor{}

	if err := db.client.Get(ctx, k, vendor); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get Vendor: %v", err)
	}

	//vendor.ID = id
	return vendor, nil
}

func (db *datastoreDB) AddVendor(v *Vendor) (id int64, err error) {
	ctx := context.Background()

	k := datastore.NewIncompleteKey(ctx, "Vendor", nil)
	_, err = db.client.Put(ctx, k, v)
	if err != nil {
		return 0, fmt.Errorf("datastoredb: could not put Vendor: %v", err)
	}

	return id, nil
}

/*
func (db *datastoreDB) UpdateVendor(v *Vendor) error {
	ctx := context.Background()
	k := db.datastoreKey(v.ID)
	if _, err := db.client.Put(ctx, k, v); err != nil {
		return fmt.Errorf("datastoredb: could not update Vendor: %v", err)
	}

	return nil
}
*/

func (db *datastoreDB) ToggleAvailability(email string, lat float64, long float64) {
	v, k := db.lookupByEmail(email)
	ctx := context.Background()

	v.Open = !v.Open
	v.Latitude = lat
	v.Longitude = long

	if _, err := db.client.Put(ctx, k, v); err != nil {
		fmt.Print("error in the put")
		return
	}

}

func (db *datastoreDB) GetVendorByEmail(email string) (*Vendor, error) {
	ctx := context.Background()
	vendors := make([]*Vendor, 0)
	q := datastore.NewQuery("Vendor").
		Filter("Email =", email).
		Limit(1)

	key, err := db.client.GetAll(ctx, q, &vendors)

	if err != nil || len(key) == 0 {
		return nil, fmt.Errorf("datastoredb: could not find vendor: %v", err)
	}

	//	vendors[0].ID = key[0].ID()

	return vendors[0], nil
}

func filterBetweenLongitudes(max float64, min float64, db *datastoreDB) ([]int64, error) {
	ctx := context.Background()

	vendors := make([]*Vendor, 0)

	find := datastore.NewQuery("Vendor").
		//	Filter("Open =", true).
		Filter("Longitude <=", max).
		Filter("Longitude >=", min)

	keys, err := db.client.GetAll(ctx, find, &vendors)

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, fmt.Errorf("datastoredb: could not find vendors: %v", err)
	}

	ids := make([]int64, len(vendors))

	for i, j := range keys {
		ids[i] = j.ID()
	}

	return ids, nil
}

func filterBetweenLatitudes(max float64, min float64, db *datastoreDB) ([]int64, error) {
	ctx := context.Background()

	vendors := make([]*Vendor, 0)

	find := datastore.NewQuery("Vendor").
		Filter("Latitude <=", max).
		Filter("Latitude >=", min).
		Filter("Open =", true)

	keys, err := db.client.GetAll(ctx, find, &vendors)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not find vendors: %v", err)
	}
	fmt.Print("hi from lats\n")

	ids := make([]int64, len(vendors))

	fmt.Print("Can you find anything?", keys[0], "\n")
	for i, j := range keys {
		ids[i] = j.ID()
	}

	return ids, nil
}

func (db *datastoreDB) GetNearbyVendors(lat float64, long float64, dist float64) ([]*Vendor, error) {
	ctx := context.Background()
	latDist := calc.GetLatDelta(dist)
	longDist := calc.GetLongDelta(long, dist)

	xMax := lat + latDist
	xMin := lat - latDist
	yMax := long + longDist
	yMin := long - longDist

	fmt.Printf("xMax: %v, xMin: %v\nyMax: %v, yMinx: %v\n", xMax, xMin, yMax, yMin)

	// filter by latitude parameter
	lats, err := filterBetweenLatitudes(xMax, xMin, db)
	if err != nil {
		fmt.Print("hello", err)
		return nil, fmt.Errorf("databasedb: could not find vendors: %v\n", err)
	}

	fmt.Print("hi", len(lats), "\n")
	longs, err := filterBetweenLongitudes(yMax, yMin, db)
	if err != nil {
		fmt.Print("really")
		return nil, fmt.Errorf("databasedb: could not find vendors: %v\n", err)
	}

	fmt.Print("bye", len(longs), "\n")

	// return the intersect of the keys returned
	intersect := calc.Intersect(lats, longs)

	fmt.Printf("%v\n", intersect[0])

	// create keys to get the full vendor
	retKeys := make([]*datastore.Key, len(intersect))
	for i, k := range intersect {
		retKeys[i] = datastore.NewKey(ctx, "Vendor", "", k, nil)
		fmt.Printf("Key: %v\n", retKeys[i].ID())
	}

	// look up all the keys created
	retVendors := make([]*Vendor, len(retKeys))
	err2 := db.client.GetMulti(ctx, retKeys, retVendors)
	if err2 != nil {
		fmt.Printf("error here: %v\n", err2)
		return nil, fmt.Errorf("datastoredb: could not find nearby vendors: %v", err)
	}

	for i, l := range retVendors {
		retVendors[i] = l
	}

	return retVendors, nil
}

func (db *datastoreDB) ListVendors() ([]*Vendor, error) {
	ctx := context.Background()
	vendors := make([]*Vendor, 0)
	q := datastore.NewQuery("Vendor").
		Order("Name")

	_, err := db.client.GetAll(ctx, q, &vendors)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list books: %v", err)
	}

	/*
		for i, k := range keys {
			vendors[i].ID = k.ID()
		}
	*/

	return vendors, nil
}
