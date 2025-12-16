package dependency

import (
	"good-todo-go/internal/infrastructure/database"
	"good-todo-go/internal/infrastructure/environment"
	"good-todo-go/internal/infrastructure/repository"
	"good-todo-go/internal/pkg"
	"good-todo-go/internal/presentation/public/controller"
	"good-todo-go/internal/presentation/public/presenter"
	"good-todo-go/internal/usecase"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// environment
	container.Provide(environment.LoadConfig)

	// infrastructure
	container.Provide(database.NewEntClient)

	// pkg
	container.Provide(func(cfg *environment.Config) *pkg.JWTService {
		return pkg.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn, cfg.JWTRefreshExpiresIn)
	})
	container.Provide(pkg.NewUUIDGenerator)

	// repository
	container.Provide(repository.NewAuthRepository)

	// usecase
	container.Provide(usecase.NewAuthInteractor)

	// presenter
	container.Provide(presenter.NewAuthPresenter)

	// controller
	container.Provide(controller.NewAuthController)

	return container
}
