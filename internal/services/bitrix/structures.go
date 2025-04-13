package bitrix

import "time"

type ReturnDataDeal struct {
	Result InfoDeal `json:"result"`
	Time   TimeInfo `json:"time"`
}

type ReturnDataCreate struct {
	Result int      `json:"result"`
	Time   TimeInfo `json:"time"`
}

type InfoDeal struct {
	ID              string `json:"ID"`
	Title           string `json:"TITLE"`
	Status          string `json:"STAGE_ID"`
	ContactID       string `json:"CONTACT_ID"`
	InternetPayment string `json:"UF_CRM_1737451536004"`
	CleaningPayment string `json:"UF_CRM_1744353480781"`
	ShippingPayment string `json:"UF_CRM_1744354030686"`
}

type TimeInfo struct {
	Start            float64    `json:"start"`
	Finish           float64    `json:"finish"`
	Duration         float64    `json:"duration"`
	Processing       float64    `json:"processing"`
	Date_start       string     `json:"date_start"`
	Date_finish      string     `json:"date_finish"`
	OperatingResetAt *time.Time `json:"operating_reser_at"`
	Operating        float64    `json:"operating"`
}

type ErrorData struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
