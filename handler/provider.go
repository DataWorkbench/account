package handler

import (
	"context"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
)

func GetProvider(ctx context.Context, req *pbrequest.GetProvider) (*pbmodel.Provider, error) {
	provider, err := cache.GetProvider(req.Name)
	if err != nil {
		return nil, err
	}
	if provider != nil {
		return provider, nil
	} else {
		return nil, qerror.InvalidProvider.Format(req.Name)
	}
}

func CreateProvider(ctx context.Context, req *pbrequest.CreateProvider) (*pbmodel.Provider, error) {
	provider := &pbmodel.Provider{
		Name:         req.Name,
		ClientId:     req.ClientId,
		ClientSecret: req.ClientSecret,
		TokenUrl:     req.TokenUrl,
		RedirectUrl:  req.RedirectUrl,
	}
	err := cache.CacheProvider(provider)
	if err != nil{
		return nil, err
	}else {
		return provider, nil
	}
}


func DeleteProvider(ctx context.Context, req *pbrequest.DeleteProvider) error {
	err := cache.DelProvider(req.Name)
	return err
}
