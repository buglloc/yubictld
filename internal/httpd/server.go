package httpd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/gofiber/fiber/v2/middleware/logger"
	mrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/yubictld/internal/touchctl"
	"github.com/buglloc/yubictld/internal/xnet"
	"github.com/buglloc/yubictld/internal/ykman"
	"github.com/buglloc/yubictld/pkg/yubictl"
)

const DefaultAddr = "127.0.0.1:3000"

type Server struct {
	addr  string
	touch touchctl.Toucher
	yk    *ykman.YkMan
	app   *fiber.App
	log   zerolog.Logger
}

func NewServer(opts ...Option) (*Server, error) {
	l := log.With().
		Str("source", "server").
		Logger()

	s := &Server{
		addr: DefaultAddr,
		app: fiber.New(fiber.Config{
			ErrorHandler: errorHandler,
		}),
		log: l,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, s.init()
}

func (s *Server) ListenAndServe() error {
	ln, err := xnet.NewListener(s.addr)
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}
	defer func() {
		_ = ln.Close()
	}()

	return s.app.Listener(ln)
}

func (s *Server) Shutdown(_ context.Context) error {
	return s.app.Shutdown()
}

func (s *Server) init() error {
	s.app.Use(
		mrecover.New(),
		mlog.New(),
		requestid.New(),
	)

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("https://yubictl.prj.buglloc.com", http.StatusSeeOther)
	})

	s.app.Route("/v1", func(router fiber.Router) {
		router.Use(func(c *fiber.Ctx) error {
			if !c.Is("json") {
				return c.SendString("only JSON supported")
			}

			return c.Next()
		})

		router.Post("/acquire", func(c *fiber.Ctx) error {
			id := utils.UUIDv4()
			if s.yk == nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: "ykman not initialized",
				}
			}

			yk, err := s.yk.Acquire(id)
			if err != nil {
				if errors.Is(err, ykman.ErrNoFreeYubikey) {
					return &yubictl.ServiceError{
						HttpCode: fiber.StatusGone,
						Code:     yubictl.ServiceErrorNoFreeYubikey,
						Msg:      fmt.Sprintf("acuire free yubikey: %v", err),
					}
				}

				return fmt.Errorf("acquire free yubikey: %w", err)
			}

			s.log.Info().
				Str("client_id", id).
				Str("path", yk.Path()).
				Uint32("yk_serial", yk.Serial()).
				Msg("acquired yubikey")

			return c.JSON(yubictl.AcquireRsp{
				ID:     id,
				Serial: yk.Serial(),
			})
		})

		router.Post("/touch", func(c *fiber.Ctx) error {
			var req yubictl.TouchReq
			if err := c.BodyParser(&req); err != nil {
				return fmt.Errorf("parse body: %w", err)
			}

			yk, err := s.ykByClient(req.ID)
			if err != nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("lookup yubikey: %v", err),
				}
			}

			if s.touch == nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: "touchctl not initialized",
				}
			}

			port := yk.Port()
			if port == 0 {
				return &fiber.Error{
					Code:    fiber.StatusNotAcceptable,
					Message: "yubikey have no port configured",
				}
			}

			if err := s.touch.Touch(port, req.Delay, req.Duration); err != nil {
				s.log.Error().
					Str("client_id", req.ID).
					Str("path", yk.Path()).
					Uint32("yk_serial", yk.Serial()).
					Msg("touch failed")
				return err
			}

			s.log.Info().
				Str("client_id", req.ID).
				Str("path", yk.Path()).
				Uint32("yk_serial", yk.Serial()).
				Msg("touch yubikey")

			return nil
		})

		router.Post("/reboot", func(c *fiber.Ctx) error {
			var req yubictl.TouchReq
			if err := c.BodyParser(&req); err != nil {
				return fmt.Errorf("parse body: %w", err)
			}

			yk, err := s.ykByClient(req.ID)
			if err != nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("lookup yubikey: %v", err),
				}
			}

			if err := yk.Reboot(); err != nil {
				s.log.Error().
					Str("client_id", req.ID).
					Str("path", yk.Path()).
					Uint32("yk_serial", yk.Serial()).
					Msg("reboot failed")
				return err
			}

			s.log.Info().
				Str("client_id", req.ID).
				Str("path", yk.Path()).
				Uint32("yk_serial", yk.Serial()).
				Msg("reboot yubikey")

			return nil
		})

		router.Post("/ping", func(c *fiber.Ctx) error {
			var req yubictl.TouchReq
			if err := c.BodyParser(&req); err != nil {
				return fmt.Errorf("parse body: %w", err)
			}

			yk, err := s.ykByClient(req.ID)
			if err != nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("lookup yubikey: %v", err),
				}
			}

			s.log.Info().
				Str("client_id", req.ID).
				Str("path", yk.Path()).
				Uint32("yk_serial", yk.Serial()).
				Msg("ping yubikey")

			return yk.Ping()
		})

		router.Post("/release", func(c *fiber.Ctx) error {
			var req yubictl.TouchReq
			if err := c.BodyParser(&req); err != nil {
				return fmt.Errorf("parse body: %w", err)
			}

			yk, err := s.ykByClient(req.ID)
			if err != nil {
				return &fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("lookup yubikey: %v", err),
				}
			}

			if err := yk.Release(); err != nil {
				s.log.Error().
					Str("client_id", req.ID).
					Str("path", yk.Path()).
					Uint32("yk_serial", yk.Serial()).
					Msg("release failed")
				return err
			}

			s.log.Info().
				Str("client_id", req.ID).
				Str("path", yk.Path()).
				Uint32("yk_serial", yk.Serial()).
				Msg("release yubikey")

			return nil
		})
	})

	return nil
}

func (s *Server) ykByClient(clientID string) (*ykman.Yubikey, error) {
	if s.yk == nil {
		return nil, errors.New("ykman not initialized")
	}

	if clientID == "" {
		return nil, errors.New("clientID is empty")
	}

	return s.yk.ForClient(clientID)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var svcErr *yubictl.ServiceError
	var fiberErr *fiber.Error
	switch {
	case errors.As(err, &svcErr):
		err = ctx.Status(fiber.StatusInternalServerError).JSON(svcErr)

	case errors.As(err, &fiberErr):
		err = ctx.Status(code).JSON(yubictl.ServiceError{
			Code: yubictl.ServiceErrorInternalError,
			Msg:  err.Error(),
		})

	default:
		err = ctx.Status(fiber.StatusInternalServerError).JSON(yubictl.ServiceError{
			Code: yubictl.ServiceErrorInternalError,
			Msg:  err.Error(),
		})
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
