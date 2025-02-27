package controller

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/openziti/zrok/rest_model_zrok"
	"github.com/openziti/zrok/rest_server_zrok/operations/metadata"
	"github.com/sirupsen/logrus"
)

func overviewHandler(_ metadata.OverviewParams, principal *rest_model_zrok.Principal) middleware.Responder {
	tx, err := str.Begin()
	if err != nil {
		logrus.Errorf("error starting transaction: %v", err)
		return metadata.NewOverviewInternalServerError()
	}
	defer func() { _ = tx.Rollback() }()
	envs, err := str.FindEnvironmentsForAccount(int(principal.ID), tx)
	if err != nil {
		logrus.Errorf("error finding environments for '%v': %v", principal.Email, err)
		return metadata.NewOverviewInternalServerError()
	}
	var out rest_model_zrok.EnvironmentSharesList
	for _, env := range envs {
		shrs, err := str.FindSharesForEnvironment(env.Id, tx)
		if err != nil {
			logrus.Errorf("error finding shares for environment '%v': %v", env.ZId, err)
			return metadata.NewOverviewInternalServerError()
		}
		es := &rest_model_zrok.EnvironmentShares{
			Environment: &rest_model_zrok.Environment{
				Address:     env.Address,
				CreatedAt:   env.CreatedAt.UnixMilli(),
				Description: env.Description,
				Host:        env.Host,
				UpdatedAt:   env.UpdatedAt.UnixMilli(),
				ZID:         env.ZId,
			},
		}

		for _, shr := range shrs {
			feEndpoint := ""
			if shr.FrontendEndpoint != nil {
				feEndpoint = *shr.FrontendEndpoint
			}
			feSelection := ""
			if shr.FrontendSelection != nil {
				feSelection = *shr.FrontendSelection
			}
			beProxyEndpoint := ""
			if shr.BackendProxyEndpoint != nil {
				beProxyEndpoint = *shr.BackendProxyEndpoint
			}
			es.Shares = append(es.Shares, &rest_model_zrok.Share{
				Token:                shr.Token,
				ZID:                  shr.ZId,
				ShareMode:            shr.ShareMode,
				BackendMode:          shr.BackendMode,
				FrontendSelection:    feSelection,
				FrontendEndpoint:     feEndpoint,
				BackendProxyEndpoint: beProxyEndpoint,
				Reserved:             shr.Reserved,
				CreatedAt:            shr.CreatedAt.UnixMilli(),
				UpdatedAt:            shr.UpdatedAt.UnixMilli(),
			})
		}
		out = append(out, es)
	}
	return metadata.NewOverviewOK().WithPayload(out)
}
