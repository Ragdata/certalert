package certificates

import (
	"fmt"

	pkcs12 "github.com/gi8lino/go-pkcs12"
	log "github.com/sirupsen/logrus"
)

// ExtractP12CertificatesInfo reads the P12 file, extracts certificate information, and returns a list of CertificateInfo
func ExtractTrustStoreCertificatesInfo(name string, certData []byte, password string, failOnError bool) ([]CertificateInfo, error) {
	var certInfoList []CertificateInfo

	// handleError is a helper function to handle failOnError
	handleError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warningf("Failed to extract certificate information: %v", errMsg)
		certInfoList = append(certInfoList, CertificateInfo{
			Name:  name,
			Type:  "p12",
			Error: errMsg,
		})
		return nil
	}

	// Decode the P12 data
	certs, err := pkcs12.DecodeTrustStore(certData, password)
	if err != nil {
		return certInfoList, handleError(fmt.Sprintf("Failed to decode P12 file '%s': %v", name, err))
	}

	var counter int

	// Extract certificates
	for _, cert := range certs {
		counter++
		subject := cert.Subject.CommonName
		if subject == "" {
			subject = fmt.Sprint(counter)
		}

		certInfo := CertificateInfo{
			Name:    name,
			Subject: subject,
			Epoch:   cert.NotAfter.Unix(),
			Type:    "truststore",
		}
		log.Debugf("Certificate '%s' expires on %s", certInfo.Subject, certInfo.ExpiryAsTime())

		certInfoList = append(certInfoList, certInfo)
	}

	return certInfoList, nil
}
