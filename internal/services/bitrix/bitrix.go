package bitrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ia-online-golang/internal/dto"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Структура EmailService для хранения настроек SMTP
type BitrixService struct {
	log     *logrus.Logger
	webhook string
}

type BitrixServiceI interface {
	GetLead(ctx context.Context, id_deal int64) (ReturnDataDeal, error)
	SendDeal(ctx context.Context, lead dto.LeadDTO, user dto.UserDTO) (ReturnDataCreate, error)
	SendContact(ctx context.Context, dto dto.LeadDTO) (ReturnDataCreate, error)
}

// Конструктор для создания нового экземпляра EmailService
func New(log *logrus.Logger, webhook string) *BitrixService {
	return &BitrixService{
		log:     log,
		webhook: webhook,
	}
}

func (b *BitrixService) GetLead(ctx context.Context, id_deal int64) (ReturnDataDeal, error) {
	const op = "BitrixService.GetLead"

	data := map[string]any{
		"ID": id_deal,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ReturnDataDeal{}, fmt.Errorf("%s: %v", op, err)
	}

	url := b.webhook + "crm.deal.get"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ReturnDataDeal{}, fmt.Errorf("%s: %v", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReturnDataDeal{}, fmt.Errorf("%s: %v", op, err)
	}

	// Попробуем сначала распарсить как ошибку
	var errData ErrorData
	if err := json.Unmarshal(body, &errData); err == nil && errData.Error != "" {
		return ReturnDataDeal{}, fmt.Errorf("%s: %s - %s", op, errData.Error, errData.ErrorDescription)
	}

	// Если ошибки нет — пробуем распарсить как успешный ответ
	var result ReturnDataDeal
	if err := json.Unmarshal(body, &result); err != nil {
		return ReturnDataDeal{}, fmt.Errorf("%s: %v", op, err)
	}

	return result, nil
}

func (b *BitrixService) SendDeal(ctx context.Context, lead dto.LeadDTO, user dto.UserDTO) (ReturnDataCreate, error) {
	op := "BitrixService.SendLead"

	contact, err := b.SendContact(ctx, lead)
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	var services []int
	if lead.IsInternet {
		services = append(services, 510)
	}
	if lead.IsCleaning {
		services = append(services, 512)
	}
	if lead.IsShipping {
		services = append(services, 514)
	}

	data := map[string]any{
		"fields": map[string]any{
			"TITLE":                 "Заявка с сайта ia-on.ru",
			"TYPE_ID":               "SALE",
			"STAGE_ID":              "C42:NEW",
			"IS_MANUAL_OPPORTUNITY": "Y",
			"CATEGORY_ID":           42,
			"CONTACT_ID":            contact.Result,
			"UF_CRM_1697646751446":  lead.Address,
			"UF_CRM_1743744405443":  services,
			"UF_CRM_1697294923031":  lead.Comment,
			"UF_CRM_1703703644316":  user.City,
			"UF_CRM_1697357613372":  user.Name,
			"UF_CRM_1700909419606":  user.PhoneNumber,
			"UF_CRM_1701035680304":  user.ID,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	url := b.webhook + "crm.deal.add"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	// Попробуем сначала распарсить как ошибку
	var errData ErrorData
	if err := json.Unmarshal(body, &errData); err == nil && errData.Error != "" {
		return ReturnDataCreate{}, fmt.Errorf("%s: %s - %s", op, errData.Error, errData.ErrorDescription)
	}

	// Если ошибки нет — пробуем распарсить как успешный ответ
	var result ReturnDataCreate
	if err := json.Unmarshal(body, &result); err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	return result, nil
}

func (b *BitrixService) SendContact(ctx context.Context, dto dto.LeadDTO) (ReturnDataCreate, error) {
	op := "BitrixService.SendContact"

	data := map[string]any{
		"fields": map[string]any{
			"NAME": dto.Name,
			"PHONE": []map[string]string{
				{
					"VALUE":      dto.PhoneNumber,
					"VALUE_TYPE": "WORK",
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	url := b.webhook + "crm.contact.add"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	// Попробуем сначала распарсить как ошибку
	var errData ErrorData
	if err := json.Unmarshal(body, &errData); err == nil && errData.Error != "" {
		return ReturnDataCreate{}, fmt.Errorf("%s: %s - %s", op, errData.Error, errData.ErrorDescription)
	}

	// Если ошибки нет — пробуем распарсить как успешный ответ
	var result ReturnDataCreate
	if err := json.Unmarshal(body, &result); err != nil {
		return ReturnDataCreate{}, fmt.Errorf("%s: %v", op, err)
	}

	return result, nil
}
