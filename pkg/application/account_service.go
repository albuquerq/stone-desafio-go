package application

import (
	"github.com/sirupsen/logrus"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

type accountService struct {
	log                           *logrus.Entry
	accountRepo                   account.Repository
	accountCreateValidator        account.ValidationStrategy
	accountBalanceUpdateValidator account.ValidationStrategy
}

// NewAccountService returns a new application implementation for account.Service.
func NewAccountService(
	repositoryRegistry domain.RepositoryRegistry,
	accountCreateValid account.ValidationStrategy,
	accountBalanceUpdateValid account.ValidationStrategy,
) account.Service {

	return &accountService{
		accountRepo:                   repositoryRegistry.AccountRepository(),
		accountCreateValidator:        accountCreateValid,
		accountBalanceUpdateValidator: accountBalanceUpdateValid,
		log:                           common.Logger().WithField("source", "accountService"),
	}
}

func (as *accountService) checkDependencies() {
	if as.accountCreateValidator == nil {
		as.log.Fatal("AccountCreationValidator strategy is not defined")
	}

	if as.accountBalanceUpdateValidator == nil {
		as.log.Fatal("BalanceUpdateValidator strategy is not defined")
	}

	if as.accountRepo == nil {
		as.log.Fatal("AccountRepository is not defined")
	}
}

func (as *accountService) CreateAccount(ac *account.Account) error {
	as.checkDependencies()
	log := as.log.WithField("op", "CreateAccount")

	ac.ID = as.accountRepo.GenerateIdentifier()

	err := ac.Validate(as.accountCreateValidator)
	if err != nil {
		log.WithError(err).WithField("accountID", ac.ID).Error("invalid data for account creation")
		return err
	}

	ac.Secret = common.HashSecret(ac.Secret)

	err = as.accountRepo.Store(ac)
	if err != nil {
		log.WithError(err).WithField("accountID", ac.ID).Error("error storing account data")
		return err
	}

	return nil
}

func (as *accountService) UpdateBalance(ac *account.Account) error {
	as.checkDependencies()
	log := as.log.WithField("op", "UpdateBalance")

	err := ac.Validate(as.accountBalanceUpdateValidator)
	if err != nil {
		log.WithError(err).WithField("accountID", ac.ID).Error("invalid data for balance update")
		return err
	}

	err = as.accountRepo.UpdateBalance(*ac)
	if err != nil {
		log.WithError(err).WithField("accountID", ac.ID).Error("error updating account balance")
		return err
	}

	return nil
}

func (as *accountService) AccountBalance(acID string) (balance account.BalanceValue, err error) {
	as.checkDependencies()
	log := as.log.WithField("op", "AccountBalance")

	ac, err := as.accountRepo.GetByID(acID)
	if err != nil {
		log.WithError(err).WithField("accountID", acID).Error("error retrieving account data")
		return
	}

	balance = account.BalanceValue{
		Balance: ac.Balance,
	}
	return
}

func (as *accountService) GetAccount(acID string) (ac account.Account, err error) {
	log := as.log.WithField("op", "GetAccount")

	ac, err = as.accountRepo.GetByID(acID)
	if err != nil {
		log.WithError(err).Error("account not found")
	}
	return
}

func (as *accountService) ListAccounts() ([]account.Account, error) {
	as.checkDependencies()
	log := as.log.WithField("op", "ListAccounts")

	acs, err := as.accountRepo.ListAll()
	if err != nil {
		log.WithError(err).Error("error when listing registered accounts")
		return nil, err
	}

	return acs, nil
}
