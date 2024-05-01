package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
)

type appService struct {
	repo Repository
}

type Service interface {
	getProducts(c fiber.Ctx) error
	getProduct(c fiber.Ctx) error
	createProduct(c fiber.Ctx) error
	updateProduct(c fiber.Ctx) error
	deleteProduct(c fiber.Ctx) error
}

func GetServiceInstance(r Repository) Service {
	return &appService{repo: r}
}

func (a *appService) getProducts(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	products, err := a.repo.findAll(ctx)

	return response(products, http.StatusOK, err, c)
}

func (a *appService) getProduct(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return response(nil, http.StatusUnprocessableEntity, errors.New("Id is not defined"), c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	travel, err := a.repo.findOne(ctx, id)

	return response(travel, http.StatusOK, err, c)
}

func (a *appService) createProduct(c fiber.Ctx) error {
	var p Product
	if err := c.Bind().Body(&p); err != nil {
		return response(p, http.StatusUnprocessableEntity, err, c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := a.repo.insertOne(ctx, &p)
	return response(p, http.StatusOK, err, c)
}

func (a *appService) updateProduct(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return response(nil, http.StatusUnprocessableEntity, errors.New("Id is not defined"), c)
	}

	var p Product
	if err := c.Bind().Body(&p); err != nil {
		return response(p, http.StatusUnprocessableEntity, err, c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := a.repo.updateOne(ctx, id, &p)
	return response(nil, http.StatusNoContent, err, c)
}

func (a *appService) deleteProduct(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return response(nil, http.StatusUnprocessableEntity, errors.New("Id is not defined"), c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := a.repo.deleteOne(ctx, id)
	return response(nil, http.StatusNoContent, err, c)
}

func response(data interface{}, httpStatus int, err error, c fiber.Ctx) error {
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"error": err.Error(),
		})
	} else {
		if data != nil {
			return c.Status(httpStatus).JSON(data)
		} else {
			return c.Status(http.StatusInternalServerError).JSON(nil)
		}
	}
}
