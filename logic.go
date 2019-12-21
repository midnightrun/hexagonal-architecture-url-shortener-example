package shortener

import (
	"errors"
	"time"

	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrReadirectInvalid = errors.New("redirect invalid")
)

type redirectService struct {
	redirectRepository RedirectRepository
}

func NewRedirectService(redirectRepository RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepository,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepository.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return ErrReadirectInvalid
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepository.Store(redirect)

}
