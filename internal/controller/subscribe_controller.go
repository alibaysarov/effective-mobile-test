package controller

import (
	"database/sql"
	"effective-mobile/internal/dto"
	"effective-mobile/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscribeController struct {
	subscribeService service.SubscribeService
}

func NewSubscribeController(subscribeService service.SubscribeService) *SubscribeController {
	return &SubscribeController{
		subscribeService: subscribeService,
	}
}

// GetAll получает все подписки
// @Summary Получить все подписки
// @Description Получить список всех подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param start_date_from query string false "начало подписки от:" default(null)
// @Param start_date_to query string false "начало подписки до:" default(null)
// @Param service_name query string false "название сервиса:" default(null)
// @Success 200 {object} map[string]interface{} "Список подписок"
// @Router /subscriptions [get]
func (cnt *SubscribeController) GetAll(c *gin.Context) {
	var request dto.Filter

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, reqErr := cnt.subscribeService.GetAll(request)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetSum получает сумму всех подписок
// @Summary Получить сумму всех подписок
// @Description Получить список всех подписок
// @Tags subscriptions/total
// @Accept json
// @Produce json
// @Param start_date_from query string false "начало подписки от:" default(null)
// @Param start_date_to query string false "начало подписки до:" default(null)
// @Param service_name query string false "название сервиса:" default(null)
// @Success 200 {object} map[string]int "Сумма цен всех подписок"
// @Router /subscriptions/total [get]
func (cnt *SubscribeController) GetSum(c *gin.Context) {
	var request dto.Filter

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, reqErr := cnt.subscribeService.GetSum(request)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Create создает новую подписку
// @Summary Создать новую подписку
// @Description Создать новую подписку с указанными данными
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.SubscriptionDto true "Данные подписки"
// @Success 201 {object} map[string]interface{} "Созданная подписка"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (cnt *SubscribeController) Create(c *gin.Context) {
	var data dto.SubscriptionDto

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row, err := cnt.subscribeService.Create(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": row})
}

// Update обновляет существующую подписку
// @Summary Обновить подписку
// @Description Обновить данные существующей подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body dto.SubscriptionDto true "Обновленные данные подписки"
// @Success 200 {object} map[string]interface{} "Обновленная подписка"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (cnt *SubscribeController) Update(c *gin.Context) {
	var data dto.SubscriptionDto

	id := c.Param("id")

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем ID из URL параметра
	data.ID = id

	row, err := cnt.subscribeService.Update(data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": row})
}

// Delete удаляет подписку
// @Summary Удалить подписку
// @Description Удалить существующую подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string "Подписка успешно удалена"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (cnt *SubscribeController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := cnt.subscribeService.Delete(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("subscription with id %s not found", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// GetOne получает одну подписку
// @Summary Получить подписку по ID
// @Description Получить данные одной подписки по указанному ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]interface{} "Данные подписки"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (cnt *SubscribeController) GetOne(c *gin.Context) {
	id := c.Param("id")

	row, err := cnt.subscribeService.GetOne(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": row})
}
