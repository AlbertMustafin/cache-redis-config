package cache_redis_config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// GenerateRSAKey generates a new RSA key pair and stores them in files named
// "private_key.pem" and "public_key.pem".
func GenerateRSAKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	b, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	privatePEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	publicPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: b}
	privateOut, err := os.Create("private_key.pem")
	if err != nil {
		return err
	}
	defer privateOut.Close()
	if err := pem.Encode(privateOut, privatePEM); err != nil {
		return err
	}
	publicOut, err := os.Create("public_key.pem")
	if err != nil {
		return err
	}
	defer publicOut.Close()
	if err := pem.Encode(publicOut, publicPEM); err != nil {
		return err
	}
	return nil
}

// GenerateCA generates a self-signed certificate for use with the Redis
// proxy. It saves the certificate to a file named "ca.crt" and the private
// key to a file named "ca.key".
func GenerateCA() error {
	serialNumberLimit := new(big.Int).Lshift(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Cache Redis Config"},
		},
		IsCA: true,
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &template.Subject.PublicKey, template.PrivateKey)
	if err != nil {
		return err
	}

	certOut, err := os.Create("ca.crt")
	if err != nil {
		return err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	keyOut, err := os.Create("ca.key")
	if err != nil {
		return err
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(template.PrivateKey)}); err != nil {
		return err
	}

	return nil
}

// GetRedisClient returns a connection to the Redis server at the specified
// address and port.
func GetRedisClient(host string, port int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, strconv.Itoa(port)),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, "config")
}

// ParseHost returns the host and port as integers from a string in the form
// "host:port".
func ParseHost(host string) (int, error) {
	parts := strings.Split(host, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid host format")
	}
	host, port := parts[0], parts[1]
	if host == "" {
		return 0, errors.New("invalid host format")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return portInt, nil
}