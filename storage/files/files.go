package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

// all users have permission for read and write
const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.Wrap("can't save page", err) }()

	// формируем путь до директории куда будет сохраняться файл
	fPath := filepath.Join(s.basePath, p.UserName)

	// создаем все нужные директории на этом пути
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	// формируем имя файла в данном случае с помощью хэша, чтобы отвечать за уникальность файлов в одной директории (так как одинвковые названия файлов запрещены)
	fName, err := fileName(p)
	if err != nil {
		return err
	}

	// дописываем имя файла к пути
	fPath = filepath.Join(fPath, fName)

	// создаем файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// записываем в него страницу в нужно формате (gob подходит для сериализации и десереализации файла)
	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPage
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePath(filepath.Join(s.basePath, userName, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	const errMsg = "can't remove file"

	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap(errMsg, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("%s '%s'", errMsg, path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	const errMsg = "can't check file on existence"

	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap(errMsg, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	// check file on existence
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("%s %s", errMsg, path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePath(filePath string) (*storage.Page, error) {
	const errMsg = "can't decode page"

	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page
	err = gob.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return &p, err
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
