package collect

import (
	"context"
	"cspm/pkg/platform/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"go.uber.org/zap"
)

type Acm struct{}

type AcmResult struct {
	CertificateDetails []acmTypes.CertificateDetail
}

func (*Acm) GetConfig(credentials aws.CredentialsProvider, regions []string) (types.RegionCollectMap, error) {
	result := types.RegionCollectMap{}
	for _, region := range regions {
		instance, err := getCertificates(credentials, region)
		if err != nil {
			return nil, err
		}

		result[region] = &AcmResult{CertificateDetails: instance}
	}

	return result, nil
}

// getCertificates 获取AWS的ACM Certificates
// https://docs.aws.amazon.com/acm/latest/APIReference/API_ListCertificates.html
func getCertificates(credentials aws.CredentialsProvider, region string) ([]acmTypes.CertificateDetail, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region),
		config.WithCredentialsProvider(credentials))
	if err != nil {
		return nil, err
	}
	svc := acm.NewFromConfig(cfg)
	request := acm.ListCertificatesInput{}
	out, err := svc.ListCertificates(context.TODO(), &request)
	if err != nil {
		return nil, err
	}

	zap.S().Debug("%v", out)

	var certificates []acmTypes.CertificateDetail
	for _, item := range out.CertificateSummaryList {
		certificate, err := svc.DescribeCertificate(context.TODO(), &acm.DescribeCertificateInput{
			CertificateArn: item.CertificateArn,
		})
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, *certificate.Certificate)
	}

	return certificates, nil
}
