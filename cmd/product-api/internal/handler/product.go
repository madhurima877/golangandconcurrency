package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"product-api/internal/model"
	"product-api/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Repo  *repository.ProductRepository
	Cache *repository.RedisCache
}

func NewProductHandler(repo *repository.ProductRepository, cache *repository.RedisCache) *ProductHandler {
	return &ProductHandler{Repo: repo, Cache: cache}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Repo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	key := fmt.Sprintf("product:%d", id)
	data, err := h.Cache.Get(c, key)
	if err == nil {
		var product model.Product
		if err := json.Unmarshal([]byte(data), &product); err == nil {
			c.JSON(200, product)
			return
		}
	}
	product, err := h.Repo.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}
	productData, _ := json.Marshal(product)
	_ = h.Cache.Set(c, key, string(productData))
	c.JSON(200, product)
}
