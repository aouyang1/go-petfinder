package petfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type header struct {
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

type breedListResponse struct {
	Petfinder struct {
		Breeds struct {
			Breed []struct {
				T string `json:"$t"`
			} `json:"breed"`
			Animal string `json:"@animal"`
		} `json:"breeds"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

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
		Pets struct {
			Pet petSingle `json:"pet"`
		} `json:"pets"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

type petFindResponses struct {
	Petfinder struct {
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
		Shelters struct {
			Shelter shelterSingle `json:"shelter"`
		} `json:"shelters"`
		Header header `json:"header"`
	} `json:"petfinder"`
}

type shelterFindResponses struct {
	Petfinder struct {
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

const (
	retryMax = 4
	minWait  = 600 * time.Millisecond
	maxWait  = 5 * time.Second
)

//Client is the Petfinder API client entrypoint
type Client struct {
	apiKey     string
	baseURL    string
	format     string
	HTTPClient *http.Client
}

//NewClient creates a new Petfinder API client as an entrypoint with a given api key
func NewClient(apiKey string) Client {
	p := Client{
		apiKey:     apiKey,
		baseURL:    "http://api.petfinder.com/",
		format:     "json",
		HTTPClient: &http.Client{},
	}
	return p
}

//Options are input arguments to the Petfind API
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

func (c Client) submitRequestWithRetry(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error
	var sleep time.Duration

	for i := 0; ; i++ {
		response, err = c.HTTPClient.Do(request)
		if err == nil || i == retryMax {
			break
		}
		sleep = time.Duration(math.Pow(2, float64(i)) * float64(minWait))
		if sleep > maxWait {
			sleep = maxWait
		}
		log.Printf("Retrying in %v with %d retries remaining", sleep, retryMax-i-1)
		time.Sleep(sleep)
	}

	return response, err
}

//ListBreeds returns a slice of breed names for a specified animal
func (c Client) ListBreeds(opt Options) ([]string, error) {
	var bl []string

	// Check required options
	if opt.Animal == "" {
		return bl, fmt.Errorf("Require animal type in options")
	}

	endpoint := c.baseURL + "breed.list"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return bl, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return bl, err
	}

	var breedList breedListResponse
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

//GetRandomPetID return a string id of a random pet
func (c Client) GetRandomPetID(opt Options) (string, error) {
	var id string

	endpoint := c.baseURL + "pet.getRandom"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return id, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return id, err
	}

	var petID petIDResponse
	err = json.Unmarshal(body, &petID)
	if err != nil {
		return id, err
	}

	return petID.Petfinder.PetIds.ID.T, nil
}

//GetRandomPet return a single random Pet
func (c Client) GetRandomPet(opt Options) (Pet, error) {
	var pet Pet

	endpoint := c.baseURL + "pet.getRandom"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return pet, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pet, err
	}

	var petResp petResponse
	err = json.Unmarshal(body, &petResp)
	if err != nil {
		return pet, err
	}

	pet.mapPetResponse(petResp.Petfinder.Pet)

	return pet, nil
}

//GetPet retrieves information about a single Pet given an ID
func (c Client) GetPet(opt Options) (Pet, error) {
	var pet Pet

	endpoint := c.baseURL + "pet.get"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return pet, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pet, err
	}

	var petResp petResponse
	err = json.Unmarshal(body, &petResp)
	if err != nil {
		return pet, err
	}

	pet.mapPetResponse(petResp.Petfinder.Pet)

	return pet, nil
}

//FindPet returns a slice of Pets with their information given a location and other search options
func (c Client) FindPet(opt Options) ([]Pet, error) {
	var pets []Pet

	endpoint := c.baseURL + "pet.find"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return pets, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pets, err
	}

	var pet Pet
	var petFindResp petFindResponse
	err = json.Unmarshal(body, &petFindResp)
	if err != nil {
		var petFindResps petFindResponses
		err = json.Unmarshal(body, &petFindResps)
		if err != nil {
			return pets, err
		}

		for _, petR := range petFindResps.Petfinder.Pets.Pet {
			pet = Pet{}
			pet.mapPetResponse(petR)
			pets = append(pets, pet)
		}

		return pets, nil
	}

	pet = Pet{}
	pet.mapPetResponse(petFindResp.Petfinder.Pets.Pet)
	pets = append(pets, pet)

	return pets, nil
}

//FindShelter resturns a slice of Shelter information given a location and other search options
func (c Client) FindShelter(opt Options) ([]Shelter, error) {
	var shelters []Shelter

	endpoint := c.baseURL + "shelter.find"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return shelters, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return shelters, err
	}

	var shelter Shelter
	var shelterFindResp shelterFindResponse
	err = json.Unmarshal(body, &shelterFindResp)
	if err != nil {
		var shelterFindResps shelterFindResponses
		err = json.Unmarshal(body, &shelterFindResps)
		if err != nil {
			return shelters, err
		}

		for _, shelterR := range shelterFindResps.Petfinder.Shelters.Shelter {
			shelter = Shelter{}
			shelter.mapShelterResponse(shelterR)
			shelters = append(shelters, shelter)
		}

		return shelters, nil
	}

	shelter = Shelter{}
	shelter.mapShelterResponse(shelterFindResp.Petfinder.Shelters.Shelter)
	shelters = append(shelters, shelter)

	return shelters, nil
}

//GetShelter retrieves Shelter information given an shelter ID
func (c Client) GetShelter(opt Options) (Shelter, error) {
	var shelter Shelter

	endpoint := c.baseURL + "shelter.get"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return shelter, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return shelter, err
	}

	var shelterResp shelterResponse
	err = json.Unmarshal(body, &shelterResp)
	if err != nil {
		return shelter, nil
	}

	shelter.mapShelterResponse(shelterResp.Petfinder.Shelter)

	return shelter, nil
}

//GetShelterPets retrieves a slice of Pet information for a shelter ID
func (c Client) GetShelterPets(opt Options) ([]Pet, error) {
	var pets []Pet

	endpoint := c.baseURL + "shelter.getPets"
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

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	response, err := c.submitRequestWithRetry(request)
	if err != nil {
		return pets, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pets, err
	}

	var pet Pet
	var petFindResp petFindResponse
	err = json.Unmarshal(body, &petFindResp)
	if err != nil {
		var petFindResps petFindResponses
		err = json.Unmarshal(body, &petFindResps)
		//out, _ := json.MarshalIndent(petFindResps, "", "  ")
		//fmt.Println(string(out))
		if err != nil {
			return pets, err
		}

		for _, petR := range petFindResps.Petfinder.Pets.Pet {
			pet = Pet{}
			pet.mapPetResponse(petR)
			pets = append(pets, pet)
		}

		return pets, nil
	}

	pet = Pet{}
	pet.mapPetResponse(petFindResp.Petfinder.Pets.Pet)
	pets = append(pets, pet)

	return pets, nil
}
