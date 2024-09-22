package server

import (
	"context"
	"fmt"
	"net/http"
)

// Server описывает веб сервер
type Server interface {
	// Start запускает веб сервер
	Start() <-chan error
	// Stop останавливает веб сервер
	Stop() error
}

// impl реализация веб сервера
type impl struct {
	settings   Settings
	router     *gin.Engine
	server     *http.Server
	controller controllers.Controller
}

// Start запускает веб сервер
func (srv *impl) Start() <-chan error {
	ch := make(chan error, 1)

	go func() {
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- err
		}
	}()

	return ch
}

// Stop останавливает веб сервер
func (srv *impl) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.settings.ShutdownTimeout)
	defer cancel()
	if err := srv.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// NewServer создает новый веб сервер
func NewServer(appLog service.AppLogger, settings Settings) (Server, error) {
	set := &settings
	if err := set.validate(); err != nil {
		return nil, err
	}

	fmt.Println(settings)

	if settings.GIN.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if settings.GIN.UseLogger {
		r.Use(gin.Logger())
	}
	if settings.GIN.UseRecovery {
		r.Use(gin.Recovery())
	}

	if settings.ProfilingEnabled {
		pprof.Register(r)
	}

	err := settings.Controllers.Parse()
	if err != nil {
		return nil, err
	}

	//r.Use(ginmiddlewares.CORS(settings.CorsSettings))

	controller, err := controllers.NewController(appLog, settings.Controllers)
	if err != nil {
		return nil, err
	}

	srv := &impl{
		settings: settings,
		router:   r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", settings.Port),
			Handler: r,
		},
		controller: controller,
	}

	srv.routes()

	return srv, nil
}
