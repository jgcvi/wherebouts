package wherebouts

import (
	//	"errors"
	"fmt"
	"sort"
	"sync"
)

var _ VendorDatabase = &memoryDB{}

type memoryDB struct {
	mu      sync.Mutex
	nextID  int64
	vendors map[int64]*Vendor
}

func newMemoryDB() *memoryDB {
	return &memoryDB{
		vendors: make(map[int64]*Vendor),
		nextID:  1,
	}
}

func (db *memoryDB) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.vendors = nil
}

func (db *memoryDB) GetVendorByID(id int64) (*Vendor, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	vendor, ok := db.vendors[id]
	if !ok {
		return nil, fmt.Errorf("memorydb: book not found with ID %d", id)
	}
	return vendor, nil
}

func (db *memoryDB) AddVendor(v *Vendor) (id int64, err error) {
	fmt.Print("In Add! Memory")
	db.mu.Lock()
	defer db.mu.Unlock()

	//	v.ID = db.nextID
	//	db.vendors[v.ID] = v

	db.nextID++

	return 0, nil
}

func (db *memoryDB) UpdateVendor(v *Vendor) error {
	/*	if v.ID == 0 {
			return errors.New("memorydb: vendor with unassigned ID passed into updateVendor")
		}

	*/
	db.mu.Lock()
	defer db.mu.Unlock()

	//	db.vendors[v.ID] = v
	return nil
}

type vendorsByName []*Vendor

func (s vendorsByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s vendorsByName) Len() int           { return len(s) }
func (s vendorsByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (db *memoryDB) ListVendors() ([]*Vendor, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	vendors := make([]*Vendor, 0)
	for _, v := range db.vendors {
		vendors = append(vendors, v)
	}

	sort.Sort(vendorsByName(vendors))
	return vendors, nil
}

func (db *memoryDB) GetNearbyVendors(long float64, lat float64, dist float64) ([]*Vendor, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	vendors := make([]*Vendor, 0)
	for _, v := range db.vendors {
		vendors = append(vendors, v)
	}

	return vendors, nil
}

func (db *memoryDB) GetVendorByEmail(email string) (*Vendor, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, v := range db.vendors {
		if v.Email == email {
			return v, nil
		}
	}

	return nil, fmt.Errorf("could not find vendor")
}

func (db *memoryDB) ToggleAvailability(email string, lat float64, long float64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	vendor, ok := db.GetVendorByEmail(email)
	if ok != nil {
		return
	}

	vendor.Latitude = lat
	vendor.Longitude = long
	vendor.Open = !vendor.Open
}
