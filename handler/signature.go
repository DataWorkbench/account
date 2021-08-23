package handler

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/DataWorkbench/account/executor"
	"github.com/DataWorkbench/account/internal/source"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/pkg/accountpb"
)

func CalculateSignature(stringToSign string, secretAccessKey string) string {
	h := hmac.New(sha256.New, []byte(secretAccessKey))
	h.Write([]byte(stringToSign))
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)
	return signature
}

func ValidateRequestSignature(ctx context.Context, req *accountpb.ValidateRequestSignatureRequest) (*executor.AccessKey, error) {
	secretAccessKey, err := source.SelectSource(req.ReqSource, cfg, ctx).GetSecretAccessKey(req.ReqAccessKeyId)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write([]byte(req.ReqBody))
	stringToSign := strings.ToUpper(req.ReqMethod) + "\n" + req.ReqPath + "\n" + req.ReqQueryString + "\n" + hex.EncodeToString(h.Sum(nil))
	signature := CalculateSignature(stringToSign, secretAccessKey.SecretAccessKey)
	if signature != req.ReqSignature {
		logger.Info().String(qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature).String(), "").Fire()
		return nil, qerror.ValidateSignatureFailed.Format(req.ReqSignature, signature)
	}
	return secretAccessKey, nil
}
