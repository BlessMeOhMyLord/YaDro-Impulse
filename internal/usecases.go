package internal

import (
	"context"
	"net"
	"strings"
)

type Repository interface {
	Add(ctx context.Context, newDns string) error
	Delete(ctx context.Context, dnsToDelete string) error
	List(ctx context.Context) ([]string, error)
}

type Service struct {
	Repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func dnsValidate(dns string) (string, error) {
	validatedDns := strings.Trim(dns, " \n\r\t")
	ip := net.ParseIP(validatedDns)
	if ip == nil {
		return "", ErrIsIncorrect
	}

	return validatedDns, nil
}

func (service *Service) Add(ctx context.Context, newDns string) (string, error) {
	dns, err := dnsValidate(newDns)
	if err != nil {
		return "", err
	}

	if err := service.Repository.Add(ctx, dns); err != nil {
		return "", err
	}

	return dns, nil
}

func (service *Service) Delete(ctx context.Context, dnsToDelete string) (string, error) {
	dns, err := dnsValidate(dnsToDelete)
	if err != nil {
		return "", err
	}

	if err := service.Repository.Delete(ctx, dns); err != nil {
		return "", err
	}
	return dns, nil
}

func (service *Service) List(ctx context.Context) ([]string, error) {
	list, err := service.Repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}
