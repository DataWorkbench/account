package handler

import (
	"context"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/common/utils/signer"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func ValidateRequestSignature(ctx context.Context, req *accountpb.ValidateRequestSignatureRequest) (*executor.AccessKey, error) {
	s := signer.CreateSigner(req.ReqUserAgent)
	secretAccessKey, err := source.SelectSource(req.ReqSource, cfg, ctx).GetSecretAccessKey(req.ReqAccessKeyId)
	if err != nil {
		return nil, err
	}
	s.Init(secretAccessKey.AccessKeyID, secretAccessKey.SecretAccessKey, "")
	signature := s.CalculateSignature(req)
	if signature != req.ReqSignature {
		logger.Info().String(qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature).String(), "").Fire()
		return nil, qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature)
	}
	return secretAccessKey, nil
}
