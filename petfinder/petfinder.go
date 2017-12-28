package petfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type Header struct {
	Timestamp struct {
		T time.Time `json:"$t"`
	} `json:"timestamp"`
	Status struct {
		Message struct {
			T string `json:"$t"`
		} `json:"message"`
		Code struct {
			T string `json:"$t"`
		} `json:"code"`
	} `json:"status"`
	Version struct {
		T string `json:"$t"`
	} `json:"version"`
}

type BreedListResponse struct {
	Petfinder struct {
		Breeds struct {
			Breed []struct {
				T string `json:"$t"`
			} `json:"breed"`
			Animal string `json:"@animal"`
		} `json:"breeds"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

type PetIDResponse struct {
	Petfinder struct {
		PetIds struct {
			ID struct {
				T string `json:"$t"`
			} `json:"id"`
		} `json:"petIds"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

type PetSingle struct {
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

type PetResponse struct {
	Petfinder struct {
		Pet    PetSingle `json:"pet"`
		Header Header    `json:"header"`
	} `json:"petfinder"`
}

type PetFindResponse struct {
	Petfinder struct {
		Pets struct {
			Pet PetSingle `json:"pet"`
		} `json:"pets"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

type PetFindResponses struct {
	Petfinder struct {
		Pets struct {
			Pet []PetSingle `json:"pet"`
		} `json:"pets"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

type Pet struct {
	Status  string
	Options []string
	Contact struct {
		Phone    string
		State    string
		Address2 string
		Email    string
		City     string
		Zip      string
		Fax      string
		Address1 string
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

func (p *Pet) mapPetResponse(petR PetSingle) {
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

type ShelterSingle struct {
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

type ShelterResponse struct {
	Petfinder struct {
		Shelter ShelterSingle `json:"shelter"`
		Header  Header        `json:"header"`
	} `json:"petfinder"`
}

type ShelterFindResponse struct {
	Petfinder struct {
		Shelters struct {
			Shelter ShelterSingle `json:"shelter"`
		} `json:"shelters"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

type ShelterFindResponses struct {
	Petfinder struct {
		Shelters struct {
			Shelter []ShelterSingle `json:"shelter"`
		} `json:"shelters"`
		Header Header `json:"header"`
	} `json:"petfinder"`
}

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

func (s *Shelter) mapShelterResponse(shelterR ShelterSingle) {
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

type PetFinderClient struct {
	apiKey  string
	baseURL string
	format  string
	Client  *http.Client
}

func NewPetFinderClient(apiKey string) PetFinderClient {
	p := PetFinderClient{
		apiKey:  apiKey,
		baseURL: "http://api.petfinder.com/",
		format:  "json",
		Client:  &http.Client{},
	}
	return p
}

type Options struct {
	ID          string `url:"id"`
	Animal      string `url:"animal"`
	Breed       string `url:"breed"`
	Size        string `url:"size"`
	Sex         string `url:"sex"`
	Location    string `url:"location"`
	Age         string `url:"age"`
	Offset      int    `url:"offset"`
	Count       int    `url:"count"`
	Output      string `url:"output"`
	ShelterID   string `url:"shelterid"`
	ShelterName string `url:"name"`
	Status      string `url:"status"`
}

func (o Options) validate() error {
	validAnimals := map[string]struct{}{
		"barnyard":   struct{}{},
		"bird":       struct{}{},
		"cat":        struct{}{},
		"dog":        struct{}{},
		"horse":      struct{}{},
		"reptile":    struct{}{},
		"smallfurry": struct{}{},
	}
	if o.Animal != "" {
		if _, ok := validAnimals[o.Animal]; !ok {
			return fmt.Errorf("Invalid animal specified")
		}
	}

	validSizes := map[string]struct{}{
		"S":  struct{}{},
		"M":  struct{}{},
		"L":  struct{}{},
		"XL": struct{}{},
	}
	if o.Size != "" {
		if _, ok := validSizes[o.Size]; !ok {
			return fmt.Errorf("Invalid size specified")
		}
	}

	validSex := map[string]struct{}{
		"M": struct{}{},
		"F": struct{}{},
	}
	if o.Sex != "" {
		if _, ok := validSex[o.Sex]; !ok {
			return fmt.Errorf("Invalid sex specified")
		}
	}

	validAges := map[string]struct{}{
		"Baby":   struct{}{},
		"Young":  struct{}{},
		"Adult":  struct{}{},
		"Senior": struct{}{},
	}
	if o.Age != "" {
		if _, ok := validAges[o.Age]; !ok {
			return fmt.Errorf("Invalid age specified")
		}
	}

	validOutputs := map[string]struct{}{
		"basic": struct{}{},
		"full":  struct{}{},
		"id":    struct{}{},
	}
	if o.Output != "" {
		if _, ok := validOutputs[o.Output]; !ok {
			return fmt.Errorf("Invalid output specified")
		}
	}

	validStatuses := map[string]struct{}{
		"A": struct{}{},
		"H": struct{}{},
		"P": struct{}{},
		"X": struct{}{},
	}
	if o.Status != "" {
		if _, ok := validStatuses[o.Status]; !ok {
			return fmt.Errorf("Invalid status specified")
		}
	}

	return nil
}

func (p PetFinderClient) ListBreeds(opt Options) ([]string, error) {
	var bl []string

	// Check required options
	if opt.Animal == "" {
		return bl, fmt.Errorf("Require animal type in options")
	}

	endpoint := p.baseURL + "breed.list"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return bl, err
	}

	err = opt.validate()
	if err != nil {
		return bl, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return bl, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return bl, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return bl, err
	}

	var breedList BreedListResponse
	err = json.Unmarshal(body, &breedList)
	if err != nil {
		return bl, err
	}

	bl = make([]string, 0, len(breedList.Petfinder.Breeds.Breed))
	for _, breed := range breedList.Petfinder.Breeds.Breed {
		bl = append(bl, breed.T)
	}

	return bl, err
}

func (p PetFinderClient) GetRandomPetID(opt Options) (string, error) {
	var id string

	endpoint := p.baseURL + "pet.getRandom"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return id, err
	}

	// Override for id output
	opt.Output = "id"

	err = opt.validate()
	if err != nil {
		return id, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return id, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return id, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return id, err
	}

	var petID PetIDResponse
	err = json.Unmarshal(body, &petID)
	if err != nil {
		return id, err
	}

	return petID.Petfinder.PetIds.ID.T, nil
}

func (p PetFinderClient) GetRandomPet(opt Options) (Pet, error) {
	var pet Pet

	endpoint := p.baseURL + "pet.getRandom"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return pet, err
	}

	// Override for id output
	if opt.Output == "id" || opt.Output == "" {
		return pet, fmt.Errorf("Output must be basic or full")
	}

	err = opt.validate()
	if err != nil {
		return pet, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return pet, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return pet, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pet, err
	}

	var petResponse PetResponse
	err = json.Unmarshal(body, &petResponse)
	if err != nil {
		return pet, err
	}

	pet.mapPetResponse(petResponse.Petfinder.Pet)

	return pet, nil
}

func (p PetFinderClient) GetPet(opt Options) (Pet, error) {
	var pet Pet

	endpoint := p.baseURL + "pet.get"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return pet, err
	}

	// Override for id output
	if opt.ID == "" {
		return pet, fmt.Errorf("Must specify pet ID")
	}

	err = opt.validate()
	if err != nil {
		return pet, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return pet, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return pet, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pet, err
	}

	var petResponse PetResponse
	err = json.Unmarshal(body, &petResponse)
	if err != nil {
		return pet, err
	}

	pet.mapPetResponse(petResponse.Petfinder.Pet)

	return pet, nil
}

func (p PetFinderClient) FindPet(opt Options) ([]Pet, error) {
	var pets []Pet

	endpoint := p.baseURL + "pet.find"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return pets, err
	}

	if opt.Location == "" {
		return pets, fmt.Errorf("Must specify zip code location string")
	}

	err = opt.validate()
	if err != nil {
		return pets, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return pets, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return pets, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pets, err
	}

	var pet Pet
	var petFindResponse PetFindResponse
	err = json.Unmarshal(body, &petFindResponse)
	if err != nil {
		var petFindResponses PetFindResponses
		err = json.Unmarshal(body, &petFindResponses)
		if err != nil {
			return pets, err
		}

		for _, petR := range petFindResponses.Petfinder.Pets.Pet {
			pet = Pet{}
			pet.mapPetResponse(petR)
			pets = append(pets, pet)
		}

		return pets, nil
	}

	pet = Pet{}
	pet.mapPetResponse(petFindResponse.Petfinder.Pets.Pet)
	pets = append(pets, pet)

	return pets, nil
}

func (p PetFinderClient) FindShelter(opt Options) ([]Shelter, error) {
	var shelters []Shelter

	endpoint := p.baseURL + "shelter.find"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return shelters, err
	}

	if opt.Location == "" {
		return shelters, fmt.Errorf("Must specify zip code or city state location string")
	}

	err = opt.validate()
	if err != nil {
		return shelters, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return shelters, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return shelters, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return shelters, err
	}

	var shelter Shelter
	var shelterFindResponse ShelterFindResponse
	err = json.Unmarshal(body, &shelterFindResponse)
	if err != nil {
		var shelterFindResponses ShelterFindResponses
		err = json.Unmarshal(body, &shelterFindResponses)
		if err != nil {
			return shelters, err
		}

		for _, shelterR := range shelterFindResponses.Petfinder.Shelters.Shelter {
			shelter = Shelter{}
			shelter.mapShelterResponse(shelterR)
			shelters = append(shelters, shelter)
		}

		return shelters, nil
	}

	shelter = Shelter{}
	shelter.mapShelterResponse(shelterFindResponse.Petfinder.Shelters.Shelter)
	shelters = append(shelters, shelter)

	return shelters, nil
}

func (p PetFinderClient) GetShelter(opt Options) (Shelter, error) {
	var shelter Shelter

	endpoint := p.baseURL + "shelter.get"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return shelter, err
	}

	if opt.ID == "" {
		return shelter, fmt.Errorf("Must specify shelter id")
	}

	err = opt.validate()
	if err != nil {
		return shelter, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return shelter, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return shelter, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return shelter, err
	}

	var shelterResponse ShelterResponse
	err = json.Unmarshal(body, &shelterResponse)
	if err != nil {
		return shelter, nil
	}

	shelter.mapShelterResponse(shelterResponse.Petfinder.Shelter)

	return shelter, nil
}

func (p PetFinderClient) GetShelterPets(opt Options) ([]Pet, error) {
	var pets []Pet

	endpoint := p.baseURL + "shelter.getPets"
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return pets, err
	}

	if opt.ID == "" {
		return pets, fmt.Errorf("Must specify pets id")
	}

	err = opt.validate()
	if err != nil {
		return pets, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return pets, err
	}

	q["key"] = []string{p.apiKey}
	q["format"] = []string{p.format}
	request.URL.RawQuery = q.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return pets, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pets, err
	}

	var pet Pet
	var petFindResponse PetFindResponse
	err = json.Unmarshal(body, &petFindResponse)
	if err != nil {
		var petFindResponses PetFindResponses
		err = json.Unmarshal(body, &petFindResponses)
		if err != nil {
			return pets, err
		}

		for _, petR := range petFindResponses.Petfinder.Pets.Pet {
			pet = Pet{}
			pet.mapPetResponse(petR)
			pets = append(pets, pet)
		}

		return pets, nil
	}

	pet = Pet{}
	pet.mapPetResponse(petFindResponse.Petfinder.Pets.Pet)
	pets = append(pets, pet)

	return pets, nil
}