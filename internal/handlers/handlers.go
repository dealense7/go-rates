package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

type WebHandler struct {
	gasService      interfaces.GasService
	storeService    interfaces.StoreService
	categoryService interfaces.CategoryService
}

func NewWebHandler(gasService interfaces.GasService, storeService interfaces.StoreService, categoryService interfaces.CategoryService) *WebHandler {
	return &WebHandler{gasService: gasService, storeService: storeService, categoryService: categoryService}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func getDumpPath() string {
	layout := "2006-01-02"

	path, _ := os.Getwd()
	basePath := path + "/static/dump/dump-%s.sql"

	// Try today's date
	today := time.Now().Format(layout)
	todayPath := fmt.Sprintf(basePath, today)

	if fileExists(todayPath) {
		return today
	}

	// Try yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1).Format(layout)

	return yesterday
}

func (h *WebHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	// Gas
	gasRates, err := h.gasService.GetAll()
	if err != nil {
		paintError(c, err)
		return
	}
	gasRatesJson, _ := json.Marshal(gasRates)

	// Products
	storeItems, err := h.storeService.GetForSlider(ctx)
	if err != nil {
		paintError(c, err)
		return
	}
	storeItemsJson, _ := json.Marshal(storeItems)

	categoryItems, err := h.storeService.GetForCategorySlider(ctx)
	if err != nil {
		paintError(c, err)
		return
	}
	categoryItemsJson, _ := json.Marshal(categoryItems)

	c.HTML(200, "index", gin.H{
		"Title":         "Products List",
		"PageType":      "index",
		"DumpPath":      getDumpPath(),
		"gasRates":      string(gasRatesJson),
		"storeItems":    string(storeItemsJson),
		"categoryItems": string(categoryItemsJson),
	})
}

func (h *WebHandler) GetProductList(c *gin.Context) {
	ctx := c.Request.Context()

	pageUrl := c.Query("page")
	categoryId := c.Query("category_id")
	name := c.Query("search")
	page, _ := strconv.Atoi(pageUrl)
	category, _ := strconv.Atoi(categoryId)

	if page == 0 {
		page = 1
	}

	// Products
	Items, err := h.storeService.GetItemsList(ctx, page, category, name)
	if err != nil {
		paintError(c, err)
		return
	}

	ItemsJson, _ := json.Marshal(Items)

	totalItems, _ := h.storeService.GetItemsCount(ctx, category, name)

	// Categories
	categories, err := h.categoryService.GetItems()
	if err != nil {
		paintError(c, err)
		return
	}

	CategoryJson, _ := json.Marshal(categories)

	c.HTML(200, "items", gin.H{
		"Title":      "Items List",
		"PageType":   "items",
		"Items":      string(ItemsJson),
		"DumpPath":   getDumpPath(),
		"Categories": string(CategoryJson),
		"totalItems": totalItems,
	})
}

func (h *WebHandler) GetProductPrices(c *gin.Context) {

	productIdStr := c.Param("id")

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	item, err := h.storeService.GetProductById(productId)
	if err != nil {
		paintError(c, err)
		return
	}

	c.JSON(200, item)
}

func paintError(c *gin.Context, err error) {
	c.HTML(500, "error.html", gin.H{"error": err.Error()})
}
