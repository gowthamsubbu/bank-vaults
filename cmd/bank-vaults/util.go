package main

import (
	"fmt"

	"github.com/banzaicloud/bank-vaults/pkg/kv"
	"github.com/banzaicloud/bank-vaults/pkg/kv/awskms"
	"github.com/banzaicloud/bank-vaults/pkg/kv/cloudkms"
	"github.com/banzaicloud/bank-vaults/pkg/kv/gcs"
	"github.com/banzaicloud/bank-vaults/pkg/kv/k8s"
	"github.com/banzaicloud/bank-vaults/pkg/kv/keyvault"
	"github.com/banzaicloud/bank-vaults/pkg/kv/s3"
	"github.com/banzaicloud/bank-vaults/pkg/vault"
	"github.com/spf13/viper"
)

func vaultConfigForConfig(cfg *viper.Viper) (vault.Config, error) {

	return vault.Config{
		SecretShares:    appConfig.GetInt(cfgSecretShares),
		SecretThreshold: appConfig.GetInt(cfgSecretThreshold),

		InitRootToken:  appConfig.GetString(cfgInitRootToken),
		StoreRootToken: appConfig.GetBool(cfgStoreRootToken),
	}, nil
}

func kvStoreForConfig(cfg *viper.Viper) (kv.Service, error) {

	if cfg.GetString(cfgMode) == cfgModeValueGoogleCloudKMSGCS {

		g, err := gcs.New(
			cfg.GetString(cfgGoogleCloudStorageBucket),
			cfg.GetString(cfgGoogleCloudStoragePrefix),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud storage kv store: %s", err.Error())
		}

		kms, err := cloudkms.New(g,
			cfg.GetString(cfgGoogleCloudKMSProject),
			cfg.GetString(cfgGoogleCloudKMSLocation),
			cfg.GetString(cfgGoogleCloudKMSKeyRing),
			cfg.GetString(cfgGoogleCloudKMSCryptoKey),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud kms kv store: %s", err.Error())
		}

		return kms, nil
	}

	if cfg.GetString(cfgMode) == cfgModeValueAWSKMS3 {
		s3, err := s3.New(
			cfg.GetString(cfgAWSS3Bucket),
			cfg.GetString(cfgAWSS3Prefix),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating AWS S3 kv store: %s", err.Error())
		}

		kms, err := awskms.New(s3, cfg.GetString(cfgAWSKMSKeyID))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS KMS ID kv store: %s", err.Error())
		}

		return kms, nil
	}

	if cfg.GetString(cfgMode) == cfgModeValueAzureKeyVault {
		kms, err := keyvault.New(cfg.GetString(cfgAzureKeyVaultName))
		if err != nil {
			return nil, fmt.Errorf("error creating Azure Key Vault kv store: %s", err.Error())
		}

		return kms, nil
	}

	if cfg.GetString(cfgMode) == cfgModeValueK8S {
		k8s, err := k8s.New(
			cfg.GetString(cfgK8SNamespace),
			cfg.GetString(cfgK8SSecret),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating K8S Secret kv store: %s", err.Error())
		}

		return k8s, nil
	}

	return nil, fmt.Errorf("Unsupported backend mode: '%s'", cfg.GetString(cfgMode))
}