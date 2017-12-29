package petfinder

import (
	"encoding/json"
)

type shelterSingle struct {
	Country struct {
		T string `json:"$t"`
	} `json:"country"`
	Longitude struct {
		T string `json:"$t"`
	} `json:"longitude"`
	Name struct {
		T string `json:"$t"`
	} `json:"name"`
	Phone struct {
		T string `json:"$t"`
	} `json:"phone"`
	State struct {
		T string `json:"$t"`
	} `json:"state"`
	Address2 struct {
		T string `json:"$t"`
	} `json:"address2"`
	Email struct {
		T string `json:"$t"`
	} `json:"email"`
	City struct {
		T string `json:"$t"`
	} `json:"city"`
	Zip struct {
		T string `json:"$t"`
	} `json:"zip"`
	Fax struct {
		T string `json:"$t"`
	} `json:"fax"`
	Latitude struct {
		T string `json:"$t"`
	} `json:"latitude"`
	ID struct {
		T string `json:"$t"`
	} `json:"id"`
	Address1 struct {
		T string `json:"$t"`
	} `json:"address1"`
}

type shelterResponse struct {
	Petfinder struct {
		Shelter shelterSingle `json:"shelter"`
		Header  header        `json:"header"`
	} `json:"petfinder"`
}

type shelterFindResponse struct {
	Petfinder struct {
		LastOffset struct {
			T string `json:"lastOffset"`
		}
		Shelters struct {
			Shelter shelterSingle `json:"shelter"`
		} `json:"shelters"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

type shelterFindResponses struct {
	Petfinder struct {
		LastOffset struct {
			T string `json:"lastOffset"`
		}
		Shelters struct {
			Shelter []shelterSingle `json:"shelter"`
		} `json:"shelters"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

//Shelter contains all information for a pet shelter
type Shelter struct {
	ID        string
	Name      string
	Longitude string
	Latitude  string
	Address1  string
	Address2  string
	City      string
	State     string
	Country   string
	Phone     string
	Email     string
	Zip       string
	Fax       string
}

func (s *Shelter) mapShelterResponse(shelterR shelterSingle) {
	s.ID = shelterR.ID.T
	s.Name = shelterR.Name.T
	s.Longitude = shelterR.Longitude.T
	s.Latitude = shelterR.Latitude.T
	s.Address1 = shelterR.Address1.T
	s.Address2 = shelterR.Address2.T
	s.City = shelterR.City.T
	s.State = shelterR.State.T
	s.Country = shelterR.Country.T
	s.Phone = shelterR.Phone.T
	s.Email = shelterR.Email.T
	s.Zip = shelterR.Zip.T
	s.Fax = shelterR.Fax.T
}

func (s *Shelter) UnmarshalJSON(buf []byte) error {
	var shelterResp shelterResponse
	err := json.Unmarshal(buf, &shelterResp)
	if err != nil {
		return err
	}

	s.mapShelterResponse(shelterResp.Petfinder.Shelter)
	return nil
}

type Shelters []Shelter

func (s *Shelters) UnmarshalJSON(buf []byte) error {
	var shelter Shelter
	var shelterFindResp shelterFindResponse
	err := json.Unmarshal(buf, &shelterFindResp)
	if err != nil {
		var shelterFindResps shelterFindResponses
		err = json.Unmarshal(buf, &shelterFindResps)
		if err != nil {
			return err
		}

		for _, shelterR := range shelterFindResps.Petfinder.Shelters.Shelter {
			shelter = Shelter{}
			shelter.mapShelterResponse(shelterR)
			*s = append(*s, shelter)
		}

		return nil
	}

	shelter = Shelter{}
	shelter.mapShelterResponse(shelterFindResp.Petfinder.Shelters.Shelter)
	*s = append(*s, shelter)

	return nil
}
