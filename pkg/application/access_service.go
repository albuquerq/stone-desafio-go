package application

import (
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/validator"
)

type accessService struct {
	log               *logrus.Entry
	accountRepository account.Repository
}

// NewAccessService returns a new access service.
func NewAccessService(repoRegistry domain.RepositoryRegistry) access.Service {
	return &accessService{
		log:               common.Logger().WithField("source", "accessService"),
		accountRepository: repoRegistry.AccountRepository(),
	}
}

func (as *accessService) Authenticate(cred access.Credential) (d access.Description, err error) {
	log := as.log.WithField("op", "Authenticate").WithField("cpf", cred.CPF)

	err = validator.ValidateAccessCredential(cred)
	if err != nil {
		return
	}

	ac, err := as.accountRepository.GetByCPF(cred.CPF)
	if err != nil {
		log.WithError(err).Error("no account found with this cpf")
		return
	}

	if !common.MatchHashAndSecret(ac.Secret, cred.Secret) { // match secrets using hash.
		err = errors.ErrInvalidAccessCredentials
		log.WithError(err).Error("mismatch secrets")
		return
	}

	d.AccountID = ac.ID
	d.CPF = ac.CPF
	d.Name = ac.Name

	return
}
