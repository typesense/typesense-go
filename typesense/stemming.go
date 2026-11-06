package typesense

type StemmingInterface interface {
	// Dictionaries manages stemming dictionaries, which map words to their root forms.
	//
	// See: https://typesense.org/docs/latest/api/stemming.html
	Dictionaries() StemmingDictionariesInterface
	// Dictionary returns a handle for the stemming dictionary with the given id.
	//
	// See: https://typesense.org/docs/latest/api/stemming.html
	Dictionary(dictionaryId string) StemmingDictionaryInterface
}

type stemming struct {
	apiClient APIClientInterface
}

func (s *stemming) Dictionaries() StemmingDictionariesInterface {
	return &stemmingDictionaries{apiClient: s.apiClient}
}

func (s *stemming) Dictionary(dictionaryId string) StemmingDictionaryInterface {
	return &stemmingDictionary{apiClient: s.apiClient, dictionaryId: dictionaryId}
}
