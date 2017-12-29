package petfinder

import (
	"encoding/json"
	"time"
)

type petIDResponse struct {
	Petfinder struct {
		PetIds struct {
			ID struct {
				T string `json:"$t"`
			} `json:"id"`
		} `json:"petIds"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

type petSingle struct {
	Options struct {
		Option interface{} `json:"option"`
	} `json:"options"`
	Status struct {
		T string `json:"$t"`
	} `json:"status"`
	Contact struct {
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
		Address1 struct {
			T string `json:"$t"`
		} `json:"address1"`
	} `json:"contact"`
	Age struct {
		T string `json:"$t"`
	} `json:"age"`
	Size struct {
		T string `json:"$t"`
	} `json:"size"`
	Media struct {
		Photos struct {
			Photo []struct {
				Size string `json:"@size"`
				T    string `json:"$t"`
				ID   string `json:"@id"`
			} `json:"photo"`
		} `json:"photos"`
	} `json:"media"`
	ID struct {
		T string `json:"$t"`
	} `json:"id"`
	ShelterPetID struct {
		T string `json:"$t"`
	} `json:"shelterPetId"`
	Breeds struct {
		Breed interface{} `json:"breed"`
	} `json:"breeds"`
	Name struct {
		T string `json:"$t"`
	} `json:"name"`
	Sex struct {
		T string `json:"$t"`
	} `json:"sex"`
	Description struct {
		T string `json:"$t"`
	} `json:"description"`
	Mix struct {
		T string `json:"$t"`
	} `json:"mix"`
	ShelterID struct {
		T string `json:"$t"`
	} `json:"shelterId"`
	LastUpdate struct {
		T time.Time `json:"$t"`
	} `json:"lastUpdate"`
	Animal struct {
		T string `json:"$t"`
	} `json:"animal"`
}

type petResponse struct {
	Petfinder struct {
		Pet    petSingle `json:"pet"`
		Header header    `json:"header"`
	} `json:"petfinder"`
}

type petFindResponse struct {
	Petfinder struct {
		LastOffset struct {
			T string `json:"lastOffset"`
		}
		Pets struct {
			Pet petSingle `json:"pet"`
		} `json:"pets"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

type petFindResponses struct {
	Petfinder struct {
		LastOffset struct {
			T string `json:"lastOffset"`
		}
		Pets struct {
			Pet []petSingle `json:"pet"`
		} `json:"pets"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

//Pet contains all the information about a single pet
type Pet struct {
	Status  string
	Options []string
	Contact struct {
		Address1 string
		Address2 string
		City     string
		Zip      string
		State    string
		Phone    string
		Email    string
		Fax      string
	}
	Age   string
	Size  string
	Media struct {
		Photos []struct {
			Size string
			URL  string
			ID   string
		}
	}
	ID           string
	ShelterPetID string
	Breeds       []string
	Name         string
	Sex          string
	Description  string
	Mix          string
	ShelterID    string
	LastUpdate   time.Time
	Animal       string
}

func (p *Pet) mapPetResponse(petR petSingle) {
	switch option := petR.Options.Option.(type) {
	case []interface{}:
		for _, o := range option {
			p.Options = append(p.Options, o.(map[string]interface{})["$t"].(string))
		}
	case interface{}:
		p.Options = []string{option.(map[string]interface{})["$t"].(string)}
	}

	p.Status = petR.Status.T

	p.Contact.Phone = petR.Contact.Phone.T
	p.Contact.State = petR.Contact.State.T
	p.Contact.Address1 = petR.Contact.Address1.T
	p.Contact.Address2 = petR.Contact.Address2.T
	p.Contact.Email = petR.Contact.Email.T
	p.Contact.City = petR.Contact.City.T
	p.Contact.Zip = petR.Contact.Zip.T
	p.Contact.Fax = petR.Contact.Fax.T

	p.Age = petR.Age.T
	p.Size = petR.Size.T

	p.ID = petR.ID.T
	p.ShelterPetID = petR.ShelterPetID.T

	switch breed := petR.Breeds.Breed.(type) {
	case []interface{}:
		for _, b := range breed {
			p.Breeds = append(p.Breeds, b.(map[string]interface{})["$t"].(string))
		}
	case interface{}:
		p.Breeds = []string{breed.(map[string]interface{})["$t"].(string)}
	}

	for _, photo := range petR.Media.Photos.Photo {
		photoStruct := struct {
			Size string
			URL  string
			ID   string
		}{Size: photo.Size, URL: photo.T, ID: photo.ID}
		p.Media.Photos = append(p.Media.Photos, photoStruct)
	}

	p.Name = petR.Name.T
	p.Sex = petR.Sex.T
	p.Description = petR.Description.T
	p.Mix = petR.Mix.T
	p.ShelterID = petR.ShelterID.T
	p.LastUpdate = petR.LastUpdate.T
	p.Animal = petR.Animal.T
}

//UnmarshalJSON is a custom unmarshaller for Pet
func (p *Pet) UnmarshalJSON(buf []byte) error {
	var petResp petResponse
	err := json.Unmarshal(buf, &petResp)
	if err != nil {
		return err
	}

	p.mapPetResponse(petResp.Petfinder.Pet)
	return nil
}

//Pets is a slice of pet
type Pets []Pet

//UnmarshalJSON is a custom unmarshaller for Pets
func (p *Pets) UnmarshalJSON(buf []byte) error {
	var pet Pet
	var petFindResp petFindResponse
	err := json.Unmarshal(buf, &petFindResp)
	if err != nil {
		var petFindResps petFindResponses
		err = json.Unmarshal(buf, &petFindResps)
		if err != nil {
			return err
		}

		for _, petR := range petFindResps.Petfinder.Pets.Pet {
			pet = Pet{}
			pet.mapPetResponse(petR)
			*p = append(*p, pet)
		}

		return nil
	}

	pet = Pet{}
	pet.mapPetResponse(petFindResp.Petfinder.Pets.Pet)
	*p = append(*p, pet)
	return nil
}
