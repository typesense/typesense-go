package typesense

type StemmingInterface interface {
	Dictionaries() StemmingDictionariesInterface
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
