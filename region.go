package cherrygo

// Region fields
type Region struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	RegionIso2 string    `json:"region_iso_2,omitempty"`
	BGP        RegionBGP `json:"bgp,omitempty"`
	Href       string    `json:"href,omitempty"`
}
