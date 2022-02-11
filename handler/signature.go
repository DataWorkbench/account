package handler

import (
	"context"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/signer"
	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
)

func getAccessKey(ctx context.Context, req *pbrequest.ValidateRequestSignature) (*executor.AccessKey, error) {
	secretAccessKey, err := cache.GetAccessKey(req.ReqAccessKeyId, req.ReqSource)
	if err != nil {
		if err == qerror.ResourceNotExists {
			logger.Debug().String("Access key not exist from cache", req.ReqAccessKeyId).Fire()
			return nil, qerror.AccessKeyNotExists.Format(req.ReqAccessKeyId)
		}
		logger.Error().String("Get access key from cache error", err.Error())
	}
	if secretAccessKey == nil {
		logger.Debug().String("Get access key from source", req.ReqAccessKeyId)
		secretAccessKey, err = source.SelectSource(cfg.Source, cfg, ctx).GetSecretAccessKey(req.ReqAccessKeyId)
		if err != nil {
			if err == qerror.ResourceNotExists {
				logger.Debug().String("Access key not exist from source", req.ReqAccessKeyId).Fire()
				if err = cache.CacheNotExistAccessKey(req.ReqAccessKeyId, req.ReqSource); err != nil {
					return nil, err
				}
				return nil, qerror.AccessKeyNotExists.Format(req.ReqAccessKeyId)
			}
			return nil, err
		}
		logger.Debug().String("Get access key not exist from source", req.ReqAccessKeyId).Fire()
		if err = cache.CacheAccessKey(secretAccessKey, req.ReqAccessKeyId, req.ReqSource); err != nil {
			logger.Warn().String("cache access key error", err.Error())
		}
	} else {
		logger.Debug().String("Get access key from cache successful", req.ReqAccessKeyId).Fire()
	}
	return secretAccessKey, nil
}

func ValidateRequestSignature(ctx context.Context, req *pbrequest.ValidateRequestSignature) (*executor.AccessKey, error) {
	if req.ReqSource == "" {
		req.ReqSource = executor.AccountExecutor.GetConf().Source
	}
	secretAccessKey, err := getAccessKey(ctx, req)
	if err != nil {
		return nil, err
	}
	s := signer.CreateSigner(req.ReqUserAgent)
	s.Init(secretAccessKey.AccessKeyID, secretAccessKey.SecretAccessKey, "")
	signature := s.CalculateSignature(req)
	if signature != req.ReqSignature {
		logger.Info().String(qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature).String(), "").Fire()
		return nil, qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature)
	}
	return secretAccessKey, nil
}
