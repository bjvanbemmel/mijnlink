package service

type URLService struct {
	IndexService IndexService
}

func (s URLService) SaveUrl(url string) (string, error) {
	return s.IndexService.SaveValue(url)
}

func (s URLService) GetURLByKey(key string) (string, error) {
	return s.IndexService.GetValueByKey(key)
}
