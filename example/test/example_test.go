package test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"strings"
	"testing"
	"time"

	"github.com/go-piv/piv-go/v2/piv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/buglloc/yubictld/pkg/yubictl"
)

type ExampleTestSuite struct {
	suite.Suite
	yc *yubictl.Yubikey
	yk *piv.YubiKey
}

func (s *ExampleTestSuite) setupYubictl() {
	svc := yubictl.NewSvcClient("http://localhost:3000")
	yk, err := svc.Acquire(context.Background())
	s.Require().NoError(err)

	s.T().Logf("acquired yubikey: %d", yk.Serial())
	s.yc = yk
}

func (s *ExampleTestSuite) setupYubikey() {
	cards, err := piv.Cards()
	require.NoError(s.T(), err)
	var yk *piv.YubiKey

	for _, card := range cards {
		if !strings.Contains(strings.ToLower(card), "yubikey") {
			continue
		}

		if yk, err = piv.Open(card); err != nil {
			s.T().Logf("unable to open yubikey %s: %v", card, err)
			continue
		}

		serial, err := yk.Serial()
		s.Require().NoError(err)

		if serial != s.yc.Serial() {
			_ = yk.Close()
			s.T().Logf("skip yubikey: %d", serial)
			continue
		}

		s.T().Logf("open yubikey: %d", serial)
		s.yk = yk
		break
	}
}

func (s *ExampleTestSuite) SetupSuite() {
	s.setupYubictl()
	s.setupYubikey()
}

func (s *ExampleTestSuite) TearDownSuite() {
	s.T().Logf("release yubikey: %d", s.yc.Serial())
	err := s.yc.Release(context.Background())
	s.Assert().NoError(err)

	s.T().Logf("close yubikey")
	err = s.yk.Close()
	s.Assert().NoError(err)
}

func (s *ExampleTestSuite) SetupTest() {
	err := s.yk.Reset()
	s.Require().NoError(err)
}

func (s *ExampleTestSuite) TestExample() {
	key := piv.Key{
		Algorithm:   piv.AlgorithmEC256,
		PINPolicy:   piv.PINPolicyAlways,
		TouchPolicy: piv.TouchPolicyAlways,
	}

	pub, err := s.yk.GenerateKey(piv.DefaultManagementKey, piv.SlotAuthentication, key)
	s.Require().NoError(err)

	auth := piv.KeyAuth{PIN: piv.DefaultPIN}
	priv, err := s.yk.PrivateKey(piv.SlotAuthentication, pub, auth)
	s.Require().NoError(err)

	signer, ok := priv.(crypto.Signer)
	s.Require().True(ok)

	data := sha256.Sum256([]byte("foo"))

	for {
		s.Run("should_ok", func() {
			err := s.yc.Touch(context.Background(), yubictl.TouchWithDelay(time.Second))
			s.Require().NoError(err)

			signature, err := signer.Sign(rand.Reader, data[:], crypto.SHA256)
			s.Require().NoError(err)
			s.Require().NotEmpty(signature)
			time.Sleep(time.Second)
		})
	}

	//s.Run("should_fail", func() {
	//	signature, err := signer.Sign(rand.Reader, data[:], crypto.SHA256)
	//	s.Require().Error(err)
	//	s.Require().Empty(signature)
	//})
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
