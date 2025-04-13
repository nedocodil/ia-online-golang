package bitrix

import (
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/services/lead"
	"net/http"

	"github.com/sirupsen/logrus"
)

type BitrixController struct {
	log         *logrus.Logger
	bitrixKey   string
	LeadService lead.LeadServiceI
}

type BitrixControllerI interface {
	СhangingDeal(w http.ResponseWriter, r *http.Request)
}

func New(log *logrus.Logger, bitrix_key string, leadService lead.LeadServiceI) *BitrixController {
	return &BitrixController{
		log:         log,
		bitrixKey:   bitrix_key,
		LeadService: leadService,
	}
}

func (c *BitrixController) СhangingDeal(w http.ResponseWriter, r *http.Request) {
	const op = "BitrixController.СhangingDeal"

	c.log.Debugf("%s: start", op)

	if r.Method != http.MethodPost {
		c.log.Infof("%s: method not allowed. method: %s", op, r.Method)

		w.Header().Set("Allow", http.MethodPost)
		responses.MethodNotAllowed(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		c.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	var hook dto.OutgoingHook

	for i := 0; ; i++ {
		key := fmt.Sprintf("document_id[%d]", i)
		if val := r.FormValue(key); val != "" {
			hook.DocumentID = append(hook.DocumentID, val)
		} else {
			break
		}
	}

	hook.Auth.Domain = r.FormValue("auth[domain]")
	hook.Auth.ClientEndpoint = r.FormValue("auth[client_endpoint]")
	hook.Auth.ServerEndpoint = r.FormValue("auth[server_endpoint]")
	hook.Auth.MemberID = r.FormValue("auth[member_id]")

	c.log.Debugf("%s: parsing form", op)

	if hook.Auth.MemberID != c.bitrixKey {
		c.log.Infof("%s: invalid member id", op)

		responses.Forbidden(w)
		return
	}

	c.log.Debugf("%s: correct member id", op)

	err = c.LeadService.EditDeal(r.Context(), hook.DocumentID)
	if err != nil {
		c.log.Errorf("%s: %v", op, err)

		responses.ServerError(w)
		return
	}

	c.log.Debugf("%s: deal changed", op)

	responses.Ok(w)
}
