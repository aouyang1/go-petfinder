package petfinder

import (
	"encoding/json"
)

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

//Breeds is a slice of breeds for an animal type
type Breeds []string

//UnmarshalJSON is a custom unmarshaller for Breeds
func (b *Breeds) UnmarshalJSON(buf []byte) error {
	var breedList breedListResponse
	err := json.Unmarshal(buf, &breedList)
	if err != nil {
		return err
	}

	for _, breed := range breedList.Petfinder.Breeds.Breed {
		*b = append(*b, breed.T)
	}
	return nil
}
