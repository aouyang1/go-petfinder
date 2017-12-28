package petfinder

import (
	"fmt"
	"os"
	"testing"
)

func fetchAPIKey() (string, error) {
	apiKey := os.Getenv("PETFINDER_API_KEY")
	if apiKey == "" {
		return apiKey, fmt.Errorf("Could not get petfinder api key from environment variable, PETFINDER_API_KEY")
	}
	return apiKey, nil

}

func TestNewPetFinderClient(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	if p.apiKey == "" {
		t.Errorf("API key not set for petfinder client")
	}
}

func TestBreedList(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	breedList, err := p.ListBreeds(Options{Animal: "dog"})
	if err != nil {
		t.Error(err)
	}

	if len(breedList) == 0 {
		t.Errorf("Did not get any breeds back")
	}
}

func TestGetRandomPetID(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	id, err := p.GetRandomPetID(Options{})
	if err != nil {
		t.Error(err)
	}

	if id == "" {
		t.Errorf("Did not get a valid ID")
	}

}

func TestGetRandomPet(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	opts := []Options{
		Options{Output: "full"},
		Options{Output: "basic"},
	}

	for _, o := range opts {
		pets, err := p.GetRandomPet(o)
		if err != nil {
			t.Error(err)
		}

		if pets.ID == "" {
			t.Errorf("Did not get a valid ID for pet, %+v with option, %+v\n", pets, o)
		}
	}

}

func TestGetPet(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	id, err := p.GetRandomPetID(Options{})
	if err != nil {
		t.Error(err)
	}

	opt := Options{ID: id}
	pets, err := p.GetPet(opt)
	if err != nil {
		t.Error(err)
	}

	if pets.ID == "" {
		t.Errorf("Did not get a valid ID for pet, %+v\n", pets)
	}
}

func TestFindPet(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	opts := []Options{
		Options{Location: "75093", Count: 10},
		Options{Location: "75093", Count: 1},
	}

	for _, o := range opts {
		pets, err := p.FindPet(o)
		if err != nil {
			t.Error(err)
		}
		if len(pets) != o.Count {
			t.Errorf("Did not receive %d pets back", o.Count)
		}
		for _, pet := range pets {
			if pet.ID == "" {
				t.Errorf("Did not get a valid ID for pet, %+v\n", pet)
			}
		}
	}
}

func TestFindShelter(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	opts := []Options{
		Options{Location: "75093", Count: 10},
		Options{Location: "75093", Count: 1},
	}

	for _, o := range opts {
		shelters, err := p.FindShelter(o)
		if err != nil {
			t.Error(err)
		}
		if len(shelters) != o.Count {
			t.Errorf("Did not receive %d shelters back", o.Count)
		}
		for _, shelter := range shelters {
			if shelter.ID == "" {
				t.Errorf("Did not get a valid ID for shelter, %+v\n", shelter)
			}
		}
	}
}

func TestGetShelter(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	shelter, err := p.GetShelter(Options{ID: "TX1203"})
	if err != nil {
		t.Error(err)
	}
	if shelter.ID == "" {
		t.Errorf("Did not get a valid ID for shelter, %+v\n", shelter)
	}
}

func TestGetShelterPets(t *testing.T) {
	apiKey, err := fetchAPIKey()
	if err != nil {
		t.Error(err)
	}

	p := NewPetFinderClient(apiKey)
	opts := []Options{
		Options{ID: "TX1203", Count: 1},
		Options{ID: "TX1203", Count: 10},
	}
	for _, o := range opts {
		pets, err := p.GetShelterPets(o)
		if err != nil {
			t.Error(err)
		}

		if len(pets) != o.Count {
			t.Errorf("Did not receive %d pets back for option, %+v\n", len(pets), o)
		}

		for _, pet := range pets {
			if pet.ID == "" {
				t.Errorf("Did not get a valid ID for pet, %+v\n", pet)
			}
		}
	}
}
