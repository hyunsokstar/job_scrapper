package mydict

import "errors"

type Dictionary map[string]string

// var errWordExists = errors.New("That word already exists")
// var errNotFound = errors.New("Not Found")
var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errWordExists = errors.New("That word already exists")
)

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word] // 2번째 exists 는 boolean <=> 존재 유무
	if exists {
		return value, nil
	}

	// word 가 없으면 에러 정보 errNotFound 를 전달
	return "", errNotFound
}

// 딕셔너리에 단어 추가 version2 by switch
func (d Dictionary) Add(word, def string) error {
	// 단어 찾기
	// 단어가 존재하면 err 로 nil 이 넘어오고 존재하지 않으면 errNotFound 에러가 넘어 온다.
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

// update dictionary
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
