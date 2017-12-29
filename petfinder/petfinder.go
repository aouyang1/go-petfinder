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

func (c Client) submitRequest(apiMethod string, opt Options) ([]byte, error) {
	var body []byte
	var response *http.Response
	var sleep time.Duration

	endpoint := c.baseURL + apiMethod
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return body, err
	}

	err = opt.validate()
	if err != nil {
		return body, err
	}

	q, err := query.Values(opt)
	if err != nil {
		return body, err
	}

	q["key"] = []string{c.apiKey}
	q["format"] = []string{c.format}
	request.URL.RawQuery = q.Encode()

	// submit request with retries
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

	if err != nil {
		return body, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

//ListBreeds returns a slice of breed names for a specified animal
//animal must be specified and one of the following:
//  barnyard, bird, cat, dog, horse, reptile, smallfurry
func (c Client) ListBreeds(opt Options) (Breeds, error) {
	var b Breeds

	// Check required options
	if opt.Animal == "" {
		return b, fmt.Errorf("Require animal type in options")
	}

	body, err := c.submitRequest("breed.list", opt)
	if err != nil {
		return b, err
	}
	err = json.Unmarshal(body, &b)
	return b, err
}

//GetRandomPetID return a string id of a random pet
//output option is overriden to id
func (c Client) GetRandomPetID(opt Options) (string, error) {
	var id string

	// Override for id output
	opt.Output = "id"

	body, err := c.submitRequest("pet.getRandom", opt)
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
//available options:
//  - animal 	string	optional 	type of animal (barnyard, bird, cat, dog, horse, reptile, smallfurry)
//  - breed 	string 	optional 	breed of animal (use breeds.list for a list of valid breeds)
//  - size 	string 	optional 	size of animal (S=small, M=medium, L=large, XL=extra-large)
//  - sex 	character 	optional 	M=male, F=female
//  - location 	string 	optional 	the ZIP/postal code or city and state the animal should be located (NOTE: the closest possible animal will be selected)
//  - shelterid 	string 	optional 	ID of the shelter that posted the pet
//  - output 	string 	optional How much of the pet record to return: basic, full
func (c Client) GetRandomPet(opt Options) (Pet, error) {
	var pet Pet

	// Override for id output
	if opt.Output == "id" || opt.Output == "" {
		return pet, fmt.Errorf("Output must be basic or full")
	}

	body, err := c.submitRequest("pet.getRandom", opt)
	if err != nil {
		return pet, err
	}
	err = json.Unmarshal(body, &pet)
	return pet, err
}

//GetPet retrieves information about a single Pet given an ID
//id option must be specified which indicates the pet ID
func (c Client) GetPet(opt Options) (Pet, error) {
	var pet Pet

	// Override for id output
	if opt.ID == "" {
		return pet, fmt.Errorf("Must specify pet ID")
	}

	body, err := c.submitRequest("pet.get", opt)
	if err != nil {
		return pet, err
	}
	err = json.Unmarshal(body, &pet)
	return pet, err
}

//FindPet returns a slice of Pets with their information given a location and other search options
//location option must be specified with represents a zip code or city/state
func (c Client) FindPet(opt Options) (Pets, error) {
	var pets Pets

	if opt.Location == "" {
		return pets, fmt.Errorf("Must specify zip code location string")
	}

	body, err := c.submitRequest("pet.find", opt)
	if err != nil {
		return pets, err
	}
	err = json.Unmarshal(body, &pets)
	return pets, err
}

//FindShelter resturns a slice of Shelter information given a location and other search options
//location option must be specified with represents a zip code or city/state
func (c Client) FindShelter(opt Options) (Shelters, error) {
	var shelters Shelters

	if opt.Location == "" {
		return shelters, fmt.Errorf("Must specify zip code or city state location string")
	}

	body, err := c.submitRequest("shelter.find", opt)
	if err != nil {
		return shelters, err
	}
	err = json.Unmarshal(body, &shelters)
	return shelters, err
}

//GetShelter retrieves Shelter information given an shelter ID
//id options must be specified which represents the shelter id
func (c Client) GetShelter(opt Options) (Shelter, error) {
	var shelter Shelter

	if opt.ID == "" {
		return shelter, fmt.Errorf("Must specify shelter id")
	}

	body, err := c.submitRequest("shelter.get", opt)
	if err != nil {
		return shelter, err
	}
	err = json.Unmarshal(body, &shelter)
	return shelter, err
}

//GetShelterPets retrieves a slice of Pet information for a shelter ID
//id options must be specified which represents the shelter id
func (c Client) GetShelterPets(opt Options) (Pets, error) {
	var pets Pets

	if opt.ID == "" {
		return pets, fmt.Errorf("Must specify pets id")
	}

	body, err := c.submitRequest("shelter.getPets", opt)
	if err != nil {
		return pets, err
	}
	err = json.Unmarshal(body, &pets)
	return pets, err
}
