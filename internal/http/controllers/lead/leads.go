package leads

import (
	"encoding/json"
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/context_keys"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/services/lead"
	"ia-online-golang/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type LeadController struct {
	log         *logrus.Logger
	validator   *validator.Validate
	LeadService lead.LeadServiceI
}

type LeadControllerI interface {
	SaveLead(w http.ResponseWriter, r *http.Request)
	Leads(w http.ResponseWriter, r *http.Request)
}

func New(log *logrus.Logger, validator *validator.Validate, leadService lead.LeadServiceI) *LeadController {
	return &LeadController{
		log:         log,
		validator:   validator,
		LeadService: leadService,
	}
}

func (c *LeadController) SaveLead(w http.ResponseWriter, r *http.Request) {
	const op = "LeadController.SaveLead"

	c.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	c.log.Debugf("%s: method id correct", op)

	var lead dto.LeadDTO
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		c.log.Infof("%s: decode error", op)

		responses.InvalidRequest(w)
		return
	}

	c.log.Debugf("%s: decode completed", op)

	// Валидируем данные
	if err := c.validator.Struct(lead); err != nil {
		c.log.Infof("%s: validation error", op)

		responses.ValidationError(w, utils.FormatValidationErrors(err))
		return
	}

	c.log.Debugf("%s: validation completed", op)

	err := c.LeadService.SaveLead(r.Context(), lead)
	if err != nil {
		c.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	responses.Ok(w)
}

func (c *LeadController) Leads(w http.ResponseWriter, r *http.Request) {
	const op = "LeadController.Leads"

	c.log.Debugf("%s: start", op)

	if r.Method != http.MethodGet {
		c.log.Infof("%s: method not allowed. method: %s", op, r.Method)

		w.Header().Set("Allow", http.MethodGet)
		responses.MethodNotAllowed(w)
		return
	}

	c.log.Debugf("%s: method allowed", op)

	// Парсим фильтры
	filter, err := parseLeadFilters(r)
	if err != nil {
		c.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	c.log.Debugf("%s: filters are received", op)

	userRolesValue := r.Context().Value(context_keys.UserRoleKey)
	userRoles, ok := userRolesValue.([]string)
	if !ok {
		c.log.Errorf("%s: user role not received", op)

		responses.ServerError(w)
		return
	}

	c.log.Debugf("%s: roles are received", op)

	if filter.UserID != nil && !utils.Contains(userRoles, "manager") {
		c.log.Infof("%s: forbidden", op)

		responses.Forbidden(w)
		return
	}

	c.log.Debugf("%s: rights checked", op)

	// Вызов сервиса
	leads, err := c.LeadService.Leads(r.Context(), filter)
	if err != nil {
		c.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	c.log.Debugf("%s: leads send", op)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}

func parseLeadFilters(r *http.Request) (dto.LeadFilterDTO, error) {
	query := r.URL.Query()

	parseInt := func(key string) (*int64, error) {
		if val := query.Get(key); val != "" {
			parsed, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid %s", key)
			}
			return &parsed, nil
		}
		return nil, nil
	}

	parseBool := func(key string) (*bool, error) {
		if val := query.Get(key); val != "" {
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid %s", key)
			}
			return &parsed, nil
		}
		return nil, nil
	}

	parseDate := func(key string) (*time.Time, error) {
		if val := query.Get(key); val != "" {
			parsed, err := time.Parse("2006-01-02", val)
			if err != nil {
				return nil, fmt.Errorf("invalid %s", key)
			}
			return &parsed, nil
		}
		return nil, nil
	}

	parseString := func(key string) *string {
		if val := query.Get(key); val != "" {
			return &val
		}
		return nil
	}

	statusID, err := parseInt("status_id")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	userID, err := parseInt("user_id")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	startDate, err := parseDate("start_date")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	endDate, err := parseDate("end_date")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	isInternet, err := parseBool("is_internet")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	isShipping, err := parseBool("is_shipping")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	isCleaning, err := parseBool("is_cleaning")
	if err != nil {
		return dto.LeadFilterDTO{}, err
	}

	limit := int64(10)
	if val := query.Get("limit"); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := int64(0)
	if val := query.Get("offset"); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	search := parseString("search")

	return dto.LeadFilterDTO{
		StatusID:   statusID,
		UserID:     userID,
		StartDate:  startDate,
		EndDate:    endDate,
		Limit:      limit,
		Offset:     offset,
		IsInternet: isInternet,
		IsShipping: isShipping,
		IsCleaning: isCleaning,
		Search:     search,
	}, nil
}
