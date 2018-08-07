package stub

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/stakater/Konfigurator/pkg/apis/konfigurator/v1alpha1"
	kContext "github.com/stakater/Konfigurator/pkg/context"
	"github.com/stakater/Konfigurator/pkg/controllers/ingress"
	"github.com/stakater/Konfigurator/pkg/controllers/konfiguratortemplate"
	"github.com/stakater/Konfigurator/pkg/controllers/pod"
	"github.com/stakater/Konfigurator/pkg/controllers/service"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func NewHandler(context *kContext.Context) sdk.Handler {
	return &Handler{
		Context: context,
	}
}

type Handler struct {
	Context *kContext.Context
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {

	case *v1alpha1.KonfiguratorTemplate:
		return h.HandleKonfiguratorTemplate(konfiguratortemplate.NewController(o, h.Context), event.Deleted)
	case *v1.Pod:
		return h.HandlePod(pod.NewController(o, h.Context), event.Deleted)
	case *v1.Service:
		return h.HandleService(service.NewController(o, h.Context), event.Deleted)
	case *v1beta1.Ingress:
		return h.HandleIngress(ingress.NewController(o, h.Context), event.Deleted)
	}

	return nil
}

func (h *Handler) HandleIngress(controller *ingress.Controller, deleted bool) error {
	if deleted {
		return controller.RemoveFromContext()
	}

	return controller.AddToContext()
}

func (h *Handler) HandleService(controller *service.Controller, deleted bool) error {
	if deleted {
		return controller.RemoveFromContext()
	}

	return controller.AddToContext()
}

func (h *Handler) HandlePod(controller *pod.Controller, deleted bool) error {
	if deleted {
		return controller.RemoveFromContext()
	}

	return controller.AddToContext()
}

func (h *Handler) HandleKonfiguratorTemplate(controller *konfiguratortemplate.Controller, deleted bool) error {
	if deleted {
		// Delegate delete calls to controller
		if err := controller.UnmountVolumes(); err != nil {
			return err
		}

		return controller.DeleteResources()
	}

	if err := controller.RenderTemplates(); err != nil {
		return err
	}
	if err := controller.CreateResources(); err != nil {
		return err
	}

	return controller.MountVolumes()
}
