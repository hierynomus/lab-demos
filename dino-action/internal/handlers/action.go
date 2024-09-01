package handlers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog/log"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/config"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/store"
	"golang.org/x/exp/rand"
)

func NextActionHandler(cfg config.Config, store store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		name_param := utils.CopyString(c.Params("name"))
		// url decode name
		name, err := url.PathUnescape(name_param)
		if err != nil {
			log.Info().Err(err).Msg("Failed unescape name parameter")
			return err
		}

		actions, err := store.Get(c.UserContext(), name)
		if err != nil {
			log.Info().Err(err).Msg("Failed to get Actions")
			c.Status(fiber.StatusBadRequest)
			return err
		}

		// Random action
		action := actions[rand.Intn(len(actions))]

		return c.JSON(action)
	}
}

func AllActionsHandler(cfg config.Config, store store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		name_param := utils.CopyString(c.Params("name"))
		// url decode name
		name, err := url.PathUnescape(name_param)
		if err != nil {
			log.Info().Err(err).Msg("Failed unescape name parameter")
			return err
		}

		actions, err := store.Get(c.UserContext(), name)
		if err != nil {
			log.Info().Err(err).Msg("Failed to get Actions")
			c.Status(fiber.StatusBadRequest)
			return err
		}

		return c.JSON(actions)
	}
}
